package utils

import (
	"github.com/persistentsys/mariadb-operator/pkg/apis/mariadb/v1alpha1"
)

var defaultBackupConfig = NewDefaultBackupConfig()

// AddBackupMandatorySpecs will add the specs which are mandatory for Backup CR in the case them
// not be applied
func AddBackupMandatorySpecs(bkp *v1alpha1.Backup) {

	/*
	 Backup Container
	*/

	if bkp.Spec.Schedule == "" {
		bkp.Spec.Schedule = defaultBackupConfig.Schedule
	}

	if bkp.Spec.BackupPath == "" {
		bkp.Spec.BackupPath = defaultBackupConfig.BackupPath
	}

}
