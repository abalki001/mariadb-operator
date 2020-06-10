package mariadbcluster

import (
	"context"

	mariadbv1alpha1 "github.com/persistentsys/mariadb-operator/pkg/apis/mariadb/v1alpha1"
	"github.com/persistentsys/mariadb-operator/pkg/resource"
	"github.com/persistentsys/mariadb-operator/pkg/service"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileMariaDBCluster) ensureStatefulSet(request reconcile.Request,
	instance *mariadbv1alpha1.MariaDBCluster,
	dep *appsv1.StatefulSet,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.StatefulSet{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		log.Info("Creating a new StatefulSet Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			log.Error(err, "Failed to create new StatefulSet Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &reconcile.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		log.Error(err, "Failed to get StatefulSet Deployment")
		return &reconcile.Result{}, err
	}

	// Check for any updates for redeployment
	applyChange := false

	// Ensure image name is correct, update image if required
	image := instance.Spec.Image
	var currentImage string = ""

	if found.Spec.Template.Spec.Containers != nil {
		currentImage = found.Spec.Template.Spec.Containers[0].Image
	}

	if image != currentImage {
		dep.Spec.Template.Spec.Containers[0].Image = image
		applyChange = true
	}

	if applyChange {
		err = r.client.Update(context.TODO(), dep)
		if err != nil {
			log.Error(err, "Failed to update StatefulSet Deployment.", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return &reconcile.Result{}, err
		}
		log.Info("Updated StatefulSet Deployment image. ")
	}

	return nil, nil
}

func (r *ReconcileMariaDBCluster) ensureService(request reconcile.Request,
	instance *mariadbv1alpha1.MariaDBCluster,
	s *corev1.Service,
) (*reconcile.Result, error) {
	found := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get Service")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcileMariaDBCluster) ensureSecret(request reconcile.Request,
	instance *mariadbv1alpha1.MariaDBCluster,
	s *corev1.Secret,
) (*reconcile.Result, error) {
	found := &corev1.Secret{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the secret
		log.Info("Creating a new secret", "Secret.Namespace", s.Namespace, "Secret.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Secret", "Secret.Namespace", s.Namespace, "Secret.Name", s.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the secret not existing
		log.Error(err, "Failed to get Secret")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// ensurePV - Ensure that PV is present. If not, create one
func (r *ReconcileMariaDBCluster) ensurePV(request reconcile.Request,
	instance *mariadbv1alpha1.MariaDBCluster,
) (*reconcile.Result, error) {
	pvName := resource.GetMariadbClusterVolumeName(instance)
	_, err := service.FetchPVByName(pvName, r.client)

	if err != nil && errors.IsNotFound(err) {
		// Create Persistent Volume
		log.Info("Creating a new PV", "PV.Name", pvName)

		pv := resource.NewMariaDbClusterPV(instance, r.scheme)
		err := r.client.Create(context.TODO(), pv)
		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new PV", "PV.Name", pvName)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get PV")
		return &reconcile.Result{}, err
	}
	return nil, nil
}

// ensurePVC - Ensure that PVC is present. If not, create one
func (r *ReconcileMariaDBCluster) ensurePVC(request reconcile.Request,
	instance *mariadbv1alpha1.MariaDBCluster,
) (*reconcile.Result, error) {
	pvcName := resource.GetMariadbClusterVolumeClaimName(instance)
	_, err := service.FetchPVCByNameAndNS(pvcName, instance.Namespace, r.client)

	if err != nil && errors.IsNotFound(err) {
		// Create Persistent Volume Claim
		log.Info("Creating a new PVC", "PVC.Name", pvcName)

		pvc := resource.NewMariaDbClusterPVC(instance, r.scheme)
		err := r.client.Create(context.TODO(), pvc)
		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new PVC", "PV.Name", pvcName, "PVC.Namespace", instance.Namespace)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get PVC")
		return &reconcile.Result{}, err
	}
	return nil, nil
}
