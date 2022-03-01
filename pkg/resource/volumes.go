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
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var volLog = logf.Log.WithName("resource_volumes")

// GetMariadbVolumeName - return name of PV used in MariaDB
func GetMariadbVolumeName(v *v1alpha1.MariaDB) string {
	return v.Name + "-" + v.Namespace + "-pv"
}

// GetMariadbVolumeClaimName - return name of PVC used in MariaDB
func GetMariadbVolumeClaimName(v *v1alpha1.MariaDB) string {
	return v.Name + "-pv-claim"
}

// GetMariadbBkpVolumeName - return name of PV used in DB Backup
func GetMariadbBkpVolumeName(bkp *v1alpha1.Backup) string {
	return bkp.Name + "-" + bkp.Namespace + "-pv"
}

// GetMariadbBkpVolumeClaimName - return name of PVC used in DB Backup
func GetMariadbBkpVolumeClaimName(bkp *v1alpha1.Backup) string {
	return bkp.Name + "-pv-claim"
}

// NewDbBackupPV Create a new PV object for Database Backup
func NewDbBackupPV(bkp *v1alpha1.Backup, v *v1alpha1.MariaDB, scheme *runtime.Scheme) *corev1.PersistentVolume {
	volLog.Info("Creating new PV for Database Backup")
	labels := utils.MariaDBBkpLabels(bkp, "mariadb-backup")
	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: GetMariadbBkpVolumeName(bkp),
			// Namespace: v.Namespace,
			Labels: labels,
		},
		Spec: corev1.PersistentVolumeSpec{
			StorageClassName: "manual",
			Capacity: corev1.ResourceList{
				corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(bkp.Spec.BackupSize),
			},
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: bkp.Spec.BackupPath},
			},
		},
	}

	volLog.Info("PV created for Database Backup ")
	controllerutil.SetControllerReference(bkp, pv, scheme)
	return pv
}

// NewDbBackupPVC Create a new PV Claim object for Database Backup
func NewDbBackupPVC(bkp *v1alpha1.Backup, v *v1alpha1.MariaDB, scheme *runtime.Scheme) *corev1.PersistentVolumeClaim {
	volLog.Info("Creating new PVC for Database Backup")
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
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(bkp.Spec.BackupSize),
				},
			},
			VolumeName: GetMariadbBkpVolumeName(bkp),
		},
	}

	volLog.Info("PVC created for Database Backup ")
	controllerutil.SetControllerReference(bkp, pvc, scheme)
	return pvc
}

// NewMariaDbPV Create a new PV object for MariaDB
func NewMariaDbPV(v *v1alpha1.MariaDB, scheme *runtime.Scheme) *corev1.PersistentVolume {
	volLog.Info("Creating new PV for MariaDB")
	labels := utils.Labels(v, "mariadb")
	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: GetMariadbVolumeName(v),
			// Namespace: v.Namespace,
			Labels: labels,
		},
		Spec: corev1.PersistentVolumeSpec{
			StorageClassName: "manual",
			Capacity: corev1.ResourceList{
				corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(v.Spec.DataStorageSize),
			},
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: v.Spec.DataStoragePath},
			},
		},
	}

	volLog.Info("PV created for MariaDB ")
	controllerutil.SetControllerReference(v, pv, scheme)
	return pv
}

// NewMariaDbPVC Create a new PV Claim object for MariaDB
func NewMariaDbPVC(v *v1alpha1.MariaDB, scheme *runtime.Scheme) *corev1.PersistentVolumeClaim {
	volLog.Info("Creating new PVC for MariaDB")
	labels := utils.Labels(v, "mariadb")
	storageClassName := "manual"
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetMariadbVolumeClaimName(v),
			Namespace: v.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(v.Spec.DataStorageSize),
				},
			},
			VolumeName: GetMariadbVolumeName(v),
		},
	}

	volLog.Info("PVC created for MariaDB ")
	controllerutil.SetControllerReference(v, pvc, scheme)
	return pvc
}
