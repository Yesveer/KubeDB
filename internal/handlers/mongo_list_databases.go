package handlers

import (
	"context"
	"dbaas-orcastrator/internal/repository"
	"dbaas-orcastrator/internal/services"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func (h *KubeDBHandler) ListMongoDatabases(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Domain  string `json:"domain"`
		Project string `json:"project"`
		Cluster string `json:"cluster"`
		Name    string `json:"name"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	rec, err := repository.GetMongoRecord(
		req.Domain, req.Project, req.Cluster, req.Name,
	)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	client, err := services.NewMongoClient(rec.ConnectionString)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	dbs, err := client.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"databases": dbs,
	})
}
