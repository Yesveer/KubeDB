package repository

import (
	"context"
	"errors"

	"dbaas-orcastrator/internal/db"
	"dbaas-orcastrator/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func GetLatestCompletedBackup(
	domain, project, cluster, name string,
) (*models.BackupInfo, error) {

	filter := bson.M{
		"domain":  domain,
		"project": project,
		"cluster": cluster,
		"name":    name,
	}

	var record models.DatabaseRecord
	err := db.Client.
		Database("compass-config").
		Collection("database").
		FindOne(context.TODO(), filter).
		Decode(&record)

	if err != nil {
		return nil, err
	}

	// 🔥 Reverse iterate = latest first
	for i := len(record.Backup) - 1; i >= 0; i-- {
		if record.Backup[i].Status == "Completed" {
			return &record.Backup[i], nil
		}
	}

	return nil, errors.New("no completed backup found")
}
