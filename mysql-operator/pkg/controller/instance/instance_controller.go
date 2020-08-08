package instance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/pkg/apis/mysql/v1alpha1"
	"github.com/robfig/cron/v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	log   = logf.Log.WithName("controller_instance")
	crond = cron.New()
)

func init() {
	reqLogger := log.WithValues("Controller", "instance")
	reqLogger.Info("Start crond")
	crond.Start()
}

// Add creates a new Instance Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileInstance{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("instance-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Instance
	err = c.Watch(&source.Kind{Type: &mysqlv1alpha1.Instance{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Instance
	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mysqlv1alpha1.Instance{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileInstance implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileInstance{}

// ReconcileInstance reconciles a Instance object
type ReconcileInstance struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Instance object and makes changes based on the state read
// and what is in the Instance.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileInstance) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Instance")

	// Fetch the Instance instance
	instance := &mysqlv1alpha1.Instance{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	storeName := instance.Spec.Restore.Store
	filePath := instance.Spec.Restore.FilePath
	store := &mysqlv1alpha1.Store{}
	if storeName != "" {
		if filePath == "" {
			if instance.Status.LastCondition == "Error" {
				return reconcile.Result{}, nil
			}
			t := metav1.Now()
			condition := mysqlv1alpha1.ConditionStatus{
				LastProbeTime: &t,
				Status:        "Error",
				Message:       "Restore file filePath should not be empty",
			}
			instance.Status.Conditions = []mysqlv1alpha1.ConditionStatus{condition}
			err = r.client.Status().Update(context.TODO(), instance)
			if err != nil {
				return reconcile.Result{}, err
			}
			// Store updated successfully - don't requeue
			return reconcile.Result{}, nil
		}
		err := r.client.Get(
			context.TODO(),
			client.ObjectKey{Namespace: request.Namespace, Name: storeName},
			store,
		)
		if err != nil {
			if instance.Status.LastCondition == "Error" {
				return reconcile.Result{}, nil
			}
			instance.Status.LastCondition = "Error"
			t := metav1.Now()
			condition := mysqlv1alpha1.ConditionStatus{
				LastProbeTime: &t,
				Status:        "Error",
				Message:       fmt.Sprintf("Could not find store %s", storeName),
			}
			instance.Status.Conditions = []mysqlv1alpha1.ConditionStatus{condition}
			err = r.client.Status().Update(context.TODO(), instance)
			if err != nil {
				return reconcile.Result{}, err
			}
			// Store updated successfully - don't requeue
			return reconcile.Result{}, nil
		}
	}
	if storeName == "" {
		store = nil
	}
	statefulSet := newStatefulSetForCR(instance, store, filePath)

	// Set Instance instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, statefulSet, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	secret := &corev1.Secret{}
	err = r.client.Get(
		context.TODO(),
		types.NamespacedName{Name: statefulSet.Name + "-exporter", Namespace: statefulSet.Namespace},
		secret,
	)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info(
			"Creating a new Secret",
			"Secret.Namespace",
			statefulSet.Namespace,
			"secret.Name",
			statefulSet.Name+"-exporter",
		)
		labels := map[string]string{
			"app": statefulSet.Name,
		}
		secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      statefulSet.Name + "-exporter",
				Namespace: statefulSet.Namespace,
				Labels:    labels,
			},
			Data: map[string][]byte{
				".my.cnf": []byte("[client]\nuser=exporter\npassword=exporter\nhost=localhost\n"),
			},
		}

		err = r.client.Create(context.TODO(), secret)
		if err != nil {
			return reconcile.Result{}, err
		}
		if err := controllerutil.SetControllerReference(instance, secret, r.scheme); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Check if this StatefulSet already exists
	found := &appsv1.StatefulSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: statefulSet.Name, Namespace: statefulSet.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new StatefulSet", "StatefulSet.Namespace", statefulSet.Namespace, "StatefulSet.Name", statefulSet.Name)
		err = r.client.Create(context.TODO(), statefulSet)
		if err != nil {
			return reconcile.Result{}, err
		}
		if instance.Spec.Maintenance.Backup && len(instance.Spec.Maintenance.WindowStart) >= 5 {
			reqLogger.Info("Create backup schedule", "StatefulSet.Namespace", statefulSet.Namespace, "StatefulSet.Name", statefulSet.Name)
			hour, _ := strconv.Atoi(instance.Spec.Maintenance.WindowStart[0:2])
			min, _ := strconv.Atoi(instance.Spec.Maintenance.WindowStart[3:5])
			//TODO: manage errors
			crond.AddFunc(fmt.Sprintf("%d %d * * *", min, hour), func() {
				currentTime := time.Now()
				reqLogger.Info(
					fmt.Sprintf("Create backup schedule at %s", currentTime.Format("2006.01.02 15:04:05")),
					"StatefulSet.Namespace", statefulSet.Namespace, "StatefulSet.Name", statefulSet.Name)
				r.KickBackup(instance, instance.Spec.Maintenance.BackupStore)
			})

		}
		// StatefulSet created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// StatefulSet already exists - don't requeue
	reqLogger.Info("Skip reconcile: StatefulSet already exists", "StatefulSet.Namespace", found.Namespace, "StatefulSet.Name", found.Name)
	return reconcile.Result{}, nil
}

// newStatefulSetForCR returns a busybox pod with the same name/namespace as the cr
func newStatefulSetForCR(cr *mysqlv1alpha1.Instance, store *mysqlv1alpha1.Store, filePath string) *appsv1.StatefulSet {
	tag := "115aaeb"
	labels := map[string]string{
		"app": cr.Name,
	}
	diskSize := resource.NewQuantity(500*1024*1024, resource.BinarySI)
	restoreDiskSize := resource.NewQuantity(500*1024*1024, resource.BinarySI)
	var replicas int32 = 1
	initContainers := []corev1.Container{}
	if store != nil {
		initContainers = []corev1.Container{
			{
				Name:  "restore",
				Image: "quay.io/blaqkube/mysql-agent:" + tag,
				Env: []corev1.EnvVar{
					corev1.EnvVar{
						Name:  "AWS_REGION",
						Value: store.Spec.S3Access.AWSConfig.Region,
					},
					corev1.EnvVar{
						Name:  "AWS_ACCESS_KEY_ID",
						Value: store.Spec.S3Access.AWSConfig.AccessKey,
					},
					corev1.EnvVar{
						Name:  "AWS_SECRET_ACCESS_KEY",
						Value: store.Spec.S3Access.AWSConfig.SecretKey,
					},
					corev1.EnvVar{
						Name:  "AGT_BUCKET",
						Value: store.Spec.S3Access.Bucket,
					},
					corev1.EnvVar{
						Name:  "AGT_PATH",
						Value: filePath,
					},
					corev1.EnvVar{
						Name:  "AGT_FILENAME",
						Value: "/docker-entrypoint-initdb.d/init-script.sql",
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					corev1.VolumeMount{
						Name:      cr.Name + "-init",
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
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			ServiceName: cr.Name,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mysql",
							Image: "mysql:8.0.20",
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "MYSQL_ALLOW_EMPTY_PASSWORD",
									Value: "1",
								},
							},
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
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
								corev1.VolumeMount{
									Name:      cr.Name + "-data",
									MountPath: "/var/lib/mysql",
								},
								corev1.VolumeMount{
									Name:      cr.Name + "-init",
									MountPath: "/docker-entrypoint-initdb.d",
								},
							},
						},
						{
							Name:  "agent",
							Image: "quay.io/blaqkube/mysql-agent:" + tag,
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      cr.Name + "-data",
									MountPath: "/var/lib/mysql",
								},
							},
						},
						{
							Name:  "exporter",
							Image: "prom/mysqld-exporter:v0.12.1",
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									Name:          "prom-mysql",
									ContainerPort: 9104,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      cr.Name + "-exporter",
									MountPath: "/home",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: cr.Name + "-exporter",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: cr.Name + "-exporter",
									Items: []corev1.KeyToPath{
										corev1.KeyToPath{
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
				corev1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{
						Name: cr.Name + "-data",
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
				corev1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{
						Name: cr.Name + "-init",
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
	if cr.Spec.Database != "" {
		sts.Spec.Template.Spec.Containers[0].Env = append(
			sts.Spec.Template.Spec.Containers[0].Env,
			corev1.EnvVar{
				Name:  "MYSQL_DATABASE",
				Value: cr.Spec.Database,
			},
		)
	}
	return sts
}

func (r *ReconcileInstance) KickBackup(instance *mysqlv1alpha1.Instance, store string) {
	currentTime := time.Now()
	b := &mysqlv1alpha1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "automatic-" + instance.ObjectMeta.Name + fmt.Sprintf("-%s", currentTime.Format("20060102-1504")),
			Namespace: instance.ObjectMeta.Namespace,
			Labels:    instance.ObjectMeta.Labels,
		},
		Spec: mysqlv1alpha1.BackupSpec{
			Store:    store,
			Instance: instance.ObjectMeta.Name,
		},
	}
	err := r.client.Create(context.TODO(), b)
	if err != nil {
		fmt.Printf(
			"Error creating backup %s: %v",
			"automatic-"+instance.ObjectMeta.Name+fmt.Sprintf("-%s", currentTime.Format("20060102-1504")),
			err,
		)
	}
}
