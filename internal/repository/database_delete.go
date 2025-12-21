package repository

import (
	"context"

	"dbaas-orcastrator/internal/db"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteDatabase(domain, project, name string) error {

	_, err := db.Client.
		Database("compass-config").
		Collection("databases").
		DeleteOne(context.TODO(), bson.M{
			"domain":  domain,
			"project": project,
			"name":    name,
		})

	return err
}
