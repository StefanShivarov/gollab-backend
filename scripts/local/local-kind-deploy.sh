#!/bin/bash
set -euo pipefail

CLUSTER_NAME="gollab-cluster"
NAMESPACE="gollab-demo-namespace"
KIND_CONFIG_PATH="../../k8s/overlays/kind/kind-config.yaml"
NAMESPACE_CONFIG_PATH="../../k8s/base/namespace.yaml"
POSTGRES_KUSTOMIZATION_PATH="../../k8s/overlays/kind/postgres"
BACKEND_KUSTOMIZATION_PATH="../../k8s/overlays/kind/backend"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MIGRATIONS_PATH="$SCRIPT_DIR/../../db/migrations"

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
kubectl wait --for=condition=ready pod -l app=gollab-postgres -n "$NAMESPACE" --timeout=120s

# ---------------------------
# 6. Run Flyway migrations
# ---------------------------
echo "Running Flyway migrations..."

kubectl port-forward svc/gollab-postgres-service 5432:5432 -n "$NAMESPACE" &
PF_PID=$!

cleanup() {
  kill $PF_PID || true
}
trap cleanup EXIT

for _ in {1..30}; do
  if nc -z localhost 5432; then
    echo "Postgres is reachable"
    break
  fi
  sleep 1
done

if ! nc -z localhost 5432; then
  echo "ERROR: Postgres did not become reachable via port-forward"
  exit 1
fi

docker run --rm \
  -v "$MIGRATIONS_PATH":/flyway/sql \
  -e FLYWAY_URL=jdbc:postgresql://localhost:5432/gollab_db \
  -e FLYWAY_USER="${DB_USER:-postgres}" \
  -e FLYWAY_PASSWORD="${DB_PASS:-postgres}" \
  -e FLYWAY_CONNECT_RETRIES=10 \
  flyway/flyway:9.20.0 \
  sh -c "flyway validate && flyway migrate"

# ---------------------------
# 7. Apply backend resources
# ---------------------------
kubectl apply -k "$BACKEND_KUSTOMIZATION_PATH"

# ---------------------------
# 8. Wait for deployments
# ---------------------------
DEPLOYMENTS=("gollab-backend" "gollab-db")
for dep in "${DEPLOYMENTS[@]}"; do
  echo "Waiting for deployment '$dep' to be ready..."
  kubectl rollout status deployment/$dep -n "$NAMESPACE" --timeout=120s
  kubectl wait --for=condition=available deployment/$dep -n "$NAMESPACE" --timeout=120s
done

# ---------------------------
# 9. Finish
# ---------------------------
echo "Backend service available at http://localhost:8080"
