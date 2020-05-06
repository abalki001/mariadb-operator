package backup

import (
	mariadbv1alpha1 "github.com/persistentsys/mariadb-operator/pkg/apis/mariadb/v1alpha1"
	"github.com/persistentsys/mariadb-operator/pkg/service"
	"github.com/persistentsys/mariadb-operator/pkg/utils"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_backup")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Backup Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileBackup{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("backup-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Backup
	err = c.Watch(&source.Kind{Type: &mariadbv1alpha1.Backup{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch CronJob resource controlled and created by it
	err = c.Watch(&source.Kind{Type: &v1beta1.CronJob{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mariadbv1alpha1.Backup{},
	})
	if err != nil {
		return err
	}

	// Watch Service resource managed by the Database
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mariadbv1alpha1.MariaDB{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileBackup implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileBackup{}

// ReconcileBackup reconciles a Backup object
type ReconcileBackup struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client    client.Client
	scheme    *runtime.Scheme
	dbPod     *corev1.Pod
	dbService *corev1.Service
}

// Reconcile reads that state of the cluster for a Backup object and makes changes based on the state read
// and what is in the Backup.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileBackup) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Backup")

	// Fetch the Backup instance
	bkp, err := service.FetchBackupCR(request.Name, request.Namespace, r.client)
	//instance := &mariadbv1alpha1.Backup{}
	//err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Backup resource not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Backup.")
		return reconcile.Result{}, err
	}

	// Add const values for mandatory specs
	log.Info("Adding backup mandatory specs")
	utils.AddBackupMandatorySpecs(bkp)

	// Create mandatory objects for the Backup
	if err := r.createResources(bkp, request); err != nil {
		log.Error(err, "Failed to create and update the secondary resource required for the Backup CR")
		return reconcile.Result{}, err
	}

	log.Info("Stop Reconciling Backup ...")
	return reconcile.Result{}, nil
}

//createResources will create and update the secondary resource which are required
//   in order to make works successfully the primary resource(CR)
func (r *ReconcileBackup) createResources(bkp *mariadbv1alpha1.Backup, request reconcile.Request) error {
	log.Info("Creating secondary Backup resources ...")

	// Check if the database instance was created
	db, err := service.FetchDatabaseCR("mariadb", request.Namespace, r.client)
	if err != nil {
		log.Error(err, "Failed to fetch Database instance/cr")
		return err
	}

	// Get the Database Pod created by the Database Controller
	if err := r.getDatabasePod(bkp, db); err != nil {
		log.Error(err, "Failed to get a Database pod")
		return err
	}

	// Get the Database Service created by the Database Controller
	if err := r.getDatabaseService(bkp, db); err != nil {
		log.Error(err, "Failed to get a Database service")
		return err
	}

	// Get the Database Backup Service created by the Backup Controller
	if err := r.getDatabaseBackupService(bkp, db); err != nil {
		log.Error(err, "Failed to get a Database Backup service")
		return err
	}

	// Check if the cronJob is created, if not create one
	if err := r.createCronJob(bkp, db); err != nil {
		log.Error(err, "Failed to create the CronJob")
		return err
	}

	return nil
}
