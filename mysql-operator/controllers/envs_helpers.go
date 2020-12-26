package controllers

import (
	"context"
	"errors"
	"fmt"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EnvManager reconciles an Environment object
type EnvManager struct {
	client.Client
	Log logr.Logger
}

var (
	// ErrMissingVariable shows when a variable is missing from the store
	ErrMissingVariable = errors.New("MissingVariable")
)

// GetEnvVars returns the environment variables for the store
func (em *EnvManager) GetEnvVars(ctx context.Context, store mysqlv1alpha1.Store) (map[string]string, error) {
	log := em.Log.WithValues("store", types.NamespacedName{Namespace: store.Namespace, Name: store.Name})
	output := map[string]string{}
	configMaps := map[string]corev1.ConfigMap{}
	secrets := map[string]corev1.Secret{}
	for _, envVar := range store.Spec.Envs {
		if envVar.Name == "" {
			return nil, ErrMissingVariable
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
					if err := em.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, &configMap); err != nil {
						if !optional {
							log.Info(fmt.Sprintf("Error getting configmap, %v", err), "configmap", name)
							return nil, ErrMissingVariable
						}
					} else {
						configMaps[name] = configMap
						ok = true
					}
				}
				if ok {
					value, ok = configMap.Data[key]
					if !ok && !optional {
						return nil, ErrMissingVariable
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
					if err := em.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, &secret); err != nil {
						if !optional {
							log.Info(fmt.Sprintf("Error getting secret, %v", err), "secret", name)
							return nil, ErrMissingVariable
						}
					} else {
						secrets[name] = secret
						ok = true
					}
				}
				if ok {
					valueBytes, ok := secret.Data[key]
					if !ok && !optional {
						return nil, ErrMissingVariable
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
		return nil, ErrMissingVariable
	}
	return output, nil
}
