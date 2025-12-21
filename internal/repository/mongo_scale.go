package repository

import (
	"context"
	"time"

	"dbaas-orcastrator/internal/db"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateMongoScale(
	domain, project, name string,
	replicas int,
	storage string,
) error {

	_, err := db.Client.
		Database("compass-config").
		Collection("databases").
		UpdateOne(
			context.TODO(),
			bson.M{
				"domain":  domain,
				"project": project,
				"name":    name,
			},
			bson.M{
				"$set": bson.M{
					"replicas":  replicas,
					"storage":   storage,
					"updatedAt": time.Now(),
				},
			},
		)

	return err
}
