package backup

import (
	"context"
	"fmt"

	"github.com/persistentsys/mariadb-operator/pkg/apis/mariadb/v1alpha1"
	"github.com/persistentsys/mariadb-operator/pkg/resource"
	"github.com/persistentsys/mariadb-operator/pkg/service"
)

// Set in the ReconcileBackup the Pod database created by Database
// NOTE: This data is required in order to create the secrets which will access the database container to do the backup
func (r *ReconcileBackup) getDatabasePod(bkp *v1alpha1.Backup, db *v1alpha1.MariaDB) error {
	dbPod, err := service.FetchDatabasePod(bkp, db, r.client)
	if err != nil || dbPod == nil {
		r.dbPod = nil
		err := fmt.Errorf("Unable to find the Database Pod")
		return err
	}
	r.dbPod = dbPod
	return nil
}

// Set in the ReconcileBackup the service database created by Database
// NOTE: This data is required in order to create the secrets which will access the database container to do the backup
func (r *ReconcileBackup) getDatabaseService(bkp *v1alpha1.Backup, db *v1alpha1.MariaDB) error {
	dbService, err := service.FetchDatabaseService(bkp, db, r.client)
	if err != nil || dbService == nil {
		r.dbService = nil
		err := fmt.Errorf("Unable to find the Database Service")
		return err
	}
	r.dbService = dbService
	return nil
}

// Check if the cronJob is created, if not create one
func (r *ReconcileBackup) createCronJob(bkp *v1alpha1.Backup, db *v1alpha1.MariaDB) error {
	if _, err := service.FetchCronJob(bkp.Name, bkp.Namespace, r.client); err != nil {
		if err := r.client.Create(context.TODO(), resource.NewBackupCronJob(bkp, db, r.scheme)); err != nil {
			return err
		}
	}
	return nil
}
