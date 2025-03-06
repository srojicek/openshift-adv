# **Velero Backup and Restore Setup for CRC with PVC Storage**

This guide outlines the deployment of **Velero in CRC**, storing backups in a **MINIO** instead of external storage. The setup ensures **backup and restore functionality**.



---

## ** Deploying Velero in CRC**
Run the following commands to deploy Velero step by step:

## ** Install Velero CLI**
To interact with Velero, you need to install the Velero CLI on your local machine.

### ** Installation on Linux/macOS**
```sh
curl -LO https://github.com/vmware-tanzu/velero/releases/download/v1.10.0/velero-v1.10.0-linux-amd64.tar.gz
tar -xvf velero-v1.10.0-linux-amd64.tar.gz
sudo mv velero-v1.10.0-linux-amd64/velero /usr/local/bin/
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

### ** Prepare ENV **
```sh
oc new-project velero
oc create sa velero      

oc adm policy add-scc-to-user privileged -z velero -n velero 
```

### ** Deploy MINIO **
```sh
oc apply -f minio.yaml
```

### ** Install Velero **
```sh
oc create secret generic cloud-credentials --namespace velero   --from-literal=AWS_ACCESS_KEY_ID=minioadmin   --from-literal=AWS_SECRET_ACCESS_KEY=minioadmin

velero install   --provider aws   --plugins velero/velero-plugin-for-aws:v1.7.1   --bucket velero-backups   --secret-file cloud-credentials   --backup-location-config region=minio,s3ForcePathStyle=true,s3Url=http://minio.velero.svc:9000   --use-node-agent   --namespace velero

oc patch deployment velero -n velero --type='json' -p '[{"op": "add", "path": "/spec/template/spec/volumes/-", "value": {"name": "cloud-credentials", "secret": {"secretName": "cloud-credentials"}}}]'
```

### ** Provide right secret **
```sh
cat <<EOF | oc apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: cloud-credentials
  namespace: velero
type: Opaque
data:
  cloud: $(echo -e "[default]\naws_access_key_id=minioadmin\naws_secret_access_key=minioadmin" | base64 -w0)
EOF
```

## ** Creating a Backup of the `db` Namespace**
Once Velero is deployed, create a backup of the `db` namespace:
```sh
velero create backup dev --include-namespaces dev
```

To check the backup status:
```sh
velero backup describe dev --details
```

To list all backups:
```sh
velero backup get
```

---

## ** Restoring the `db` Namespace from Backup**
If needed, restore the `db` namespace from the last backup:
```sh
velero restore create --from-backup dev
```

To check the restore status:
```sh
velero restore describe dev --details
```

To list all restores:
```sh
velero restore get
```

---
**Copyright (c) 2025 by Alexander Kolin. All rights reserved.**