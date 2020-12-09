package controllers

import (
	"context"
	"fmt"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	"github.com/blaqkube/mysql-operator/mysql-operator/helpers"
	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/api/errors"
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
		if store.Spec.Backend == nil || *store.Spec.Backend == "s3" {

			s, err := r.BackupStore.New(store.Spec.S3.Env)
			if err != nil {
				store.Status.Reason = "Error"
				log.Error(err, "Unable to Initialize store, setting to Error")
				if err := r.Status().Update(ctx, &store); err != nil {
					log.Error(err, "Unable to update store status to Error")
					return ctrl.Result{}, err
				}
				return ctrl.Result{}, nil
			}
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
func (r *StoreReconciler) GetEnvVars(ctx context.Context, name types.NamespacedName, Env []mysqlv1alpha1.EnvVar) (map[string]string, error) {
	output := map[string]string{}
	for _, envVar := range Env {
		if envVar.Name == "" {
			continue
		}
		if envVar.Value != "" {
			output[envVar.Name] = envVar.Value
			continue
		}
		if envVar.ValueFrom != nil {
			switch {
			case envVar.ValueFrom.ConfigMapKeyRef != nil:
				cm := envVar.ValueFrom.ConfigMapKeyRef
				name := cm.Name
				key := cm.Key
				optional := cm.Optional != nil && *cm.Optional
				configMap, ok := configMaps[name]
				if !ok {
					if kl.kubeClient == nil {
						return result, fmt.Errorf("couldn't get configMap %v/%v, no kubeClient defined", pod.Namespace, name)
					}
					configMap, err = kl.configMapManager.GetConfigMap(pod.Namespace, name)
					if err != nil {
						if errors.IsNotFound(err) && optional {
							// ignore error when marked optional
							continue
						}
						return result, err
					}
					configMaps[name] = configMap
				}
				runtimeVal, ok = configMap.Data[key]
				if !ok {
					if optional {
						continue
					}
					return result, fmt.Errorf("couldn't find key %v in ConfigMap %v/%v", key, pod.Namespace, name)
				}
			case envVar.ValueFrom.SecretKeyRef != nil:
				s := envVar.ValueFrom.SecretKeyRef
				name := s.Name
				key := s.Key
				optional := s.Optional != nil && *s.Optional
				secret, ok := secrets[name]
				if !ok {
					if kl.kubeClient == nil {
						return result, fmt.Errorf("couldn't get secret %v/%v, no kubeClient defined", pod.Namespace, name)
					}
					secret, err = kl.secretManager.GetSecret(pod.Namespace, name)
					if err != nil {
						if errors.IsNotFound(err) && optional {
							// ignore error when marked optional
							continue
						}
						return result, err
					}
					secrets[name] = secret
				}
				runtimeValBytes, ok := secret.Data[key]
				if !ok {
					if optional {
						continue
					}
					return result, fmt.Errorf("couldn't find key %v in Secret %v/%v", key, pod.Namespace, name)
				}
				runtimeVal = string(runtimeValBytes)
			}
		}
	}
	return output, nil
}
