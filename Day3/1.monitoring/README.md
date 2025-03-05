# **Deploying Prometheus and Grafana for Full Cluster Monitoring in CRC**

This guide outlines how to deploy **Prometheus and Grafana** in a dedicated namespace while ensuring Prometheus has **cluster-wide monitoring capabilities** in CRC. The setup includes:

✅ Deploying **Prometheus and Grafana**
✅ Setting up **RBAC permissions** for Prometheus to monitor the cluster
✅ Enabling **Service Discovery** to automatically detect Nodes, Pods, and Services
✅ **Automating Prometheus integration in Grafana** using a ConfigMap

---

## **1️⃣ Deployment Files**
This setup includes the following YAML files:

### **1. `namespace.yaml`** (Creates `monitoring` namespace)
```sh
oc apply -f namespace.yaml
```

### **2. `prometheus.yaml`** (Deploys Prometheus in the `monitoring` namespace)
```sh
oc apply -f prometheus.yaml
```

### **3. `grafana.yaml`** (Deploys Grafana in the `monitoring` namespace)
```sh
oc apply -f grafana.yaml
```

### **4. `cr.yaml`** (Creates ClusterRole & ClusterRoleBinding for Prometheus)
```sh
oc apply -f cr.yaml
```

### **5. `service-discovery.yaml`** (Configures Prometheus to scrape Cluster metrics)
```sh
oc apply -f service-discovery.yaml
```

### **6. `grafana-datasource.yaml`** (Automates Prometheus integration in Grafana)
```sh
oc apply -f grafana-datasource.yaml
```

## **7. `dashboard.yaml`** (add Dashboard to Grafana)
```sh
oc apply -f dashboard.yaml
```

After applying all these manifests, restart Prometheus to ensure the updated configurations take effect:
```sh
oc delete pod -n monitoring -l app=prometheus
oc delete pod -n monitoring -l app=grafana
```

---

## **2️⃣ Checking the Deployment**
Once deployed, verify that all components are running:
```sh
oc get pods -n monitoring
```
Expected output:
```
NAME                           READY   STATUS    RESTARTS   AGE
prometheus-xxx                 1/1     Running   0          Xm
grafana-xxx                    1/1     Running   0          Xm
```

---

## **3️⃣ Accessing Grafana**
Retrieve the **Grafana URL**:
```sh
oc get route grafana-service -n monitoring
```
Open the URL in a browser and log in with:
- **Username:** `admin`
- **Password:** `admin` (or your configured password)

### **🔹 Prometheus Integration is now Automated**
Grafana will automatically have Prometheus as a data source configured via the applied ConfigMap.

If you need to verify, go to **Configuration → Data Sources**, and you should see Prometheus already set up.

### **🔹 Importing Kubernetes Cluster Dashboard**
1️⃣ Navigate to **Create → Import**
2️⃣ Enter **Dashboard ID:** `315`
3️⃣ Select **Prometheus** as the data source and click **Import**

**Grafana now visualizes real-time CRC cluster metrics!**

---

## ** Summary**
- **`namespace.yaml`** → Creates `monitoring` namespace
- **`prometheus.yaml`** → Deploys Prometheus
- **`grafana.yaml`** → Deploys Grafana
- **`cr.yaml`** → Grants Prometheus cluster monitoring permissions
- **`service-discovery.yaml`** → Enables Prometheus Service Discovery
- **`grafana-datasource.yaml`** → Configures Prometheus in Grafana automatically

---
**Copyright (c) 2025 by Alexander Kolin. All rights reserved.**

