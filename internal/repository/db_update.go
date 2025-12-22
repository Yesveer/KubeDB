package repository

import (
	"context"
	"time"

	"dbaas-orcastrator/internal/db"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateDBVersion(
	domain, project, cluster, name, version string,
) error {

	filter := bson.M{
		"domain":  domain,
		"project": project,
		"cluster": cluster,
		"name":    name,
	}

	update := bson.M{
		"$set": bson.M{
			"version":   version,
			"updatedAt": time.Now(),
		},
	}

	_, err := db.Client.
		Database("compass-config").
		Collection("databases").
		UpdateOne(context.TODO(), filter, update)

	return err
}
