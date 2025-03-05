# **Velero Backup and Restore Setup for CRC with PVC Storage**

This guide outlines the deployment of **Velero in CRC**, storing backups in a **PersistentVolumeClaim (PVC)** instead of external storage. The setup ensures **backup and restore functionality** for the `db` namespace.

## **1️⃣ Deployment Files Overview**
This setup includes the following YAML files:

| File | Purpose |
|------|---------|
| `namespace.yaml` | Creates the `velero` namespace |
| `pvc.yaml` | Defines a PVC (`velero-backup-pvc`) for backup storage |
| `rbac.yaml` | Grants Velero the necessary permissions |
| `deployment.yaml` | Deploys Velero using the PVC as storage |
| `backup.yaml` | Defines a backup job for the `db` namespace |
| `restore.yaml` | Defines a restore job from a backup |

---

## **2️⃣ Deploying Velero in CRC**
Run the following commands to deploy Velero step by step:

1️⃣ **Create the namespace:**
```sh
oc apply -f namespace.yaml
```

2️⃣ **Create the PVC for storing backups:**
```sh
oc apply -f pvc.yaml
```

3️⃣ **Grant necessary permissions:**
```sh
oc apply -f rbac.yaml
```

4️⃣ **Deploy Velero:**
```sh
oc apply -f deployment.yaml
```

---

## **3️⃣ Install Velero CLI**
To interact with Velero, you need to install the Velero CLI on your local machine.

### ** Installation on Linux/macOS**
```sh
curl -LO https://github.com/vmware-tanzu/velero/releases/download/v1.15.2/velero-v1.15.2-linux-amd64.tar.gz
tar -xvf velero-v1.15.2-linux-amd64.tar.gz
sudo mv velero-v1.15.2-linux-amd64/velero /usr/local/bin/
```

### ** Verify Installation**
```sh
velero version
```
➡ **Expected Output:**
```
Client:
        Version: v1.xx.x
```

---

### ** Install CRDs **
```sh
velero install --crds-only --dry-run -o yaml | oc apply -f -
```



## **4️⃣ Creating a Backup of the `db` Namespace**
Once Velero is deployed, create a backup of the `db` namespace:
```sh
oc apply -f backup.yaml
```

To check the backup status:
```sh
velero backup describe db-backup --details
```

To list all backups:
```sh
velero backup get
```

---

## **5️⃣ Restoring the `db` Namespace from Backup**
If needed, restore the `db` namespace from the last backup:
```sh
oc apply -f restore.yaml
```

To check the restore status:
```sh
velero restore describe db-restore --details
```

To list all restores:
```sh
velero restore get
```

---

## ** Summary**
- **Velero is deployed in CRC with PVC-based storage**
- **Backups and restores are managed via YAML manifests**
- **Permissions are granted via RBAC to allow full access**
- **No external storage is needed** – everything runs inside CRC
- **Velero CLI is required for managing backups and restores**

---
**Copyright (c) 2025 by Alexander Kolin. All rights reserved.**