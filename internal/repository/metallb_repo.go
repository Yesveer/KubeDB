package repository

import (
	"context"
	"time"

	"dbaas-orcastrator/internal/db"
	"dbaas-orcastrator/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const metallbCollection = "dbaas-ip"

// Add / Update pool
func SaveMetalLBPool(pool models.MetalLBPool) error {

	pool.CreatedAt = time.Now()

	_, err := db.Client.
		Database("compass-config").
		Collection(metallbCollection).
		UpdateOne(
			context.TODO(),
			bson.M{
				"domain":    pool.Domain,
				"project":   pool.Project,
				"cluster":   pool.Cluster,
				"pool_name": pool.PoolName,
			},
			bson.M{
				"$set": pool,
			},
			options.Update().SetUpsert(true), // ✅ CORRECT
		)

	return err
}

// List all pools

func ListMetalLBPoolsFiltered(
	domain, project, cluster string,
) ([]models.MetalLBPool, error) {

	filter := bson.M{}

	if domain != "" {
		filter["domain"] = domain
	}
	if project != "" {
		filter["project"] = project
	}
	if cluster != "" {
		filter["cluster"] = cluster
	}

	cursor, err := db.Client.
		Database("compass-config").
		Collection("dbaas-ip").
		Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var pools []models.MetalLBPool
	if err := cursor.All(context.TODO(), &pools); err != nil {
		return nil, err
	}

	return pools, nil
}
