package controllers

import (
	"context"
	"fmt"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	maxInstanceConditions = 10
)

// StatefulSetProperties defines the default agent and mysql versions
type StatefulSetProperties struct {
	AgentVersion string
	MySQLVersion string
}

// InstanceManager provides methods to manage the instance subcomponents
type InstanceManager struct {
	Context     context.Context
	Reconciler  *InstanceReconciler
	Properties  *StatefulSetProperties
	TimeManager *TimeManager
}

func (im *InstanceManager) setInstanceCondition(instance *mysqlv1alpha1.Instance, condition metav1.Condition) (ctrl.Result, error) {
	if condition.Reason == instance.Status.Reason {
		c := len(instance.Status.Conditions) - 1
		d := im.TimeManager.Next(instance.Status.Conditions[c].LastTransitionTime.Time)
		return ctrl.Result{Requeue: true, RequeueAfter: d}, nil
	}
	instance.Status.Ready = condition.Status
	instance.Status.Reason = condition.Reason
	instance.Status.Message = condition.Message
	conditions := append(instance.Status.Conditions, condition)
	if len(conditions) > maxInstanceConditions {
		conditions = conditions[1:]
	}
	instance.Status.Conditions = conditions
	log := im.Reconciler.Log.WithValues("namespace", instance.Namespace, "instance", instance.Name)
	log.Info("Updating instance with new Status", "Reason", condition.Reason, "Message", condition.Message)
	if err := im.Reconciler.Status().Update(im.Context, instance); err != nil {
		log.Error(err, "Unable to update instance")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (im *InstanceManager) getExporterSecret(instance *mysqlv1alpha1.Instance) (*corev1.Secret, error) {
	log := im.Reconciler.Log.WithValues("function", "getExporterSecret", "namespace", instance.Namespace, "instance", instance.Name)

	secretName := types.NamespacedName{
		Name:      instance.Name + "-exporter",
		Namespace: instance.Namespace,
	}
	secret := &corev1.Secret{}
	err := im.Reconciler.Client.Get(im.Context, secretName, secret)
	if err != nil && !errors.IsNotFound(err) {
		log.Error(err, "Error getting secret", "secret", secretName.Name)
	}
	if err != nil {
		log.Info("Secret does not exist", "secret", secretName.Name)
	}
	return secret, err
}

func (im *InstanceManager) createExporterSecret(instance *mysqlv1alpha1.Instance) (ctrl.Result, error) {
	log := im.Reconciler.Log.WithValues("function", "createExporterSecret", "namespace", instance.Namespace, "instance", instance.Name)

	labels := map[string]string{
		"app": instance.Name,
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-exporter",
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		StringData: map[string]string{
			".my.cnf": "[client]\nuser=exporter\npassword=exporter\nhost=localhost\n",
		},
	}
	if err := controllerutil.SetControllerReference(instance, secret, im.Reconciler.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	log.Info("Create secret", "secret", secret.Name)
	if err := im.Reconciler.Client.Create(im.Context, secret); err != nil {
		log.Error(err, "Secret creation failed", "secret", secret.Name)
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             mysqlv1alpha1.InstanceExporterSecretFailed,
			Message:            fmt.Sprintf("Secret exporter creation failed: %v", err),
		}
		return im.setInstanceCondition(instance, condition)
	}
	log.Info("Secret create succeeded", "secret", secret.Name)
	instance.Status.ExporterSecret = corev1.ObjectReference{
		Kind:            secret.Kind,
		Namespace:       secret.Namespace,
		Name:            secret.Name,
		UID:             secret.UID,
		APIVersion:      secret.APIVersion,
		ResourceVersion: secret.ResourceVersion,
	}
	condition := metav1.Condition{
		Type:               "available",
		Status:             metav1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             mysqlv1alpha1.InstanceExporterSecretCreated,
		Message:            "Secret exporter has been successfully created",
	}
	return im.setInstanceCondition(instance, condition)
}

func (im *InstanceManager) deleteExporterSecret(instance *mysqlv1alpha1.Instance, secret *corev1.Secret) (ctrl.Result, error) {
	log := im.Reconciler.Log.WithValues("function", "deleteExporterSecret", "namespace", instance.Namespace, "instance", instance.Name)

	log.Info("Delete secret", "secret", secret.Name)
	if err := im.Reconciler.Client.Delete(im.Context, secret); err != nil {
		log.Error(err, "Secret deletion failed", "secret", secret.Name)
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             mysqlv1alpha1.InstanceExporterSecretFailed,
			Message:            fmt.Sprintf("Secret deletion failed: %v", err),
		}
		return im.setInstanceCondition(instance, condition)
	}
	log.Info("Secret deletion succeeded", "secret", secret.Name)
	instance.Status.ExporterSecret = corev1.ObjectReference{}
	condition := metav1.Condition{
		Type:               "available",
		Status:             metav1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             mysqlv1alpha1.InstanceExporterSecretDeleted,
		Message:            "Secret exporter has been successfully deleted",
	}
	return im.setInstanceCondition(instance, condition)
}

func (im *InstanceManager) getStatefulSet(instance *mysqlv1alpha1.Instance) (*appsv1.StatefulSet, error) {
	log := im.Reconciler.Log.WithValues("function", "getStatefulSet", "namespace", instance.Namespace, "statefulset", instance.Name)

	stsName := types.NamespacedName{
		Name:      instance.Name,
		Namespace: instance.Namespace,
	}
	sts := &appsv1.StatefulSet{}
	err := im.Reconciler.Client.Get(im.Context, stsName, sts)
	if err != nil && !errors.IsNotFound(err) {
		log.Error(err, "Error getting statefulset", "statefulset", stsName.Name)
	}
	if err != nil {
		log.Info("Statefulset does not exist", "secret", stsName.Name)
	}
	return sts, err
}

func (im *InstanceManager) createStatefulSet(instance *mysqlv1alpha1.Instance, store *mysqlv1alpha1.Store, location string) (ctrl.Result, error) {
	log := im.Reconciler.Log.WithValues("function", "createStatefulSet", "namespace", instance.Namespace, "instance", instance.Name)

	sts := im.Properties.NewStatefulSetForInstance(instance, store, location)

	if err := controllerutil.SetControllerReference(instance, sts, im.Reconciler.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	log.Info("Create StatefulSet", "statefulset", sts.Name)
	if err := im.Reconciler.Client.Create(im.Context, sts); err != nil {
		log.Error(err, "Statfulset creation failed", "statefulset", sts.Name)
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             mysqlv1alpha1.InstanceStatefulSetFailed,
			Message:            fmt.Sprintf("Statfulset creation failed: %v", err),
		}
		return im.setInstanceCondition(instance, condition)
	}
	log.Info("Statfulset creation succeeded", "statefulset", sts.Name)
	instance.Status.StatefulSet = corev1.ObjectReference{
		Kind:            sts.Kind,
		Namespace:       sts.Namespace,
		Name:            sts.Name,
		UID:             sts.UID,
		APIVersion:      sts.APIVersion,
		ResourceVersion: sts.ResourceVersion,
	}
	condition := metav1.Condition{
		Type:               "available",
		Status:             metav1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             mysqlv1alpha1.InstanceStatefulSetCreated,
		Message:            "StatefulSet has been successfully created",
	}
	return im.setInstanceCondition(instance, condition)
}

// NewStatefulSetForInstance returns a MySQL StatefulSet with the instance name/namespace
func (s *StatefulSetProperties) NewStatefulSetForInstance(instance *mysqlv1alpha1.Instance, store *mysqlv1alpha1.Store, location string) *appsv1.StatefulSet {
	labels := map[string]string{
		"app": instance.Name,
	}
	diskSize := resource.NewQuantity(500*1024*1024, resource.BinarySI)
	restoreDiskSize := resource.NewQuantity(500*1024*1024, resource.BinarySI)
	var replicas int32 = 1
	initContainers := []corev1.Container{}
	if store != nil {
		if store.Spec.Envs != nil {
			env := store.Spec.Envs
			env = append(env, corev1.EnvVar{
				Name:  "AGT_BUCKET",
				Value: store.Spec.Bucket,
			})
			env = append(env, corev1.EnvVar{
				Name:  "AGT_PATH",
				Value: location,
			})
			env = append(env, corev1.EnvVar{
				Name:  "AGT_FILENAME",
				Value: "/docker-entrypoint-initdb.d/init-script.sql",
			})
			initContainers = []corev1.Container{
				{
					Name:  "restore",
					Image: "quay.io/blaqkube/mysql-agent:" + s.AgentVersion,
					Env:   env,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      instance.Name + "-init",
							MountPath: "/docker-entrypoint-initdb.d",
						},
					},
					Command: []string{
						"./mysql-agent",
						"init",
						"--restore",
					},
				},
			}
		}
	}
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			ServiceName: instance.Name,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mysql",
							Image: "mysql:" + s.MySQLVersion,
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_ALLOW_EMPTY_PASSWORD",
									Value: "1",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "mysql",
									ContainerPort: 3306,
								},
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{"mysqladmin", "ping"},
									},
								},
								InitialDelaySeconds: int32(30),
								TimeoutSeconds:      int32(5),
								PeriodSeconds:       int32(10),
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{"mysql", "-h", "127.0.0.1", "-e", "SELECT 1"},
									},
								},
								InitialDelaySeconds: int32(30),
								TimeoutSeconds:      int32(5),
								PeriodSeconds:       int32(10),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      instance.Name + "-data",
									MountPath: "/var/lib/mysql",
								},
								{
									Name:      instance.Name + "-init",
									MountPath: "/docker-entrypoint-initdb.d",
								},
							},
						},
						{
							Name:  "agent",
							Image: "quay.io/blaqkube/mysql-agent:" + s.AgentVersion,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      instance.Name + "-data",
									MountPath: "/var/lib/mysql",
								},
							},
							Command: []string{
								"./mysql-agent",
								"serve",
							},
						},
						{
							Name:  "exporter",
							Image: "prom/mysqld-exporter:v0.12.1",
							Ports: []corev1.ContainerPort{
								{
									Name:          "prom-mysql",
									ContainerPort: 9104,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      instance.Name + "-exporter",
									MountPath: "/home",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: instance.Name + "-exporter",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: instance.Name + "-exporter",
									Items: []corev1.KeyToPath{
										{
											Key:  ".my.cnf",
											Path: ".my.cnf",
										},
									},
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: instance.Name + "-data",
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: *diskSize,
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: instance.Name + "-init",
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: *restoreDiskSize,
							},
						},
					},
				},
			},
		},
	}
	if store != nil {
		sts.Spec.Template.Spec.InitContainers = initContainers
	}
	if instance.Spec.Database != "" {
		sts.Spec.Template.Spec.Containers[0].Env = append(
			sts.Spec.Template.Spec.Containers[0].Env,
			corev1.EnvVar{
				Name:  "MYSQL_DATABASE",
				Value: instance.Spec.Database,
			},
		)
	}
	return sts
}
