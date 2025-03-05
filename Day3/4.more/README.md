# **Extending Velero Backup with MinIO in CRC**

This guide outlines how to extend Velero backup in CRC by deploying **MinIO** as an S3-compatible object storage and configuring Velero to use it as a backup destination.

---

## **1️⃣ Deployment Files Overview**

| File | Purpose |
|------|---------|
| `namespace.yaml` | Creates the `minio` namespace |
| `secret.yaml` | Stores MinIO access credentials |
| `deployment.yaml` | Deploys MinIO as an S3-compatible storage service |
| `route.yaml` | Exposes MinIO via OpenShift Route |

---

## **2️⃣ Deploying MinIO in CRC**

Run the following commands to deploy MinIO step by step:

1️⃣ **Create the namespace:**
```sh
oc apply -f namespace.yaml
```

2️⃣ **Apply the MinIO secret (already created in `secret.yaml`):**
```sh
oc apply -f secret.yaml
```

3️⃣ **Deploy MinIO:**
```sh
oc apply -f deployment.yaml
```

4️⃣ **Expose MinIO via OpenShift Route:**
```sh
oc apply -f route.yaml
```

To verify that MinIO is running:
```sh
oc get pods -n minio
```
Expected output:
```
NAME                   READY   STATUS    RESTARTS   AGE
minio-xxxx             1/1     Running   0          Xm
```

To get the MinIO Console URL:
```sh
oc get route minio-console -n minio
```

Login using:
- **Username:** `admin`
- **Password:** `supersecretpassword`

---

## **3️⃣ Configuring Velero to Use MinIO**

### ** 3.1 Configure MinIO Storage via Velero Config**
Instead of manually creating a bucket, we configure Velero to create it automatically.

```sh
velero backup-location create default \
  --provider aws \
  --bucket velero-backups \
  --secret-file <(oc get secret velero-secret -n velero -o jsonpath='{.data}' | base64 -d) \
  --config region=minio,s3ForcePathStyle=true,s3Url=http://minio-service.minio.svc:9000
```

To verify the backup storage location:
```sh
velero backup-location get
```

---

## **4️⃣ Creating and Restoring Backups with MinIO**

### ** 4.1 Backup the `db` Namespace to MinIO**
```sh
velero backup create db-backup --include-namespaces=db --storage-location=default
```

To check the backup status:
```sh
velero backup describe db-backup --details
```

### ** 4.2 Restore from MinIO Backup**
```sh
velero restore create --from-backup db-backup
```

To check the restore status:
```sh
velero restore describe db-restore --details
```

---

## **✅ Summary**
- **MinIO deployed as S3-compatible storage in CRC**
- **Velero configured to use MinIO for backups**
- **Velero automatically creates the backup bucket**
- **Full backup and restore functionality using MinIO instead of PVCs**

Now you can store and restore your CRC backups using MinIO with Velero!

---
**Copyright (c) 2025 by Alexander Kolin. All rights reserved.**