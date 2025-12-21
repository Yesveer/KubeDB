// package models

// import "time"

// type DatabaseRecord struct {
// 	Domain  string `bson:"domain" json:"domain"`
// 	Project string `bson:"project" json:"project"`
// 	Cluster string `bson:"cluster" json:"cluster"`

// 	DBType    string `bson:"dbType" json:"dbType"`
// 	Name      string `bson:"name" json:"name"`
// 	Namespace string `bson:"namespace" json:"namespace"`
// 	Version   string `bson:"version" json:"version"`

// 	Username string `bson:"username" json:"username"`
// 	Password string `bson:"password" json:"password"`

// 	Replicas    int    `bson:"replicas" json:"replicas"`
// 	ReplicaSet  string `bson:"replicaSet" json:"replicaSet"`
// 	Storage     string `bson:"storage" json:"storage"`
// 	MetalLBPool string `bson:"metallbPool" json:"metallbPool"`

// 	PublicIP         string `bson:"publicIP" json:"publicIP"`
// 	ConnectionString string `bson:"connectionString" json:"connectionString"`

// 	Status    string    `bson:"status" json:"status"`
// 	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
// }

package models

import "time"

type BackupInfo struct {
	BackupName  string    `bson:"backupName" json:"backupName"`
	Status      string    `bson:"status" json:"status"`
	StartedAt   time.Time `bson:"startedAt" json:"startedAt"`
	CompletedAt time.Time `bson:"completedAt" json:"completedAt"`
	TTL         string    `bson:"ttl" json:"ttl"`
}

type DatabaseRecord struct {
	Domain  string `bson:"domain" json:"domain"`
	Project string `bson:"project" json:"project"`
	Cluster string `bson:"cluster" json:"cluster"`

	DBType    string `bson:"dbType" json:"dbType"`
	Name      string `bson:"name" json:"name"`
	Namespace string `bson:"namespace" json:"namespace"`
	Version   string `bson:"version" json:"version"`

	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`

	Replicas    int    `bson:"replicas" json:"replicas"`
	ReplicaSet  string `bson:"replicaSet" json:"replicaSet"`
	Storage     string `bson:"storage" json:"storage"`
	MetalLBPool string `bson:"metallbPool" json:"metallbPool"`

	PublicIP         string `bson:"publicIP" json:"publicIP"`
	ConnectionString string `bson:"connectionString" json:"connectionString"`

	Backup []BackupInfo `bson:"backup,omitempty" json:"backup,omitempty"`

	Status    string    `bson:"status" json:"status"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
