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

func InstallKafka(db models.DatabaseRecord) {

	log.Println("🚀 Creating Kafka:", db.Name)

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
  name: kafka-custom-auth
  namespace: %s
type: kubernetes.io/basic-auth
stringData:
  username: %s
  password: %s
---
apiVersion: kubedb.com/v1
kind: Kafka
metadata:
  name: %s
  namespace: %s
spec:
  version: "%s"

  authSecret:
    name: kafka-custom-auth
    externallyManaged: true

  topology:
    controller:
      replicas: %d
      storage:
        storageClassName: local-path
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: "%s"
    broker:
      replicas: %d
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
      exporter:
        port: 9091
      serviceMonitor:
        labels:
          release: kube-prom-stack
        interval: 30s
---
apiVersion: v1
kind: Service
metadata:
  name: %s-broker-external
  namespace: %s
  annotations:
    metallb.io/address-pool: %s
spec:
  type: LoadBalancer
  selector:
    app.kubernetes.io/name: kafkas.kubedb.com
    app.kubernetes.io/instance: %s
    kubedb.com/role: broker
  ports:
    - name: kafka
      protocol: TCP
      port: 9092
      targetPort: 9092
---
apiVersion: v1
kind: Service
metadata:
  name: %s-controller-external
  namespace: %s
  annotations:
    metallb.io/address-pool: %s
spec:
  type: LoadBalancer
  selector:
    app.kubernetes.io/name: kafkas.kubedb.com
    app.kubernetes.io/instance: %s
    kubedb.com/role: controller
  ports:
    - name: kafka-controller
      protocol: TCP
      port: 9093
      targetPort: 9093
`,
		// 🔐 Secret
		db.Namespace,
		db.Username,
		db.Password,

		// 🟠 Kafka CR
		db.Name,
		db.Namespace,
		db.Version,
		db.Replicas, // controller replicas
		db.Storage,  // controller storage
		db.Replicas, // broker replicas
		db.Storage,  // broker storage

		// 🔌 Broker Service
		db.Name,
		db.Namespace,
		db.MetalLBPool,
		db.Name,

		// 🔌 Controller Service
		db.Name,
		db.Namespace,
		db.MetalLBPool,
		db.Name,
	)

	tmp := "/tmp/kafka.yaml"
	if err := os.WriteFile(tmp, []byte(yaml), 0644); err != nil {
		log.Println("❌ Failed to write yaml:", err)
		return
	}

	// 3️⃣ APPLY YAML
	log.Println("📄 Applying Kafka YAML")
	run("kubectl apply -f " + tmp)

	// 4️⃣ WAIT FOR POD
	log.Println("⏳ Waiting for Kafka to be ready...")
	if err := WaitForKafkaReady(kubeconfig, db.Name, db.Namespace); err != nil {
		log.Println("❌ Kafka not ready:", err)
		return
	}

	// 5️⃣ GET LB IP
	lbIP := strings.TrimSpace(run(
		fmt.Sprintf("kubectl get svc %s-broker-external -n %s -o jsonpath='{.status.loadBalancer.ingress[0].ip}'",
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
		"kafka://%s:%s@%s:9092/",
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
