package utils

import (
	"github.com/persistentsys/mariadb-operator/pkg/apis/mariadb/v1alpha1"
)

func Labels(v *v1alpha1.MariaDB, tier string) map[string]string {
	return map[string]string{
		"app":        "MariaDB",
		"MariaDB_cr": v.Name,
		"tier":       tier,
	}
}

func MariaDBBkpLabels(v *v1alpha1.Backup, tier string) map[string]string {
	return map[string]string{
		"app":        "MariaDB-Backup",
		"MariaDB_cr": v.Name,
		"tier":       tier,
	}
}

func MonitorLabels(v *v1alpha1.Monitor, tier string) map[string]string {
	return map[string]string{
		"app":        "MariaDB-Monitor",
		"Monitor_cr": v.Name,
		"tier":       tier,
	}
}

func ServiceMonitorLabels(v *v1alpha1.Monitor, tier string) map[string]string {
	return map[string]string{
		"app":        "ServiceMonitor",
		"Monitor_cr": v.Name,
		"tier":       tier,
	}
}
