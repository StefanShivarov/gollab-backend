# gollab-backend

This repository contains the backend service for the **Gollab** application. The project is developed as a course project for both the **Go course** and the **DevOps course** at Faculty of Mathematics and Informatics (FMI).

# DevOps Overview

This document describes the DevOps setup for the Golang Backend API with a PostgreSQL database. It covers the CI/CD workflows, Kubernetes deployment, and instructions for local development using **kind**.


## Project Architecture

* **Backend**: Golang REST API
* **Database**: PostgreSQL 15
* **Containerization**: Docker
* **Kubernetes**: `kind` (Kubernetes-in-Docker)
* **CI/CD**: GitHub Actions

**Main CI/CD components and flows:**

| Component               | Purpose                                                                                    |
| ----------------------- | ------------------------------------------------------------------------------------------ |
| CI Workflow             | Build, unit tests, linting, Docker image build                                             |
| CD Workflow             | Push image to GHCR, deploy to ephemeral Kubernetes cluster, run migrations and smoke tests |
| Local Deployment Script | Deploy full stack locally on a persistent kind cluster                                     |



## CI workflow

The CI workflow runs automatically on every push to the `main` branch and on every push inside a pull request pointing to the `main` branch.

### CI Jobs

#### 1. Build & Unit Tests

* Runs inside a `golang:1.25` container.
* Steps:

  * Checkout repository code
  * Build the Go application
  * Run unit tests
  * Upload build artifacts for later jobs

#### 2. Linting

* Runs after a successful build and unit tests.
* Uses **Super-Linter** to validate code quality.
* Lint errors are reported but do not fail the workflow (`DISABLE_ERRORS=true`).

#### 3. Docker Image Build

* Uses the build artifacts from previous jobs.
* Logs into GitHub Container Registry (GHCR).
* Builds the backend Docker image.
* The image is **not pushed** in CI (pushing is handled by CD).


## CD workflow

The CD workflow runs on:

* Pushes to the `main` branch
* Manual trigger (`workflow_dispatch`)

### CD Jobs

#### 1. Build and Push Docker Image

* Builds the Docker image from the current commit.
* Pushes the image to GHCR with two tags:

  * `latest`
  * Commit SHA

#### 2. Kubernetes Deployment (Ephemeral)

The deployment uses a **temporary kind cluster** created inside the GitHub Actions runner. The cluster exists only for the duration of the workflow.

##### Deployment Steps

1. **Create kind cluster**

   * Uses a custom kind configuration
   * Exposes:

     * Postgres on `localhost:5432`
     * Backend API on `localhost:8080`

2. **Create namespace**

   * All resources are deployed in `gollab-demo-namespace`

3. **Create Postgres credentials secret**

   * Injects database credentials securely into Kubernetes

4. **Deploy Postgres**

   * Applies PVC, Deployment, and Service using Kustomize
   * Waits for Postgres pod readiness

5. **Run Flyway migrations**

   * Uses Flyway Docker image
   * Connects to Postgres via host networking
   * Applies SQL migrations from `db/migrations`

6. **Deploy backend**

   * Applies backend Deployment and Service using Kustomize
   * Waits for rollout completion and pod availability

7. **Smoke tests**

   * Calls the `/health` endpoint on the backend
   * Verifies the application is running correctly

8. **Debug on failure**

   * Prints pod status, logs, Kubernetes events, and deployment details

> ⚠️ The Kubernetes cluster is destroyed automatically after the workflow finishes.


## Kubernetes Configuration

### Namespace

* `gollab-demo-namespace`
* Isolates all backend and database resources

### PostgreSQL

* **Deployment**: Single replica PostgreSQL 15
* **PersistentVolumeClaim**: 5Gi storage
* **Service**: NodePort (`30032`)
* Credentials are provided via Kubernetes Secrets

### Backend

* **Deployment**:

  * 2 replicas
  * Rolling update strategy
  * Readiness and liveness probes on `/health`
* **Service**: NodePort (`30080`)
* Environment variables configured for database connectivity



## Local Deployment with kind

For local development and testing, the entire stack can be deployed using a local kind cluster.

### Run Local Deployment

From the root of the repository folder, you can execute the following script:
```bash
./scripts/local/local-kind-deploy.sh
```

### What the Script Does

1. Creates a kind cluster if it does not already exist
2. Creates the Kubernetes namespace
3. Creates Postgres credentials secret
4. Deploys Postgres resources
5. Waits for database readiness
6. Runs Flyway migrations
7. Deploys backend resources
8. Waits for backend deployment to become available

After successful execution, the backend API is available at:

```
http://localhost:8080
```

---

## Summary

* Fully automated CI/CD pipelines using GitHub Actions
* Docker-based builds and deployments
* Kubernetes orchestration using kind (no paid cloud services required)
* Database schema migrations handled via Flyway
* Identical deployment logic for CI/CD and local environments
