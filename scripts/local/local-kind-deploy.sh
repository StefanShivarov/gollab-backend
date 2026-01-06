#!/bin/bash
set -euo pipefail

CLUSTER_NAME="gollab-cluster"
NAMESPACE="gollab-demo-namespace"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MIGRATIONS_PATH="$SCRIPT_DIR/../../db/migrations"
KIND_CONFIG_PATH="$SCRIPT_DIR/../../k8s/kind-config.yaml"
NAMESPACE_CONFIG_PATH="$SCRIPT_DIR/../../k8s/namespace.yaml"
POSTGRES_KUSTOMIZATION_PATH="$SCRIPT_DIR/../../k8s/postgres"
BACKEND_KUSTOMIZATION_PATH="$SCRIPT_DIR/../../k8s/backend"

# ---------------------------
# 1. Create kind cluster if not exists
# ---------------------------
if kind get clusters | grep -q "^${CLUSTER_NAME}$"; then
  echo "Kind cluster '${CLUSTER_NAME}' already exists. Skipping creation..."
else
  echo "Creating kind cluster '${CLUSTER_NAME}'..."
  kind create cluster --name "$CLUSTER_NAME" --config "$KIND_CONFIG_PATH"
fi

# ---------------------------
# 2. Create namespace if not exists
# ---------------------------
if kubectl get ns "$NAMESPACE" >/dev/null 2>&1; then
  echo "Namespace '$NAMESPACE' already exists. Skipping..."
else
  echo "Creating namespace '$NAMESPACE'..."
  kubectl apply -f "$NAMESPACE_CONFIG_PATH"
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
# 4. Apply Postgres k8s resources via Kustomize
# ---------------------------
echo "Applying Postgres resources..."
kubectl apply -k "$POSTGRES_KUSTOMIZATION_PATH"

# ---------------------------
# 5. Wait for Postgres pod
# ---------------------------
echo "Waiting for Postgres pod to be ready..."
kubectl rollout status deployment/gollab-db -n "$NAMESPACE" --timeout=120s
kubectl wait --for=condition=available deployment/gollab-db -n "$NAMESPACE" --timeout=120s

# ---------------------------
# 6. Run Flyway migrations
# ---------------------------
echo "Waiting for Postgres on localhost:5432..."
for _ in {1..60}; do
  (echo > /dev/tcp/127.0.0.1/5432) >/dev/null 2>&1 && break
  sleep 1
done

PG_URL="jdbc:postgresql://localhost:5432/gollab_db"

echo "Running Flyway migrate..."
docker run --rm \
  --network=host \
  -v "$MIGRATIONS_PATH:/flyway/sql" \
  -e FLYWAY_URL="$PG_URL" \
  -e FLYWAY_USER="${DB_USER:-postgres}" \
  -e FLYWAY_PASSWORD="${DB_PASS:-postgres}" \
  flyway/flyway:9.20.0 migrate

# ---------------------------
# 7. Apply backend resources
# ---------------------------
kubectl apply -k "$BACKEND_KUSTOMIZATION_PATH"

# ---------------------------
# 8. Wait for deployments
# ---------------------------
kubectl rollout status deployment/gollab-backend -n "$NAMESPACE" --timeout=120s
kubectl wait --for=condition=available deployment/gollab-backend -n "$NAMESPACE" --timeout=120s

# ---------------------------
# 9. Finish
# ---------------------------
echo "Backend service available at http://localhost:8080"
