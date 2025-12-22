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

func InstallRedis(db models.DatabaseRecord) {

	log.Println("🚀 Creating Redis:", db.Name)

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
  name: redis-custom-auth
  namespace: %s
type: kubernetes.io/basic-auth
stringData:
  username: %s
  password: "%s"
---
apiVersion: kubedb.com/v1
kind: Redis
metadata:
  name: %s
  namespace: %s
spec:
  version: "%s"
  mode: Standalone
  replicas: %d

  authSecret:
    name: redis-custom-auth
    externallyManaged: true

  storageType: Durable
  storage:
    storageClassName: local-path
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: "%s"

  serviceTemplates:
    - alias: primary
      metadata:
        annotations:
          metallb.io/address-pool: "%s"
      spec:
        type: LoadBalancer
        ports:
          - name: redis
            port: 6379

  monitor:
    agent: prometheus.io/operator
    prometheus:
      exporter:
        port: 9121
      serviceMonitor:
        labels:
          release: kube-prom-stack
        interval: 30s
---
apiVersion: autoscaling.kubedb.com/v1alpha1
kind: RedisAutoscaler
metadata:
  name: "%s-autoscaler"
  namespace: "%s"
spec:
  databaseRef:
    name: "%s"
  compute:
    standalone:
      minAllowed:
        cpu: 400m
        memory: 400Mi
      maxAllowed:
        cpu: "1"
        memory: 1Gi
      trigger: "On"
      controlledResources:
        - cpu
        - memory
      containerControlledValues: RequestsAndLimits

  storage:
    standalone:
      expansionMode: Online
      trigger: "On"
      usageThreshold: 70
      scalingThreshold: 60
`,
		// 🔐 Secret
		db.Namespace,
		db.Username,
		db.Password,

		// 🔴 Redis
		db.Name,
		db.Namespace,
		db.Version,
		db.Replicas,
		db.Storage,

		// 🔌 Service (MetalLB)
		db.MetalLBPool,

		// 🧮 Autoscaler
		db.Name,      // autoscaler name
		db.Namespace, // namespace
		db.Name,      // redis name
	)

	tmp := "/tmp/redis.yaml"
	if err := os.WriteFile(tmp, []byte(yaml), 0644); err != nil {
		log.Println("❌ Failed to write yaml:", err)
		return
	}

	// 3️⃣ APPLY YAML
	log.Println("📄 Applying Redis YAML")
	run("kubectl apply -f " + tmp)

	// 4️⃣ WAIT FOR POD
	log.Println("⏳ Waiting for Redis to be ready...")
	if err := WaitForRedisReady(kubeconfig, db.Name, db.Namespace); err != nil {
		log.Println("❌ Redis not ready:", err)
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
		// "rediss://%s:%s@%s:6379/",
		// `redis-cli -h %s -p 6379 -a %s`,
		`redis://%s:%s@%s:6379`,
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
