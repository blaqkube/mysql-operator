package controllers

import (
	"context"
	"errors"
	"fmt"
	"os"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	maxStoreConditions = 10
)

// StoreManager provides methods to manage the store subcomponents
type StoreManager struct {
	Context    context.Context
	Reconciler *StoreReconciler
}

func (sm *StoreManager) setStoreCondition(store *mysqlv1alpha1.Store, condition metav1.Condition) (ctrl.Result, error) {
	store.Status.Ready = condition.Status
	store.Status.Reason = condition.Reason
	store.Status.Message = condition.Message
	conditions := append(store.Status.Conditions, condition)
	if len(conditions) > maxStoreConditions {
		conditions = conditions[1:]
	}
	store.Status.Conditions = conditions
	log := sm.Reconciler.Log.WithValues("namespace", store.Namespace, "store", store.Name)
	log.Info("Updating store with new Status", "Reason", condition.Reason, "Message", condition.Message)
	if err := sm.Reconciler.Status().Update(sm.Context, store); err != nil {
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

// GetEnvVars returns the environment variables for the store
func (sm *StoreManager) GetEnvVars(store mysqlv1alpha1.Store) (map[string]string, error) {
	log := sm.Reconciler.Log.WithValues("store", types.NamespacedName{Namespace: store.Namespace, Name: store.Name})
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
					if err := sm.Reconciler.Get(sm.Context, types.NamespacedName{Namespace: namespace, Name: name}, &configMap); err != nil {
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
					if err := sm.Reconciler.Get(sm.Context, types.NamespacedName{Namespace: namespace, Name: name}, &secret); err != nil {
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
			log.Info("Reference Key", "variable", envVar.Name, "value", "***")
			output[envVar.Name] = value
			continue
		}
		return nil, errors.New("MissingVariable")
	}
	return output, nil
}
