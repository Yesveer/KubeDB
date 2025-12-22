package repository

import (
	"context"

	"dbaas-orcastrator/internal/db"
	"dbaas-orcastrator/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func GetDatabases(domain, project string) ([]models.DatabaseRecord, error) {

	cursor, err := db.Client.
		Database("compass-config").
		Collection("databases").
		Find(context.TODO(), bson.M{
			"domain":  domain,
			"project": project,
		})
	if err != nil {
		return nil, err
	}

	var list []models.DatabaseRecord
	if err := cursor.All(context.TODO(), &list); err != nil {
		return nil, err
	}

	return list, nil
}

func GetDatabaseByName(domain, project, name string) (*models.DatabaseRecord, error) {

	var rec models.DatabaseRecord

	err := db.Client.
		Database("compass-config").
		Collection("databases").
		FindOne(context.TODO(), bson.M{
			"domain":  domain,
			"project": project,
			"name":    name,
		}).
		Decode(&rec)

	if err != nil {
		return nil, nil
	}

	return &rec, nil
}

func GetAllDatabases() ([]models.DatabaseRecord, error) {

	var result []models.DatabaseRecord

	cursor, err := db.Client.
		Database("compass-config").
		Collection("databases").
		Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	return result, nil
}
