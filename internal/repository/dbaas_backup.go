package repository

import (
	"context"
	"time"

	"dbaas-orcastrator/internal/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnableBackup(
	domain, project, cluster, s3Url, accessKey, secretKey string,
) error {

	_, err := db.Client.
		Database("compass-config").
		Collection("dbaas").
		UpdateOne(
			context.TODO(),
			bson.M{
				"domain":  domain,
				"project": project,
				"cluster": cluster,
			},
			bson.M{
				"$set": bson.M{
					"backup_enabled": true,
					"s3": bson.M{
						"endpoint":   s3Url,
						"access_key": accessKey,
						"secret_key": secretKey,
					},
					"updatedAt": time.Now(),
				},
			},
			options.Update().SetUpsert(true),
		)

	return err
}
