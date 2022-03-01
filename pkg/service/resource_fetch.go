package service

import (
	"context"

	"github.com/persistentsys/mariadb-operator/pkg/apis/mariadb/v1alpha1"
	"github.com/persistentsys/mariadb-operator/pkg/utils"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var rfLog = logf.Log.WithName("resource_fetch")

// FetchDatabaseCR fetches CR of MariDB
// Request object not found, could have been deleted after reconcile request.
// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
func FetchDatabaseCR(name, namespace string, client client.Client) (*v1alpha1.MariaDB, error) {
	rfLog.Info("Fetching Database CR ...")
	db := &v1alpha1.MariaDB{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, db)
	return db, err
}

// FetchBackupCR fetches CR of Maria DB Backup object
func FetchBackupCR(name, namespace string, client client.Client) (*v1alpha1.Backup, error) {
	rfLog.Info("Fetching Backup CR ...")
	bkp := &v1alpha1.Backup{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, bkp)
	return bkp, err
}

// FetchDatabasePod search in the cluster for 1 Pod managed by the Database Controller
func FetchDatabasePod(bkp *v1alpha1.Backup, db *v1alpha1.MariaDB, client client.Client) (*corev1.Pod, error) {
	rfLog.Info("Fetching Database Pod ...")
	listOps := buildDatabaseCriteria(bkp, db)
	dbPodList := &corev1.PodList{}
	err := client.List(context.TODO(), dbPodList, listOps)
	if err != nil {
		return nil, err
	}

	if len(dbPodList.Items) == 0 {
		return nil, err
	}

	pod := dbPodList.Items[0]
	return &pod, nil
}

// FetchDatabaseService search in the cluster for 1 Service managed by the Database Controller
func FetchDatabaseService(bkp *v1alpha1.Backup, db *v1alpha1.MariaDB, client client.Client) (*corev1.Service, error) {
	rfLog.Info("Fetching Database Service ...")
	listOps := buildDatabaseCriteria(bkp, db)
	dbServiceList := &corev1.ServiceList{}
	err := client.List(context.TODO(), dbServiceList, listOps)
	if err != nil {
		return nil, err
	}

	if len(dbServiceList.Items) == 0 {
		return nil, err
	}

	srv := dbServiceList.Items[0]
	return &srv, nil
}

// FetchDatabaseBackupService search in the cluster for 1 Service managed by the Backup Controller
func FetchDatabaseBackupService(bkp *v1alpha1.Backup, db *v1alpha1.MariaDB, client client.Client) (*corev1.Service, error) {
	rfLog.Info("Fetching Database Backup Service ...")
	listOps := buildDatabaseBackupCriteria(bkp, db)
	bkpServiceList := &corev1.ServiceList{}
	err := client.List(context.TODO(), bkpServiceList, listOps)
	if err != nil {
		return nil, err
	}

	if len(bkpServiceList.Items) == 0 {
		return nil, err
	}

	srv := bkpServiceList.Items[0]
	return &srv, nil
}

// FetchCronJob returns the CronJob resource with the name in the namespace
func FetchCronJob(name, namespace string, client client.Client) (*v1beta1.CronJob, error) {
	rfLog.Info("Fetching CronJob ...")
	cronJob := &v1beta1.CronJob{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, cronJob)
	return cronJob, err
}

// FetchPVByName search in the cluster for PV managed by the Backup Controller
func FetchPVByName(name string, client client.Client) (*corev1.PersistentVolume, error) {
	reqLogger := rfLog.WithValues("PV Name", name)
	reqLogger.Info("Fetching Persistent Volume")

	pv := &corev1.PersistentVolume{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name}, pv)
	return pv, err
}

// FetchPVCByNameAndNS search in the cluster for PVC managed by the Backup Controller
func FetchPVCByNameAndNS(name, namespace string, client client.Client) (*corev1.PersistentVolumeClaim, error) {
	reqLogger := rfLog.WithValues("PVC Name", name, "PVC Namespace", namespace)
	reqLogger.Info("Fetching Persistent Volume Claim")

	pvc := &corev1.PersistentVolumeClaim{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, pvc)
	return pvc, err
}

// buildDatabaseCreteria returns client.ListOptions required to fetch the secondary resource created by
func buildDatabaseCriteria(bkp *v1alpha1.Backup, db *v1alpha1.MariaDB) *client.ListOptions {
	labelSelector := labels.SelectorFromSet(utils.Labels(db, "mariadb"))
	listOps := &client.ListOptions{Namespace: db.Namespace, LabelSelector: labelSelector}
	return listOps
}

// buildDatabaseCreteria returns client.ListOptions required to fetch the secondary resource created by
func buildDatabaseBackupCriteria(bkp *v1alpha1.Backup, db *v1alpha1.MariaDB) *client.ListOptions {
	labelSelector := labels.SelectorFromSet(utils.Labels(db, "mariadb-backup"))
	listOps := &client.ListOptions{Namespace: bkp.Namespace, LabelSelector: labelSelector}
	return listOps
}
