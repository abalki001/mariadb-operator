package mariadbcluster

import (
	"context"

	mariadbv1alpha1 "github.com/persistentsys/mariadb-operator/pkg/apis/mariadb/v1alpha1"
	"github.com/persistentsys/mariadb-operator/pkg/resource"
	"github.com/persistentsys/mariadb-operator/pkg/service"
	"github.com/persistentsys/mariadb-operator/pkg/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var clusterLog = logf.Log.WithName("mariadbcluster")

const mariadbPort = 3306
const pvStorageName = "mariadb-pv-storage"

func mariadbClusterStatefulStateName(v *mariadbv1alpha1.MariaDBCluster) string {
	return v.Name + "-statefulstate"
}

func mariadbServiceName(v *mariadbv1alpha1.MariaDBCluster) string {
	return v.Name + "-service"
}

func mariadbLBServiceName(v *mariadbv1alpha1.MariaDBCluster) string {
	return v.Namespace + "-lb-service"
}

func mariadbClusterAuthName() string {
	return "mariadb-cluster-auth"
}

func (r *ReconcileMariaDBCluster) mariadbClusterStatefulSet(v *mariadbv1alpha1.MariaDBCluster) *appsv1.StatefulSet {
	pvClaimName := resource.GetMariadbClusterVolumeClaimName(v)
	labels := utils.MariaDBClusterLabels(v, "mariadb-cluster")
	image := v.Spec.Image

	dbname := v.Spec.Database
	rootpwd := v.Spec.Rootpwd

	userSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: mariadbClusterAuthName()},
			Key:                  "username",
		},
	}

	passwordSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: mariadbClusterAuthName()},
			Key:                  "password",
		},
	}

	gcommUrl := ""
	// Get Headless Service list
	dbServiceList, err := service.FetchClusterHeadlessServiceList(v, r.client)
	if err != nil || dbServiceList == nil {
		// Service not found, this is first pod
		//gcommUrl := ""
	} else {
		// Print Service list
		reqLogger := clusterLog.WithValues("Service_List_Size", len(dbServiceList.Items))
		reqLogger.Info("MariaDB CLuster Headless Service List")

		// TODO: Get list of names from service for gcomm_url
		// serviceItems := dbServiceList.Items

		//svcItem := &corev1.Service{}
		for i := 0; i < len(dbServiceList.Items); i++ {
			svcItem := dbServiceList.Items[i]
			svcName := svcItem.ObjectMeta.Name

			if svcName != mariadbClusterStatefulStateName(v) {
				gcommUrl = gcommUrl + svcName + "." + svcItem.ObjectMeta.Namespace + ","
			}
		}
	}

	// Create StatefulSet Deployment object
	dep := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mariadbClusterStatefulStateName(v),
			Namespace: v.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: pvStorageName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: pvClaimName,
								},
							},
						},
					},
					Containers: []corev1.Container{{
						Image:           image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "mariadb-cluster-service",
						Ports: []corev1.ContainerPort{{
							ContainerPort: mariadbPort,
							Name:          "mariadb-port",
						}},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      pvStorageName,
								MountPath: "/var/lib/mysql",
							},
						},
						Args: []string{
							"--wsrep-on=ON",
							"--wsrep-cluster-address=gcomm://" + gcommUrl,
							"--wsrep-provider=/usr/lib/galera/libgalera_smm.so",
							"--binlog-format=row",
							"--default-storage-engine=InnoDB",
							"--innodb-autoinc-lock_mode=2",
							"--bind-address=0.0.0.0",
							"--wsrep-cluster-name='galera_cluster'",
						},
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_ROOT_PASSWORD",
								Value: rootpwd,
							},
							{
								Name:  "MYSQL_DATABASE",
								Value: dbname,
							},
							{
								Name:      "MYSQL_USER",
								ValueFrom: userSecret,
							},
							{
								Name:      "MYSQL_PASSWORD",
								ValueFrom: passwordSecret,
							},
						},
					}},
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": v.Spec.Cluster.NodeName,
					},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.scheme)
	return dep
}

func (r *ReconcileMariaDBCluster) mariadbClusterHeadlessService(v *mariadbv1alpha1.MariaDBCluster) *corev1.Service {
	labels := utils.MariaDBClusterHeadlessServiceLabels(v, "mariadb-cluster")
	selectorLabels := utils.MariaDBClusterHeadlessServiceSelector(v, "mariadb-cluster")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mariadbServiceName(v),
			Namespace: v.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector:  selectorLabels,
			ClusterIP: "None",
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       mariadbPort,
				TargetPort: intstr.FromInt(3306),
				Name:       "mariadb-port",
			},
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       4444,
					TargetPort: intstr.FromInt(4444),
					Name:       "sst-port",
				},
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       4567,
					TargetPort: intstr.FromInt(4567),
					Name:       "galera-replication-port",
				},
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       4568,
					TargetPort: intstr.FromInt(4568),
					Name:       "ist-port",
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	controllerutil.SetControllerReference(v, s, r.scheme)
	return s
}

// mariadbClusterLBService Create Load Balancer Service to expose DB outside Cluster
func (r *ReconcileMariaDBCluster) mariadbClusterLBService(v *mariadbv1alpha1.MariaDBCluster) *corev1.Service {
	labels := utils.MariaDBClusterLBServiceLabels(v, "mariadb-cluster")
	selectorLabels := utils.MariaDBClusterLBServiceSelector(v, "mariadb-cluster")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mariadbLBServiceName(v),
			Namespace: v.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       mariadbPort,
				TargetPort: intstr.FromInt(3306),
				NodePort:   v.Spec.Port,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	controllerutil.SetControllerReference(v, s, r.scheme)
	return s
}

func (r *ReconcileMariaDBCluster) updateMariadbStatus(v *mariadbv1alpha1.MariaDBCluster) error {
	//v.Status.BackendImage = mariadbImage
	err := r.client.Status().Update(context.TODO(), v)
	return err
}

func (r *ReconcileMariaDBCluster) mariadbAuthSecret(v *mariadbv1alpha1.MariaDBCluster) *corev1.Secret {

	username := v.Spec.Username
	password := v.Spec.Password

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mariadbClusterAuthName(),
			Namespace: v.Namespace,
		},
		Type: "Opaque",
		Data: map[string][]byte{
			"username": []byte(username),
			"password": []byte(password),
		},
	}
	controllerutil.SetControllerReference(v, secret, r.scheme)
	return secret
}
