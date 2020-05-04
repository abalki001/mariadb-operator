package resource

import (
	"github.com/persistentsys/mariadb-operator/pkg/apis/mariadb/v1alpha1"
	"github.com/persistentsys/mariadb-operator/pkg/utils"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const pvStorageName = "mariadb-bkp-pv-storage"
const pvClaimName = "mariadb-pv-claim"

// const pvBkpBackupName = "mariadb-backup-storage"
// const pvBkpClaimName = "mariadb-bkp-pv-claim"

//Returns the NewBackupCronJob object for the Database Backup
func NewBackupCronJob(bkp *v1alpha1.Backup, db *v1alpha1.MariaDB, scheme *runtime.Scheme) *v1beta1.CronJob {
	cron := &v1beta1.CronJob{
		ObjectMeta: v1.ObjectMeta{
			Name:      bkp.Name,
			Namespace: bkp.Namespace,
			Labels:    utils.Labels(db, "mariadb"),
		},
		Spec: v1beta1.CronJobSpec{
			Schedule: bkp.Spec.Schedule,
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							ServiceAccountName: "mariadb-operator",
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
							Containers: []corev1.Container{
								{
									Name:    bkp.Name,
									Image:   db.Spec.Image,
									Command: []string{"/bin/sh", "-c"},
									Args:    []string{"echo 'Starting DB Backup' && mysqldump  --lock-tables --all-databases > /var/lib/mysql/backup.sql && echo 'Completed DB Backup'"},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      pvStorageName,
											MountPath: "/var/lib/mysql",
										},
									},
									Env: []corev1.EnvVar{
										{
											Name:  "MYSQL_PWD",
											Value: db.Spec.Rootpwd,
										},
										{
											Name:  "USER",
											Value: "root",
										},
									},
								},
							},
							RestartPolicy: corev1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(bkp, cron, scheme)
	return cron
}
