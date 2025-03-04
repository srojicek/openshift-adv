# OpenShift SCC Example: Restricting Pod Access with SecurityContextConstraints

## Overview
This example demonstrates how to apply **Security Context Constraints (SCC)** in OpenShift to control pod permissions.

We define:
- A **custom SCC (`restricted-scc`)** with strict security rules.
- A **ServiceAccount (`secure-sa`)** to which the SCC is assigned.
- A **RoleBinding** to explicitly grant SCC permissions to the ServiceAccount.
- **Two pods** using `bitnami/nginx`, where:
  - One pod (`allowed-nginx`) **successfully runs** as it follows the SCC rules.
  - The other pod (`blocked-nginx`) **fails** due to violations.

---

## Steps to Deploy

### 1️⃣ **Create a ServiceAccount**
```sh
oc apply -f serviceaccount.yaml
```

### 2️⃣ **Create and Apply the SCC**
```sh
oc apply -f scc.yaml
```

### 3️⃣ **Create and Apply the RoleBinding (Required for SCC Application)**
```sh
oc apply -f rolebinding.yaml
```

### 4️⃣ **Deploy the Allowed Pod (Runs Successfully)**
```sh
oc apply -f allowed-nginx.yaml
```

### 5️⃣ **Deploy the Blocked Pod (Fails)**
```sh
oc apply -f blocked-nginx.yaml
```

---

## Expected Behavior
| Pod Name         | Status  | Reason |
|-----------------|---------|--------------------------------------------|
| `allowed-nginx` | ✅ Running | Matches SCC rules |
| `blocked-nginx` | ❌ Failed | Runs as root, escalates privileges, uses `hostPath` |

---

## Cleanup
To remove all resources:
```sh
oc delete -f serviceaccount.yaml
oc delete -f scc.yaml
oc delete -f rolebinding.yaml
oc delete -f allowed-nginx.yaml
oc delete -f blocked-nginx.yaml
```

---
**📢 Important:**  
If the `blocked-nginx` pod is still running after applying the SCC, verify that the ServiceAccount is correctly assigned and that the SCC annotation appears on the pod:

```sh
oc get pod blocked-nginx -o jsonpath='{.metadata.annotations.openshift\.io/scc}'
```

If no SCC is applied, try deleting and recreating the pod:

```sh
oc delete pod blocked-nginx
oc apply -f blocked-nginx.yaml
```

---
**© 2025 by Alexander Kolin. All rights reserved.**

