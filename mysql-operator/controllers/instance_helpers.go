package controllers

import (
	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StatefulSetProperties struct {
	AgentVersion string
	MySQLVersion string
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
		if store.Spec.S3Backup != nil && store.Spec.S3Backup.AWSConfig != nil {
			initContainers = []corev1.Container{
				{
					Name:  "restore",
					Image: "quay.io/blaqkube/mysql-agent:" + s.AgentVersion,
					Env: []corev1.EnvVar{
						corev1.EnvVar{
							Name:  "AWS_REGION",
							Value: store.Spec.S3Backup.AWSConfig.Region,
						},
						corev1.EnvVar{
							Name:  "AWS_ACCESS_KEY_ID",
							Value: store.Spec.S3Backup.AWSConfig.AccessKey,
						},
						corev1.EnvVar{
							Name:  "AWS_SECRET_ACCESS_KEY",
							Value: store.Spec.S3Backup.AWSConfig.SecretKey,
						},
						corev1.EnvVar{
							Name:  "AGT_BUCKET",
							Value: store.Spec.S3Backup.Bucket,
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
									Name:      instance.Name + "-data",
									MountPath: "/var/lib/mysql",
								},
								corev1.VolumeMount{
									Name:      instance.Name + "-init",
									MountPath: "/docker-entrypoint-initdb.d",
								},
							},
						},
						{
							Name:  "agent",
							Image: "quay.io/blaqkube/mysql-agent:" + s.AgentVersion,
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      instance.Name + "-data",
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
									Name:      instance.Name + "-exporter",
									MountPath: "/home",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: instance.Name + "-exporter",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: instance.Name + "-exporter",
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
				corev1.PersistentVolumeClaim{
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
