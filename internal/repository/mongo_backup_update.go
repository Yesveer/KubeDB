package repository

import (
	"context"
	"fmt"

	"dbaas-orcastrator/internal/db"
	"dbaas-orcastrator/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func AppendMongoBackup(
	domain, project, cluster, name string,
	backup models.BackupInfo,
) error {

	filter := bson.M{
		"domain":  domain,
		"project": project,
		"cluster": cluster,
		"name":    name, // 🔥 mongo-rs
	}

	update := bson.M{
		"$push": bson.M{
			"backup": backup,
		},
		"$set": bson.M{
			"updatedAt": backup.CompletedAt,
		},
	}

	res, err := db.Client.
		Database("compass-config").
		Collection("databases").
		UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("no database record found to attach backup")
	}

	return nil
}
