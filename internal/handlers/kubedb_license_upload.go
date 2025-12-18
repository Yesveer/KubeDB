package handlers

import (
	"encoding/json"
	"net/http"

	"dbaas-orcastrator/internal/licence"
	"dbaas-orcastrator/internal/repository"
)

type UploadLicenseRequest struct {
	Domain  string `json:"domain"`
	Project string `json:"project"`
	Cluster string `json:"cluster"`
	License string `json:"license"`
}

func (h *KubeDBHandler) UploadLicense(w http.ResponseWriter, r *http.Request) {

	var req UploadLicenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", 400)
		return
	}

	if req.License == "" {
		http.Error(w, "License content is empty", 400)
		return
	}

	// 1️⃣ Encode license
	licenseB64 := licence.EncodeToBase64(req.License)

	// 2️⃣ Overwrite local license file
	if err := licence.WriteToFile("./scripts/licence.txt", req.License); err != nil {
		http.Error(w, "Failed to write license file: "+err.Error(), 500)
		return
	}

	// 3️⃣ Update DB (license only)
	if err := repository.SaveLicenseOnly(
		req.Domain,
		req.Project,
		req.Cluster,
		licenseB64,
	); err != nil {
		http.Error(w, "Mongo save failed: "+err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "License uploaded and updated successfully",
	})
}
