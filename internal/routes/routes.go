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
	http.HandleFunc("/metallb/delete", k.DeleteMetalLBPool)
	http.HandleFunc("/db/versions", k.GetDBVersions)
	http.HandleFunc("/database/mongo/create", k.CreateMongoDB)
	http.HandleFunc("/database/clickhouse/create", k.CreateClickHouse)
	http.HandleFunc("/database/kafka/create", k.CreateKafka)
	http.HandleFunc("/database/mysql/create", k.CreateMySql)
	http.HandleFunc("/database/postgres/create", k.CreatePostgres)
	http.HandleFunc("/database/redis/create", k.RedisCreateRequest)
	http.HandleFunc("/databases", k.GetDatabases)
	http.HandleFunc("/database/delete", k.DeleteDatabase)
	http.HandleFunc("/database/scale", k.ScaleMongoDB)
	http.HandleFunc("/backup/create", k.CreateBackup)
	http.HandleFunc("/database/upgrade", k.UpgradeDatabase)
	// ✅ MONITORING ROUTES - ADD THESE AT THE END
	// The following handlers are commented out because the corresponding methods are not implemented in KubeDBHandler.
	// Uncomment and implement these methods in the handler if needed.
	http.HandleFunc("/monitoring/health", k.MonitoringHealthCheck)
	http.HandleFunc("/monitoring/query", k.QueryMetrics)
	http.HandleFunc("/monitoring/query_range", k.QueryRangeMetrics)
}
