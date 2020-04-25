package controller

import (
	"github.com/blaqkube/mysql-operator/pkg/controller/instance"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, instance.Add)
}
