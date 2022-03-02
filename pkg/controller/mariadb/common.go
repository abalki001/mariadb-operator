package mariadb

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

func (r *ReconcileMariaDB) ensureDeployment(request reconcile.Request,
	instance *mariadbv1alpha1.MariaDB,
	dep *appsv1.Deployment,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &reconcile.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		log.Error(err, "Failed to get Deployment")
		return &reconcile.Result{}, err
	}

	// Check for any updates for redeployment
	applyChange := false

	// Ensure the deployment size is same as the spec
	size := instance.Spec.Size
	if *dep.Spec.Replicas != size {
		dep.Spec.Replicas = &size
		applyChange = true
	}

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
			log.Error(err, "Failed to update Deployment.", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return &reconcile.Result{}, err
		}
		log.Info("Updated Deployment image. ")
	}

	return nil, nil
}

func (r *ReconcileMariaDB) ensureService(request reconcile.Request,
	instance *mariadbv1alpha1.MariaDB,
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

func (r *ReconcileMariaDB) ensureSecret(request reconcile.Request,
	instance *mariadbv1alpha1.MariaDB,
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
func (r *ReconcileMariaDB) ensurePV(request reconcile.Request,
	instance *mariadbv1alpha1.MariaDB,
) (*reconcile.Result, error) {
	pvName := resource.GetMariadbVolumeName(instance)
	_, err := service.FetchPVByName(pvName, r.client)

	if err != nil && errors.IsNotFound(err) {
		// Create Persistent Volume
		log.Info("Creating a new PV", "PV.Name", pvName)

		pv := resource.NewMariaDbPV(instance, r.scheme)
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
func (r *ReconcileMariaDB) ensurePVC(request reconcile.Request,
	instance *mariadbv1alpha1.MariaDB,
) (*reconcile.Result, error) {
	pvcName := resource.GetMariadbVolumeClaimName(instance)
	_, err := service.FetchPVCByNameAndNS(pvcName, instance.Namespace, r.client)

	if err != nil && errors.IsNotFound(err) {
		// Create Persistent Volume Claim
		log.Info("Creating a new PVC", "PVC.Name", pvcName)

		pvc := resource.NewMariaDbPVC(instance, r.scheme)
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
