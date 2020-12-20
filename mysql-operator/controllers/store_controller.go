package controllers

import (
	"context"
	"errors"
	"fmt"
	"os"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/blaqkube/mysql-operator/agent/backend"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StoreReconciler reconciles a Store object
type StoreReconciler struct {
	client.Client
	Log     logr.Logger
	Scheme  *runtime.Scheme
	Storage backend.Storage
}

const (
	// StoreReasonCheckRequest set the reason for the change in conditions to check requested
	StoreReasonCheckRequest = "CheckRequest"

	maxConditions = 10
)

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=stores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=stores/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=stores/finalizers,verbs=update

func setStoreCondition(ctx context.Context, r *StoreReconciler, store *mysqlv1alpha1.Store, condition metav1.Condition) (ctrl.Result, error) {
	store.Status.Ready = condition.Status
	store.Status.Reason = condition.Reason
	store.Status.Message = condition.Message
	conditions := append(store.Status.Conditions, condition)
	if len(conditions) > maxConditions {
		conditions = conditions[1:]
	}
	store.Status.Conditions = conditions
	log := r.Log.WithValues("store", store.Namespace+"."+store.Name)
	log.Info("Updating store with new Status", "Reason", condition.Reason, "Message", condition.Message)
	if err := r.Status().Update(ctx, store); err != nil {
		log.Error(err, "Unable to update store")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func initTestFile() (*string, error) {
	name := ".mysql-operator.out"
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	_, err = fmt.Fprintf(file, "[blaqkube]")
	return &name, err
}

// Reconcile implement the reconciliation loop for stores
func (r *StoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("store", req.NamespacedName)

	// TODO:
	// - Reconciler should be able to detect a change in the ConfigMap or
	//   Secret and reload the associated data
	var store mysqlv1alpha1.Store
	if err := r.Get(ctx, req.NamespacedName, &store); err != nil {
		log.Error(err, "unable to fetch Store")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if store.Status.Reason == "" || store.Status.CheckRequested == true {
		store.Status.CheckRequested = false
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionUnknown,
			LastTransitionTime: metav1.Now(),
			Reason:             mysqlv1alpha1.StateCheckRequested,
			Message:            "A new check has been requested",
		}
		return setStoreCondition(ctx, r, &store, condition)
	}

	if store.Status.Reason == mysqlv1alpha1.StateCheckRequested {
		if store.Spec.Backend == nil || *store.Spec.Backend == "s3" {
			store.Status.CheckRequested = false
			envs, err := r.GetEnvVars(ctx, store)
			if err != nil {
				condition := metav1.Condition{
					Type:               "available",
					Status:             metav1.ConditionFalse,
					LastTransitionTime: metav1.Now(),
					Reason:             mysqlv1alpha1.StateCheckFailed,
					Message:            "Cannot access values for envs",
				}
				return setStoreCondition(ctx, r, &store, condition)
			}
			e := []openapi.EnvVar{}
			for k := range envs {
				e = append(e, openapi.EnvVar{Name: k, Value: envs[k]})
			}
			filename, err := initTestFile()
			request := &openapi.BackupRequest{
				Bucket:   store.Spec.Bucket,
				Location: "/blaqkube/.mysql-operator.out",
				Envs:     e,
			}
			if err != nil {
				condition := metav1.Condition{
					Type:               "available",
					Status:             metav1.ConditionFalse,
					LastTransitionTime: metav1.Now(),
					Reason:             mysqlv1alpha1.StateCheckFailed,
					Message:            fmt.Sprintf("Cannot initialize local file, error: %v", err),
				}
				return setStoreCondition(ctx, r, &store, condition)
			}
			err = r.Storage.Push(request, *filename)
			if err == nil {
				err = r.Storage.Delete(request)
			}
			if err != nil {
				condition := metav1.Condition{
					Type:               "available",
					Status:             metav1.ConditionFalse,
					LastTransitionTime: metav1.Now(),
					Reason:             mysqlv1alpha1.StateCheckFailed,
					Message:            fmt.Sprintf("Cannot write to bucket, error: %v", err),
				}
				return setStoreCondition(ctx, r, &store, condition)
			}
			condition := metav1.Condition{
				Type:               "available",
				Status:             metav1.ConditionTrue,
				LastTransitionTime: metav1.Now(),
				Reason:             mysqlv1alpha1.StateCheckSucceeded,
				Message:            "The check has succeeded",
			}
			return setStoreCondition(ctx, r, &store, condition)
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
	for _, envVar := range store.Spec.Envs {
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
					configMap = corev1.ConfigMap{}
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
					secret = corev1.Secret{}
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
			r.Log.Info("Get SecretRef", "variable", envVar.Name, "value", value)
			output[envVar.Name] = value
			continue
		}
		return nil, errors.New("MissingVariable")
	}
	r.Log.Info("Map with all values built, continue...")
	return output, nil
}
