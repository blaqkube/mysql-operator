package controllers

import (
	"context"

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

// StatefulSetProperties defines an interface with the default agent/mysql versions
type StatefulSetProperties struct {
	AgentVersion string
	MySQLVersion string
}

// CreateOrUpdateStafefulSet creates a secret for the instance
func (r *InstanceReconciler) CreateOrUpdateStafefulSet(instance *mysqlv1alpha1.Instance, store *mysqlv1alpha1.Store, filePath string) (ctrl.Result, error) {
	ctx := context.Background()

	secretName := types.NamespacedName{Name: instance.Name + "-exporter", Namespace: instance.Namespace}
	secret := &corev1.Secret{}
	err := r.Client.Get(ctx, secretName, secret)
	if err != nil && !errors.IsNotFound(err) {
		return ctrl.Result{}, err
	}
	if err != nil {
		r.Log.Info("Creating a new Secret", "Secret.Namespace", instance.Namespace, "secret.Name", instance.Name+"-exporter")
		newSecret := r.Properties.NewSecretForInstance(instance)
		if err := controllerutil.SetControllerReference(instance, newSecret, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Client.Create(ctx, newSecret); err != nil {
			return ctrl.Result{}, err
		}
	}

	sts := r.Properties.NewStatefulSetForInstance(instance, store, filePath)
	// Check if this StatefulSet already exists
	found := &appsv1.StatefulSet{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: sts.Name, Namespace: sts.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("Creating a new StatefulSet", "StatefulSet.Namespace", sts.Namespace, "StatefulSet.Name", sts.Name)
		if err := controllerutil.SetControllerReference(instance, sts, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		err = r.Client.Create(ctx, sts)
		if err != nil {
			return ctrl.Result{}, err
		}
		// StatefulSet created successfully - don't requeue
		instance.Status.LastCondition = "Success"
		if err := r.Status().Update(ctx, instance); err != nil {
			r.Log.Error(err, "unable to update instance status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// StatefulSet already exists - don't requeue
	r.Log.Info("Skip reconcile: StatefulSet already exists", "StatefulSet.Namespace", found.Namespace, "StatefulSet.Name", found.Name)
	return ctrl.Result{}, nil
}

// NewSecretForInstance returns a secret that stores the mysql configuration file
func (s *StatefulSetProperties) NewSecretForInstance(instance *mysqlv1alpha1.Instance) *corev1.Secret {
	labels := map[string]string{
		"app": instance.Name,
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-exporter",
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Data: map[string][]byte{
			".my.cnf": []byte("[client]\nuser=exporter\npassword=exporter\nhost=localhost\n"),
		},
	}
	return secret
}

// NewStatefulSetForInstance returns a MySQL StatefulSet with the instance name/namespace
func (s *StatefulSetProperties) NewStatefulSetForInstance(instance *mysqlv1alpha1.Instance, store *mysqlv1alpha1.Store, filePath string) *appsv1.StatefulSet {
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
				Value: filePath,
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
