package repository

import (
	"context"
	"time"

	"dbaas-orcastrator/internal/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName    = "compass-config"
	dbaasCollection = "dbaas"
)

// 🔹 Install / Check / Status update (LICENSE SAFE)
func SaveDBaaSRecord(
	domain, project, cluster, status, kubeconfigB64 string,
) error {

	_, err := db.Client.
		Database(databaseName).
		Collection(dbaasCollection).
		UpdateOne(
			context.TODO(),
			bson.M{
				"domain":  domain,
				"project": project,
				"cluster": cluster,
			},
			bson.M{
				"$set": bson.M{
					"domain":         domain,
					"project":        project,
					"cluster":        cluster,
					"status":         status,
					"kubeconfig_b64": kubeconfigB64,
					"updatedAt":      time.Now(),
				},
			},
			options.Update().SetUpsert(true),
		)

	return err
}

// 🔹 License upload only (DO NOT TOUCH OTHER FIELDS)
func SaveLicenseOnly(
	domain, project, cluster, licenseB64 string,
) error {

	_, err := db.Client.
		Database(databaseName).
		Collection(dbaasCollection).
		UpdateOne(
			context.TODO(),
			bson.M{
				"domain":  domain,
				"project": project,
				"cluster": cluster,
			},
			bson.M{
				"$set": bson.M{
					"license_b64": licenseB64,
					"updatedAt":   time.Now(),
				},
			},
			options.Update().SetUpsert(true),
		)

	return err
}
