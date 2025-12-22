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

func InstallClickHouse(db models.DatabaseRecord) {

	log.Println("🚀 Creating ClickHouse:", db.Name)

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
  name: ch-custom-auth
  namespace: "%s"
type: Opaque
stringData:
  username: "%s"
  password: "%s"
---
apiVersion: kubedb.com/v1alpha2
kind: ClickHouse
metadata:
  name: "%s"
  namespace: "%s"
spec:
  version: "%s"
  replicas: %d

  authSecret:
    name: ch-custom-auth

  storageType: Durable
  storage:
    storageClassName: local-path
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: "%s"

  monitor:
    agent: prometheus.io/operator
    prometheus:
      serviceMonitor:
        labels:
          release: kube-prom-stack
        interval: 30s
---
apiVersion: v1
kind: Service
metadata:
  name: "%s-external"
  namespace: "%s"
  annotations:
    metallb.io/address-pool: "%s"
spec:
  type: LoadBalancer
  selector:
    app.kubernetes.io/name: clickhouses.kubedb.com
    app.kubernetes.io/instance: "%s"
  ports:
    - name: http
      port: 8123
      targetPort: 8123
    - name: native
      port: 9000
      targetPort: 9000
---
apiVersion: autoscaling.kubedb.com/v1alpha1
kind: ClickHouseAutoscaler
metadata:
  name: "%s-autoscaler"
  namespace: "%s"
spec:
  databaseRef:
    name: "%s"
  compute:
    clickhouse:
      minAllowed:
        cpu: 1
        memory: 2Gi
      maxAllowed:
        cpu: 2
        memory: 3Gi
      trigger: "On"
      controlledResources:
        - cpu
        - memory
      containerControlledValues: RequestsAndLimits

  storage:
    clickhouse:
      expansionMode: Online
      trigger: "On"
      usageThreshold: 60
      scalingThreshold: 100
`,
		// Secret
		db.Namespace,
		db.Username,
		db.Password,

		// ClickHouse
		db.Name,
		db.Namespace,
		db.Version,
		db.Replicas,
		db.Storage,

		// Service
		db.Name,
		db.Namespace,
		db.MetalLBPool,
		db.Name,

		// Autoscaler
		db.Name,
		db.Namespace,
		db.Name,
	)

	tmp := "/tmp/clickhouse.yaml"
	if err := os.WriteFile(tmp, []byte(yaml), 0644); err != nil {
		log.Println("❌ Failed to write yaml:", err)
		return
	}

	// 3️⃣ APPLY YAML
	log.Println("📄 Applying ClickHouse YAML")
	run("kubectl apply -f " + tmp)

	// 4️⃣ WAIT FOR POD
	log.Println("⏳ Waiting for ClickHouse to be ready...")
	if err := WaitForClickHouseReady(kubeconfig, db.Name, db.Namespace); err != nil {
		log.Println("❌ ClickHouse not ready:", err)
		return
	}

	// 5️⃣ GET LB IP
	lbIP := strings.TrimSpace(run(
		fmt.Sprintf("kubectl get svc %s-external -n %s -o jsonpath='{.status.loadBalancer.ingress[0].ip}'",
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
		"clickhouse://%s:%s@%s:8123/",
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
