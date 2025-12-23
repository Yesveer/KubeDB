package repository

import (
	"context"
	"errors"

	"dbaas-orcastrator/internal/db"
	"dbaas-orcastrator/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func GetMongoRecord(domain, project, cluster, name string) (*models.DatabaseRecord, error) {

	filter := bson.M{
		"domain":  domain,
		"project": project,
		"cluster": cluster,
		"name":    name,
	}

	var rec models.DatabaseRecord
	err := db.Client.
		Database("compass-config").
		Collection("databases").
		FindOne(context.TODO(), filter).
		Decode(&rec)

	if err != nil {
		return nil, err
	}

	if rec.DBType != "mongo" {
		return nil, errors.New("this is not mongo cluster")
	}

	return &rec, nil
}
