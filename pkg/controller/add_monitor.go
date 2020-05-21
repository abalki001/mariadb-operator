package controller

import (
	"github.com/persistentsys/mariadb-operator/pkg/controller/monitor"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, monitor.Add)
}
