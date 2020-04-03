package mariadb

import (
	"context"
	"time"

	mariadbv1alpha1 "github.com/persistentsys/mariadb-operator/pkg/apis/mariadb/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const mariadbPort = 80
//const mariadbNodePort = 80
const mariadbImage = "mariadb/server:10.3"

func mariadbDeploymentName(v *mariadbv1alpha1.MariaDB) string {
	return v.Name + "-deployment"
}

func mariadbServiceName(v *mariadbv1alpha1.MariaDB) string {
	return v.Name + "-service"
}

func (r *ReconcileMariaDB) mariadbDeployment(v *mariadbv1alpha1.MariaDB) *appsv1.Deployment {
	labels := labels(v, "mariadb")
	size := v.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:		mariadbDeploymentName(v),
			Namespace: 	v.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:	mariadbImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:	"mariadb-service",
						Ports:	[]corev1.ContainerPort{{
							ContainerPort: 	mariadbPort,
							Name:			"mariadb",
						}},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.scheme)
	return dep
}

func (r *ReconcileMariaDB) mariadbService(v *mariadbv1alpha1.MariaDB) *corev1.Service {
	labels := labels(v, "mariadb")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:		mariadbServiceName(v),
			Namespace: 	v.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol: corev1.ProtocolTCP,
				Port: mariadbPort,
				TargetPort: intstr.FromInt(mariadbPort),
				NodePort: 30685,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	controllerutil.SetControllerReference(v, s, r.scheme)
	return s
}

func (r *ReconcileMariaDB) updateMariadbStatus(v *mariadbv1alpha1.MariaDB) (error) {
	//v.Status.BackendImage = mariadbImage
	err := r.client.Status().Update(context.TODO(), v)
	return err
}

func (r *ReconcileMariaDB) handleMariadbChanges(v *mariadbv1alpha1.MariaDB) (*reconcile.Result, error) {
	found := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      mariadbDeploymentName(v),
		Namespace: v.Namespace,
	}, found)
	if err != nil {
		// The deployment may not have been created yet, so requeue
		return &reconcile.Result{RequeueAfter:5 * time.Second}, err
	}

	size := v.Spec.Size

	if size != *found.Spec.Replicas {
		found.Spec.Replicas = &size
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			log.Error(err, "Failed to update Deployment.", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return &reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return &reconcile.Result{Requeue: true}, nil
	}

	return nil, nil
}