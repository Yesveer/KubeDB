package repository

import (
	"context"
	"time"

	"dbaas-orcastrator/internal/db"
	"dbaas-orcastrator/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertMongoDB(d models.DatabaseRecord) error {
	d.Status = "CREATING"
	d.CreatedAt = time.Now()

	_, err := db.Client.
		Database("compass-config").
		Collection("databases").
		InsertOne(context.TODO(), d)

	return err
}

func UpdateMongoRunning(
	d models.DatabaseRecord,
	publicIP string,
	conn string,
) error {

	_, err := db.Client.
		Database("compass-config").
		Collection("databases").
		UpdateOne(
			context.TODO(),
			bson.M{
				"domain":  d.Domain,
				"project": d.Project,
				"cluster": d.Cluster,
				"name":    d.Name,
			},
			bson.M{
				"$set": bson.M{
					"status":           "RUNNING",
					"publicIP":         publicIP,
					"connectionString": conn,
				},
			},
		)

	return err
}
