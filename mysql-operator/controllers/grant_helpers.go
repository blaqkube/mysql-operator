package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/blaqkube/mysql-operator/mysql-operator/agent"
	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	"github.com/prometheus/common/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	maxGrantConditions = 10
)

// GrantManager provides methods to manage grant subcomponents
type GrantManager struct {
	Context     context.Context
	Reconciler  *GrantReconciler
	TimeManager *TimeManager
}

func (gm *GrantManager) setGrantCondition(grant *mysqlv1alpha1.Grant, condition metav1.Condition) (ctrl.Result, error) {
	if condition.Reason == grant.Status.Reason {
		c := len(grant.Status.Conditions) - 1
		d := gm.TimeManager.Next(grant.Status.Conditions[c].LastTransitionTime.Time)
		if condition.Reason != mysqlv1alpha1.DatabaseSucceeded {
			return ctrl.Result{Requeue: true, RequeueAfter: d}, nil
		}
		return ctrl.Result{}, nil
	}
	grant.Status.Ready = condition.Status
	grant.Status.Reason = condition.Reason
	grant.Status.Message = condition.Message
	conditions := append(grant.Status.Conditions, condition)
	if len(conditions) > maxGrantConditions {
		conditions = conditions[1:]
	}
	grant.Status.Conditions = conditions
	log := gm.Reconciler.Log.WithValues("namespace", grant.Namespace, "grant", grant.Name)
	log.Info("Updating grant with new Status", "Reason", condition.Reason, "Message", condition.Message)
	if err := gm.Reconciler.Status().Update(gm.Context, grant); err != nil {
		log.Error(err, "Unable to update grant")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// CreateGrant is the script that creates a grant
func (gm *GrantManager) CreateGrant(grant *mysqlv1alpha1.Grant) error {
	_ = gm.Reconciler.Log.WithValues("namespace", grant.Namespace, "grant", grant.Name)
	a := &APIReconciler{
		Client: gm.Reconciler.Client,
		Log:    gm.Reconciler.Log,
	}
	user, err := a.GetUser(gm.Context, types.NamespacedName{Namespace: grant.Namespace, Name: grant.Spec.User})
	if err != nil {
		return err
	}
	database, err := a.GetDatabase(gm.Context, types.NamespacedName{Namespace: grant.Namespace, Name: grant.Spec.Database})
	if err != nil {
		return err
	}
	if user.Spec.Instance != database.Spec.Instance {
		return ErrUserDatabaseMismatch
	}
	api, err := a.GetAPI(
		gm.Context,
		types.NamespacedName{
			Name:      user.Spec.Instance,
			Namespace: grant.Namespace,
		},
	)
	if err != nil {
		return err
	}
	_, response, err := api.MysqlApi.CreateGrantForUserDatabase(
		gm.Context,
		database.Spec.Name,
		user.Spec.Username,
		agent.Grant{
			AccessMode: string(grant.Spec.AccessMode),
		},
		nil)
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
