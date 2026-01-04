#!/bin/bash
set -euo pipefail

CLUSTER_NAME="gollab-cluster"
NAMESPACE="gollab-demo-namespace"

# ---------------------------
# 1. Create kind cluster if not exists
# ---------------------------
if kind get clusters | grep -q "^${CLUSTER_NAME}$"; then
  echo "Kind cluster '${CLUSTER_NAME}' already exists. Skipping creation..."
else
  echo "Creating kind cluster '${CLUSTER_NAME}'..."
  kind create cluster --name "$CLUSTER_NAME" --config k8s/overlays/kind/kind-config.yaml
fi

# ---------------------------
# 2. Create namespace if not exists
# ---------------------------
if kubectl get ns "$NAMESPACE" >/dev/null 2>&1; then
  echo "Namespace '$NAMESPACE' already exists. Skipping..."
else
  echo "Creating namespace '$NAMESPACE'..."
  kubectl apply -f k8s/base/namespace.yaml
fi

# ---------------------------
# 3. Create Postgres secret (always apply)
# ---------------------------
kubectl create secret generic postgres-credentials \
  --from-literal=POSTGRES_USER="${DB_USER:-postgres}" \
  --from-literal=POSTGRES_PASSWORD="${DB_PASS:-postgres}" \
  --from-literal=POSTGRES_DB="gollab_db" \
  -n "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# ---------------------------
# 4. Apply k8s resources via Kustomize
# ---------------------------
echo "Applying Kubernetes resources..."
kubectl apply -k ../../k8s/overlays/kind

# ---------------------------
# 5. Wait for deployments
# ---------------------------
DEPLOYMENTS=("gollab-backend" "gollab-db")
for dep in "${DEPLOYMENTS[@]}"; do
  echo "Waiting for deployment '$dep' to be ready..."
  kubectl rollout status deployment/$dep -n "$NAMESPACE" --timeout=120s
  kubectl wait --for=condition=available deployment/$dep -n "$NAMESPACE" --timeout=120s
done

# ---------------------------
# 6. Finish
# ---------------------------
echo "Backend service available at http://localhost:8080"
