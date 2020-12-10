package controllers

import (
	"context"
	"errors"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	"github.com/blaqkube/mysql-operator/mysql-operator/helpers"
	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StoreReconciler reconciles a Store object
type StoreReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	BackupStore helpers.StoreInitializer
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=stores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=stores/status,verbs=get;update;patch

// Reconcile implement the reconciliation loop for stores
func (r *StoreReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("store", req.NamespacedName)

	// your logic here
	var store mysqlv1alpha1.Store
	if err := r.Get(ctx, req.NamespacedName, &store); err != nil {
		log.Error(err, "unable to fetch Store")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if store.Status.Reason == "" {
		store.Status.Reason = "Initializing"
		if err := r.Status().Update(ctx, &store); err != nil {
			log.Error(err, "Unable to update store status to Initializing")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}
	if store.Status.Reason == "Initializing" {
		if (store.Spec.Backend == nil || *store.Spec.Backend == "s3") && store.Spec.S3 != nil {

			keys, err := r.GetEnvVars(ctx, store)
			if err != nil {
				store.Status.Reason = "Error"
				if err.Error() == "MissingVariable" {
					store.Status.Reason = "MissingVariable"
				}
				log.Errorf("Error reading environment variable: %v", err)
				if err := r.Status().Update(ctx, &store); err != nil {
					log.Errorf("Unable to update store status: %s", err)
					return ctrl.Result{}, err
				}
				return ctrl.Result{}, nil
			}
			accessKey, _ := keys["AWS_ACCESS_KEY_ID"]
			secretKey, _ := keys["AWS_SECRET_ACCESS_KEY"]
			region, ok := keys["AWS_DEFAULT_REGION"]
			if !ok {
				region, ok = keys["AWS_REGION"]
			}
			s, err := r.BackupStore.New(&helpers.AWSConfig{AccessKey: accessKey, SecretKey: secretKey, Region: region})
			err = s.TestS3Access(store.Spec.S3.Bucket, "/validation")
			if err != nil {
				store.Status.Reason = "AccessDenied"
				log.Error(err, "Unable to access store, setting to AccessDenied")
				if err := r.Status().Update(ctx, &store); err != nil {
					log.Error(err, "Unable to update store status to AccessDenied")
					return ctrl.Result{}, err
				}
				return ctrl.Result{}, nil
			}
			store.Status.Reason = "Succeeded"
			log.Error(err, "Access to store Succeeded")
			if err := r.Status().Update(ctx, &store); err != nil {
				log.Error(err, "Unable to update store status to Succeeded")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager configure type of events the manager should watch
func (r *StoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Store{}).
		Complete(r)
}

// GetEnvVars returns the environment variables for the store
func (r *StoreReconciler) GetEnvVars(ctx context.Context, store mysqlv1alpha1.Store) (map[string]string, error) {
	output := map[string]string{}
	configMaps := map[string]corev1.ConfigMap{}
	secrets := map[string]corev1.Secret{}
	for _, envVar := range store.Spec.S3.Env {
		if envVar.Name == "" {
			return nil, errors.New("MissingVariable")
		}
		if envVar.Value != "" {
			output[envVar.Name] = envVar.Value
			continue
		}
		if envVar.ValueFrom != nil {
			value := ""
			switch {
			case envVar.ValueFrom.ConfigMapKeyRef != nil:
				cm := envVar.ValueFrom.ConfigMapKeyRef
				namespace := store.Namespace
				name := cm.Name
				key := cm.Key
				optional := cm.Optional != nil && *cm.Optional
				configMap, ok := configMaps[name]
				if !ok {
					configMap := corev1.ConfigMap{}
					if err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, &configMap); err != nil {
						if !optional {
							return nil, err
						}
					} else {
						configMaps[name] = configMap
						ok = true
					}
				}
				if ok {
					value, ok = configMap.Data[key]
					if !ok && !optional {
						return nil, errors.New("MissingVariable")
					}
				}
			case envVar.ValueFrom.SecretKeyRef != nil:
				s := envVar.ValueFrom.SecretKeyRef
				namespace := store.Namespace
				name := s.Name
				key := s.Key
				optional := s.Optional != nil && *s.Optional
				secret, ok := secrets[name]
				if !ok {
					secret := corev1.Secret{}
					if err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, &secret); err != nil {
						if !optional {
							return nil, err
						}
					} else {
						secrets[name] = secret
						ok = true
					}
				}
				if ok {
					valueBytes, ok := secret.Data[key]
					if !ok && !optional {
						return nil, errors.New("MissingVariable")
					}
					if ok {
						value = string(valueBytes)
					}
				}
			}
			output[envVar.Name] = value
			continue
		}
		return nil, errors.New("MissingVariable")
	}
	return output, nil
}
