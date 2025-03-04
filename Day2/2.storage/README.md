# Advanced Lab: Storage Monitoring & Debugging in OpenShift/CRC

## Objective
- Test storage performance using `fio`
- Monitor storage utilization via **CRC Console (Prometheus UI)**
- Configure **storage alerts for PVCs**
- Debug storage issues in OpenShift

---

## 1. Persistent Volume and Persistent Volume Claim
Ensure a **PV and PVC** exist before proceeding.

Apply the following YAML files:
```sh
oc apply -f pv.yaml
oc apply -f pvc.yaml
```

Verify PVC binding:
```sh
oc get pv,pvc
```

---

## 2. Storage Performance Analysis with `fio`
Deploy a `fio` test pod using `fio-pod.yaml`.
```sh
oc apply -f fio-pod.yaml
```

Run performance tests inside the pod:
```sh
oc exec -it fio-test -- sh
```
Then execute:
```sh
fio --name=randwrite --ioengine=libaio --rw=randwrite --bs=4k --size=1G --numjobs=4 --runtime=60 --group_reporting
```

---

## 3. Storage Monitoring using CRC Console (Prometheus UI)
### 3.1 Access the OpenShift Web Console
- Open **https://console-openshift-console.apps-crc.testing/**
- Navigate to **Observe → Metrics**
- Use the following PromQL queries:

### 3.2 Monitor PVC Storage Usage
```text
kubelet_volume_stats_used_bytes{persistentvolumeclaim="pvc-example"}
```
Displays the current PVC storage consumption.

### 3.3 Monitor Total Storage Usage per Namespace
```text
sum by (namespace) (kubelet_volume_stats_used_bytes)
```
Summarizes storage consumption for all PVCs in a namespace.

### 3.4 Measure Storage IOPS
```text
rate(node_disk_io_time_seconds_total[5m])
```
Displays IOPS for node storage systems.

---

## 4. Configure Storage Alerts
Apply `pvc-alert.yaml` to configure alerts in Prometheus:
```sh
oc apply -f pvc-alert.yaml
```

This triggers an alert when **PVC usage exceeds 50%**.

---

## 5. Storage Debugging
If a pod is stuck in **ContainerCreating**:
```sh
oc describe pod <pod-name>
```

### Common Issues & Solutions
| Issue | Solution |
|--------|--------|
| **PVC pending** | Check `oc describe pvc`, validate storage class |
| **Mount failures** | Check `oc logs -p <pod>`, `/var/log/messages` on the node |
| **IO errors** | Verify storage class, run `fio` test |
| **Slow read/write** | Use PromQL: `rate(node_disk_io_time_seconds_total[5m])` |

To debug storage in a running pod:
```sh
oc debug pod/fio-test
```

---
**Copyright (c) 2025 by Alexander Kolin. All rights reserved.**
