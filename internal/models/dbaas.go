package models

import "time"

type DBaaSRecord struct {
	Domain        string    `bson:"domain" json:"domain"`
	Project       string    `bson:"project" json:"project"`
	Cluster       string    `bson:"cluster" json:"cluster"`
	Status        string    `bson:"status" json:"status"`
	KubeconfigB64 string    `bson:"kubeconfig_b64" json:"kubeconfig_b64"`
	LicenseB64    string    `bson:"license_b64" json:"license_b64"`
	BackupEnabled bool      `bson:"backup_enabled,omitempty" json:"backup_enabled"`
	UpdatedAt     time.Time `bson:"updatedAt" json:"updatedAt"`
}
