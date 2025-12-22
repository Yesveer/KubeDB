// package repository

// import (
// 	"context"

// 	"dbaas-orcastrator/internal/db"
// 	"dbaas-orcastrator/internal/models"

// 	"go.mongodb.org/mongo-driver/bson"
// )

// func ListAllDBaaS() ([]models.DBaaSRecord, error) {

// 	cursor, err := db.Client.
// 		Database("compass-config").
// 		Collection("dbaas").
// 		Find(context.TODO(), bson.M{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(context.TODO())

// 	var records []models.DBaaSRecord
// 	if err := cursor.All(context.TODO(), &records); err != nil {
// 		return nil, err
// 	}

// 	return records, nil
// }

package repository

import (
	"context"

	"dbaas-orcastrator/internal/db"
	"dbaas-orcastrator/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func ListDBaaSByDomainProject(domain, project string) ([]models.DBaaSRecord, error) {

	filter := bson.M{
		"domain":  domain,
		"project": project,
	}

	cursor, err := db.Client.
		Database("compass-config").
		Collection("dbaas").
		Find(context.TODO(), filter)
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

func ListAllDBaaS() ([]models.DBaaSRecord, error) {

	var result []models.DBaaSRecord

	cursor, err := db.Client.
		Database("compass-config").
		Collection("dbaas").
		Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	return result, nil
}
