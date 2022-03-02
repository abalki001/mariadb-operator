package mariadb

import (
	"context"

	mariadbv1alpha1 "github.com/persistentsys/mariadb-operator/pkg/apis/mariadb/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	//"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_mariadb")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new MariaDB Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMariaDB{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mariadb-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MariaDB
	err = c.Watch(&source.Kind{Type: &mariadbv1alpha1.MariaDB{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner MariaDB
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mariadbv1alpha1.MariaDB{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mariadbv1alpha1.MariaDB{},
	})

	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileMariaDB implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileMariaDB{}

// ReconcileMariaDB reconciles a MariaDB object
type ReconcileMariaDB struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a MariaDB object and makes changes based on the state read
// and what is in the MariaDB.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMariaDB) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MariaDB")

	// Fetch the MariaDB instance
	instance := &mariadbv1alpha1.MariaDB{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	var result *reconcile.Result

	result, err = r.ensureSecret(request, instance, r.mariadbAuthSecret(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureDeployment(request, instance, r.mariadbDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(request, instance, r.mariadbService(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensurePV(request, instance)
	if result != nil {
		return *result, err
	}

	result, err = r.ensurePVC(request, instance)
	if result != nil {
		return *result, err
	}

	err = r.updateMariadbStatus(instance)
	if err != nil {
		// Requeue the request if the status could not be updated
		return reconcile.Result{}, err
	}

	// Everything went fine, don't requeue
	return reconcile.Result{}, nil
}
