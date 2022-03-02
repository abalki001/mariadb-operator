package resource

import (
	"github.com/persistentsys/mariadb-operator/pkg/apis/mariadb/v1alpha1"
	"github.com/persistentsys/mariadb-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const dbBakupServicePort = 3306
const dbBakupServiceTargetPort = 3306

func mariadbBkpServiceName(bkp *v1alpha1.Backup) string {
	return bkp.Name + "-service"
}

// NewDbBackupService Create a new service object for Database Backup
func NewDbBackupService(bkp *v1alpha1.Backup, v *v1alpha1.MariaDB, scheme *runtime.Scheme) *corev1.Service {
	labels := utils.Labels(v, "mariadb-backup")
	selectorLabels := utils.Labels(v, "mariadb")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mariadbBkpServiceName(bkp),
			Namespace: v.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       dbBakupServicePort,
				TargetPort: intstr.FromInt(dbBakupServiceTargetPort),
			}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	controllerutil.SetControllerReference(v, s, scheme)
	return s
}
