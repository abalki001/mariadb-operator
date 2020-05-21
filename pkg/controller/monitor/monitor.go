package monitor

import (
	"context"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	grafanav1alpha1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	mariadbv1alpha1 "github.com/persistentsys/mariadb-operator/pkg/apis/mariadb/v1alpha1"
	"github.com/persistentsys/mariadb-operator/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const monitorPort = 9104
const monitorPortName = "monitor"
const monitorApp = "monitor-app"

func (r *ReconcileMonitor) monitorDeployment(v *mariadbv1alpha1.Monitor) *appsv1.Deployment {

	labels := utils.MonitorLabels(v, monitorApp)

	size := v.Spec.Size
	image := v.Spec.Image
	dataSourceName := v.Spec.DataSourceName

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      monitorDeploymentName(v),
			Namespace: v.Namespace,
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
						Image:           image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            monitorApp,
						Ports: []corev1.ContainerPort{{
							ContainerPort: monitorPort,
							Name:          monitorPortName,
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "DATA_SOURCE_NAME",
								Value: dataSourceName,
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.scheme)
	return dep
}

func (r *ReconcileMonitor) monitorService(v *mariadbv1alpha1.Monitor) *corev1.Service {
	labels := utils.MonitorLabels(v, monitorApp)

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      monitorServiceName(v),
			Namespace: v.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       monitorPort,
				TargetPort: intstr.FromInt(9104),
				Name:       monitorPortName,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	controllerutil.SetControllerReference(v, s, r.scheme)
	return s
}

func (r *ReconcileMonitor) monitorServiceMonitor(v *mariadbv1alpha1.Monitor) *monitoringv1.ServiceMonitor {
	labels := utils.ServiceMonitorLabels(v, monitorApp)

	s := &monitoringv1.ServiceMonitor{

		ObjectMeta: v12.ObjectMeta{
			Name:      monitorServiceMonitorName(v),
			Namespace: v.Namespace,
			Labels:    labels,
		},
		Spec: monitoringv1.ServiceMonitorSpec{

			Endpoints: []monitoringv1.Endpoint{{
				Path: "/metrics",
				Port: monitorPortName,
			}},
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"tier": monitorApp,
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, s, r.scheme)
	return s
}

func (r *ReconcileMonitor) monitorGrafanaDashboard(v *mariadbv1alpha1.Monitor) *grafanav1alpha1.GrafanaDashboard {

	labels := utils.ServiceMonitorLabels(v, monitorApp)

	s := &grafanav1alpha1.GrafanaDashboard{
		ObjectMeta: v12.ObjectMeta{
			Name:      "GrafanaDashboard",
			Namespace: v.Namespace,
			Labels:    labels,
		},
		Spec: grafanav1alpha1.GrafanaDashboardSpec{
			Json: DashboardJSON,
			Name: "mariadb.json",
			Plugins: []grafanav1alpha1.GrafanaPlugin{
				{
					Name:    "grafana-piechart-panel",
					Version: "1.5.0",
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, s, r.scheme)
	return s
}

func monitorDeploymentName(v *mariadbv1alpha1.Monitor) string {
	return v.Name + "-deployment"
}

func monitorServiceName(v *mariadbv1alpha1.Monitor) string {
	return v.Name + "-service"
}

func monitorServiceMonitorName(v *mariadbv1alpha1.Monitor) string {
	return v.Name + "-serviceMonitor"
}

func (r *ReconcileMonitor) updateMonitorStatus(v *mariadbv1alpha1.Monitor) error {
	err := r.client.Status().Update(context.TODO(), v)
	return err
}
