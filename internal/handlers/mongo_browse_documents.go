package handlers

import (
	"context"
	"dbaas-orcastrator/internal/repository"
	"dbaas-orcastrator/internal/services"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *KubeDBHandler) BrowseMongoDocuments(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Domain     string `json:"domain"`
		Project    string `json:"project"`
		Cluster    string `json:"cluster"`
		Name       string `json:"name"`
		Database   string `json:"database"`
		Collection string `json:"collection"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	if req.Database == "" || req.Collection == "" {
		http.Error(w, "database and collection required", 400)
		return
	}

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

	cursor, err := client.
		Database(req.Database).
		Collection(req.Collection).
		Find(context.TODO(), bson.M{}, options.Find().SetLimit(50))

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var docs []map[string]interface{}
	cursor.All(context.TODO(), &docs)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(docs),
		"data":  docs,
	})
}
