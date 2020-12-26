package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/blaqkube/mysql-operator/mysql-operator/agent"
)

const (
	maxUserConditions = 10
)

var (
	// ErrMissingPassword is reported when a password is missing for a user
	ErrMissingPassword = errors.New("MissingPassword")

	// ErrKeyNotFound is reported when a key is not found in configMap or secret
	ErrKeyNotFound = errors.New("KeyNotFound")

	// ErrConfigMapNotFound is reported when a configMap is not found
	ErrConfigMapNotFound = errors.New("ConfigMapNotFound")

	// ErrSecretNotFound is reported when a secret is not found
	ErrSecretNotFound = errors.New("SecretNotFound")
)

// UserManager provides methods to manage user subcomponents
type UserManager struct {
	Context     context.Context
	Reconciler  *UserReconciler
	TimeManager *TimeManager
}

func (um *UserManager) setUserCondition(user *mysqlv1alpha1.User, condition metav1.Condition) (ctrl.Result, error) {
	if condition.Reason == user.Status.Reason {
		c := len(user.Status.Conditions) - 1
		d := um.TimeManager.Next(user.Status.Conditions[c].LastTransitionTime.Time)
		if condition.Reason != mysqlv1alpha1.DatabaseSucceeded {
			return ctrl.Result{Requeue: true, RequeueAfter: d}, nil
		}
		return ctrl.Result{}, nil
	}
	user.Status.Ready = condition.Status
	user.Status.Reason = condition.Reason
	user.Status.Message = condition.Message
	conditions := append(user.Status.Conditions, condition)
	if len(conditions) > maxUserConditions {
		conditions = conditions[1:]
	}
	user.Status.Conditions = conditions
	log := um.Reconciler.Log.WithValues("namespace", user.Namespace, "user", user.Name)
	log.Info("Updating user with new Status", "Reason", condition.Reason, "Message", condition.Message)
	if err := um.Reconciler.Status().Update(um.Context, user); err != nil {
		log.Error(err, "Unable to update user")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// CreateUser is the script that creates a user
func (um *UserManager) CreateUser(user *mysqlv1alpha1.User) error {
	password, err := um.GetPassword(user)
	if err != nil {
		return err
	}
	log := um.Reconciler.Log.WithValues("namespace", user.Namespace, "user", user.Name)
	a := &APIReconciler{
		Client: um.Reconciler.Client,
		Log:    um.Reconciler.Log,
	}
	api, err := a.GetAPI(
		um.Context,
		types.NamespacedName{
			Name:      user.Spec.Instance,
			Namespace: user.Namespace,
		},
	)
	if err != nil {
		return err
	}

	payload := agent.User{
		Username: user.Spec.Username,
		Password: password,
	}

	_, response, err := api.MysqlApi.CreateUser(um.Context, payload, nil)
	if err != nil || response == nil {
		msg := "NoResponse"
		if err != nil {
			msg = err.Error()
		}
		log.Info(fmt.Sprintf("Could not access agent, error: %s", msg))
		return ErrAgentAccessFailed
	}
	if response.StatusCode != http.StatusCreated {
		log.Info("Agent returned unexpected response", "httpcode", response.StatusCode)
		return ErrAgentRequestFailed
	}
	return nil
}

// GetPassword returns the password for the user
func (um *UserManager) GetPassword(user *mysqlv1alpha1.User) (string, error) {
	log := um.Reconciler.Log.WithValues("user", types.NamespacedName{Namespace: user.Namespace, Name: user.Name})
	if user.Spec.Password != "" {
		return user.Spec.Password, nil
	}

	if user.Spec.PasswordFrom != nil {
		switch {
		case user.Spec.PasswordFrom.ConfigMapKeyRef != nil:
			cm := user.Spec.PasswordFrom.ConfigMapKeyRef
			namespace := user.Namespace
			name := cm.Name
			key := cm.Key
			configMap := corev1.ConfigMap{}
			if err := um.Reconciler.Get(um.Context, types.NamespacedName{Namespace: namespace, Name: name}, &configMap); err != nil {
				log.Info(fmt.Sprintf("Read Configmap, error: %s", err), "namespace", namespace, "configmap", name)
				return "", ErrConfigMapNotFound
			}
			value, ok := configMap.Data[key]
			if !ok {
				log.Info(fmt.Sprintf("Key %s in Configmap does not exist", key), "namespace", namespace, "configmap", name)
				return "", ErrKeyNotFound
			}
			return value, nil
		case user.Spec.PasswordFrom.SecretKeyRef != nil:
			s := user.Spec.PasswordFrom.SecretKeyRef
			namespace := user.Namespace
			name := s.Name
			key := s.Key
			secret := corev1.Secret{}
			if err := um.Reconciler.Get(um.Context, types.NamespacedName{Namespace: namespace, Name: name}, &secret); err != nil {
				log.Info(fmt.Sprintf("Read secret, error: %s", err), "namespace", namespace, "secret", name)
				return "", ErrSecretNotFound
			}
			valueBytes, ok := secret.Data[key]
			if !ok {
				log.Info(fmt.Sprintf("Key %s in secret does not exist", key), "namespace", namespace, "configmap", name)
				return "", ErrKeyNotFound
			}
			return string(valueBytes), nil
		}
	}
	return "", ErrMissingPassword
}
