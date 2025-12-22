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

func InstallMongo(db models.DatabaseRecord) {

	log.Println("🚀 Creating MongoDB:", db.Name)

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
	yaml := fmt.Sprintf(`
apiVersion: v1
kind: Secret
metadata:
  name: mongo-custom-auth
  namespace: %s
type: kubernetes.io/basic-auth
stringData:
  username: %s
  password: "%s"
---
apiVersion: kubedb.com/v1
kind: MongoDB
metadata:
  name: %s
  namespace: %s
spec:
  version: "%s"
  replicas: %d

  replicaSet:
    name: %s

  storage:
    storageClassName: local-path
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: %s

  serviceTemplates:
    - alias: primary
      metadata:
        annotations:
          metallb.io/address-pool: %s
      spec:
        type: LoadBalancer
        ports:
          - name: mongodb
            port: 27017

    - alias: standby
      metadata:
        annotations:
          metallb.io/address-pool: %s
      spec:
        type: LoadBalancer
        ports:
          - name: mongodb
            port: 27017

  monitor:
    agent: prometheus.io/operator
    prometheus:
      exporter:
        port: 9216
      serviceMonitor:
        labels:
          release: kube-prom-stack
        interval: 30s
---
apiVersion: ops.kubedb.com/v1alpha1
kind: MongoDBOpsRequest
metadata:
  name: %s-rotate-auth
  namespace: %s
spec:
  type: RotateAuth
  databaseRef:
    name: %s
  authentication:
    secretRef:
      kind: Secret
      name: mongo-custom-auth
  timeout: 5m
  apply: IfReady
`,
		db.Namespace, // namespace
		db.Username,  // username
		db.Password,  // password

		db.Name,        // mongodb name
		db.Namespace,   // namespace
		db.Version,     // version
		db.Replicas,    // replicas (INT ✅)
		db.ReplicaSet,  // rs name
		db.Storage,     // storage
		db.MetalLBPool, // primary pool
		db.MetalLBPool, // standby pool

		db.Name,      // opsrequest name prefix
		db.Namespace, // opsrequest namespace
		db.Name,      // databaseRef name
	)
	tmp := "/tmp/mongo.yaml"
	if err := os.WriteFile(tmp, []byte(yaml), 0644); err != nil {
		log.Println("❌ Failed to write yaml:", err)
		return
	}

	// 3️⃣ APPLY YAML
	log.Println("📄 Applying MongoDB YAML")
	run("kubectl apply -f " + tmp)

	// 4️⃣ WAIT FOR POD
	log.Println("⏳ Waiting for MongoDB to be ready...")
	if err := WaitForMongoReady(kubeconfig, db.Name, db.Namespace); err != nil {
		log.Println("❌ Mongo not ready:", err)
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
		"mongodb://%s:%s@%s:27017/admin?authSource=admin",
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
