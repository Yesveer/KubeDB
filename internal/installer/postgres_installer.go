package installer

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"dbaas-orcastrator/internal/models"
	"dbaas-orcastrator/internal/repository"
)

func InstallPostgres(db models.DatabaseRecord) {

	log.Println("🚀 Creating Postgres:", db.Name)

	kubeconfig := os.Getenv("KUBECONFIG_PATH")

	run := func(cmdStr string) string {
		cmd := exec.Command("bash", "-c", cmdStr)
		cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfig)
		out, err := cmd.CombinedOutput()
		log.Println(string(out))
		if err != nil {
			log.Println("❌ command failed:", err)
		}
		return string(out)
	}

	// 1️⃣ CREATE NAMESPACE (FIXED)
	log.Println("📦 Creating namespace:", db.Namespace)
	run(fmt.Sprintf(`
kubectl apply -f - <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: %s
EOF
`, db.Namespace))

	// 2️⃣ WRITE YAML (FIXED)
	// 2️⃣ WRITE YAML (FIXED & CLEAN)
	yaml := fmt.Sprintf(`
apiVersion: v1
kind: Secret
metadata:
  name: pg-custom-auth
  namespace: %s
type: kubernetes.io/basic-auth
stringData:
  username: %s
  password: "%s"
---
apiVersion: kubedb.com/v1
kind: Postgres
metadata:
  name: %s
  namespace: %s
spec:
  version: "%s"

  authSecret:
    name: pg-custom-auth
    externallyManaged: true

  replicas: %d
  standbyMode: Hot
  streamingMode: Asynchronous

  storageType: Durable
  storage:
    storageClassName: local-path
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: "%s"

  leaderElection:
    leaseDurationSeconds: 15
    renewDeadlineSeconds: 10
    retryPeriodSeconds: 2

  serviceTemplates:
    - alias: primary
      metadata:
        annotations:
          metallb.io/address-pool: "%s"
      spec:
        type: LoadBalancer
        ports:
          - name: postgres
            port: 5432

    - alias: standby
      spec:
        type: ClusterIP
        ports:
          - name: postgres
            port: 5432

  monitor:
    agent: prometheus.io/operator
    prometheus:
      exporter:
        port: 9187
      serviceMonitor:
        labels:
          release: kube-prom-stack
        interval: 30s
`,
		// 🔐 Secret
		db.Namespace,
		db.Username,
		db.Password,

		// 🐘 Postgres
		db.Name,
		db.Namespace,
		db.Version,
		db.Replicas,
		db.Storage,

		// 🔌 Service (MetalLB)
		db.MetalLBPool,
	)

	tmp := "/tmp/postgres.yaml"
	if err := os.WriteFile(tmp, []byte(yaml), 0644); err != nil {
		log.Println("❌ Failed to write yaml:", err)
		return
	}

	// 3️⃣ APPLY YAML
	log.Println("📄 Applying Postgres YAML")
	run("kubectl apply -f " + tmp)

	// 4️⃣ WAIT FOR POD
	log.Println("⏳ Waiting for Postgres to be ready...")
	if err := WaitForPostgresReady(kubeconfig, db.Name, db.Namespace); err != nil {
		log.Println("❌ Postgres not ready:", err)
		return
	}

	// 5️⃣ GET LB IP
	lbIP := strings.TrimSpace(run(
		fmt.Sprintf("kubectl get svc %s -n %s -o jsonpath='{.status.loadBalancer.ingress[0].ip}'",
			db.Name, db.Namespace),
	))

	if lbIP == "" {
		log.Println("❌ LoadBalancer IP not assigned")
		return
	}

	// 6️⃣ GET NODE IP
	nodeIP, err := GetFirstNodeIP(kubeconfig)
	if err != nil {
		log.Println("❌ Failed to get node IP:", err)
		return
	}

	// 7️⃣ ADD IP TO NODE
	log.Println("🔌 Adding IP to node:", lbIP)
	if err := AddIPToNode(
		nodeIP,
		lbIP,
	); err != nil {
		log.Println("❌ Failed to add IP:", err)
	}

	// 8️⃣ CONNECTION STRING
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/",
		db.Username,
		db.Password,
		lbIP,
	)

	// 9️⃣ UPDATE DB
	if err := repository.UpdateMongoRunning(db, lbIP, conn); err != nil {
		log.Println("❌ DB update failed:", err)
		return
	}

	log.Println("✅ MongoDB RUNNING")
	log.Println("🔗 Connection:", conn)
}
