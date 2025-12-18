package repository

import (
	"context"

	"dbaas-orcastrator/internal/db"
	"dbaas-orcastrator/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindDBaaSRecord(domain, project, cluster string) (*models.DBaaSRecord, error) {

	var record models.DBaaSRecord

	err := db.Client.
		Database("compass-config").
		Collection("dbaas").
		FindOne(
			context.TODO(),
			bson.M{
				"domain":  domain,
				"project": project,
				"cluster": cluster,
			},
		).
		Decode(&record)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &record, nil
}
