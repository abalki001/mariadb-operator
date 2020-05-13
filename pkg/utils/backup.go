package utils

const (
	schedule   = "0 0 * * *"
	backupPath = "/mnt/backup"
)

type DefaultBackupConfig struct {
	Schedule       string `json:"schedule"`
	BackupPath string `json:"backupPath"`
}

func NewDefaultBackupConfig() *DefaultBackupConfig {
	return &DefaultBackupConfig{
		Schedule:   schedule,
		BackupPath: backupPath,
	}
}
