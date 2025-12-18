package routes

import (
	"net/http"

	"dbaas-orcastrator/internal/handlers"
)

func Register(k *handlers.KubeDBHandler) {
	http.HandleFunc("/install/kubedb", k.Install)
	http.HandleFunc("/check/kubedb", k.Check)
	http.HandleFunc("/license/kubedb", k.GenerateLicense)
	http.HandleFunc("/license/kubedb/upload", k.UploadLicense)
	http.HandleFunc("/backup/enable", k.EnableBackup)
	http.HandleFunc("/dbaas/check", k.CheckDBaaS)
	http.HandleFunc("/dbaas/list", k.ListDBaaS)
	http.HandleFunc("/metallb/add", k.AddMetalLBPool)
	http.HandleFunc("/metallb/list", k.ListMetalLBPools)
}
