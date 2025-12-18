package repository

import (
	"context"

	"dbaas-orcastrator/internal/db"
	"dbaas-orcastrator/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func ListAllDBaaS() ([]models.DBaaSRecord, error) {

	cursor, err := db.Client.
		Database("compass-config").
		Collection("dbaas").
		Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var records []models.DBaaSRecord
	if err := cursor.All(context.TODO(), &records); err != nil {
		return nil, err
	}

	return records, nil
}
