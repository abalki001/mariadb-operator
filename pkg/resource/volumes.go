package resource

import (
	"github.com/persistentsys/mariadb-operator/pkg/apis/mariadb/v1alpha1"
	"github.com/persistentsys/mariadb-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//const dbBakupServicePort = 3306
//const dbBakupServiceTargetPort = 3306

// GetMariadbBkpVolumeName - return name of PV used in DB Backup
func GetMariadbBkpVolumeName(bkp *v1alpha1.Backup) string {
	return bkp.Name + "-pv-volume-test"
}

// GetMariadbBkpVolumeClaimName - return name of PVC used in DB Backup
func GetMariadbBkpVolumeClaimName(bkp *v1alpha1.Backup) string {
	return bkp.Name + "-pv-claim-test"
}

// NewDbBackupPV Create a new PV object for Database Backup
func NewDbBackupPV(bkp *v1alpha1.Backup, v *v1alpha1.MariaDB, scheme *runtime.Scheme) *corev1.PersistentVolume {

	labels := utils.MariaDBBkpLabels(bkp, "mariadb-backup")
	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetMariadbBkpVolumeName(bkp),
			Namespace: v.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeSpec{
			StorageClassName: "manual",
			Capacity: corev1.ResourceList{
				corev1.ResourceName(corev1.ResourceStorage): resource.MustParse("1Gi"),
			},
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/mnt/backup"},
			},
		},
	}

	controllerutil.SetControllerReference(bkp, pv, scheme)
	return pv
}

// NewDbBackupPVC Create a new PV Claim object for Database Backup
func NewDbBackupPVC(bkp *v1alpha1.Backup, v *v1alpha1.MariaDB, scheme *runtime.Scheme) *corev1.PersistentVolumeClaim {

	labels := utils.MariaDBBkpLabels(bkp, "mariadb-backup")
	storageClassName := "manual"
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetMariadbBkpVolumeClaimName(bkp),
			Namespace: v.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse("1Gi"),
				},
			},
			VolumeName: GetMariadbBkpVolumeName(bkp),
		},
	}

	controllerutil.SetControllerReference(bkp, pvc, scheme)
	return pvc
}
