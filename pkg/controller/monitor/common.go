package monitor

import (
	"context"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	grafanav1alpha1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	mariadbv1alpha1 "github.com/persistentsys/mariadb-operator/pkg/apis/mariadb/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileMonitor) ensureDeployment(request reconcile.Request,
	instance *mariadbv1alpha1.Monitor,
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

func (r *ReconcileMonitor) ensureService(request reconcile.Request,
	instance *mariadbv1alpha1.Monitor,
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

func (r *ReconcileMonitor) ensureServiceMonitor(request reconcile.Request,
	instance *mariadbv1alpha1.Monitor,
	s *monitoringv1.ServiceMonitor,
) (*reconcile.Result, error) {
	found := &monitoringv1.ServiceMonitor{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new ServiceMonitor", "ServiceMonitor.Namespace", s.Namespace, "ServiceMonitor.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new ServiceMonitor", "ServiceMonitor.Namespace", s.Namespace, "ServiceMonitor.Name", s.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get ServiceMonitor")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcileMonitor) ensureGrafanaDashboard(request reconcile.Request,
	instance *mariadbv1alpha1.Monitor,
	s *grafanav1alpha1.GrafanaDashboard,
) (*reconcile.Result, error) {

	found := &grafanav1alpha1.GrafanaDashboard{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the GrafanaDashboard
		log.Info("Creating a new GrafanaDashboard", "GrafanaDashboard.Namespace", s.Namespace, "GrafanaDashboard.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new GrafanaDashboard", "GrafanaDashboard.Namespace", s.Namespace, "GrafanaDashboard.Name", s.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get GrafanaDashboard")
		return &reconcile.Result{}, err
	}

	return nil, nil
}
