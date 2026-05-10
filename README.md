# go-portfolio

Go web application with a fully automated CI/CD pipeline deployed on Amazon EKS using GitHub Actions, Docker, Helm, and ArgoCD.

## Overview

A static portfolio website served by a Go HTTP file server. The application is containerized with Docker and deployed to a private EKS cluster via a GitOps workflow: GitHub Actions builds and pushes a Docker image, updates the Helm chart tag in the `go-app-devops` repo, and ArgoCD automatically syncs the new version to the cluster.

## Architecture

```
Developer → git push → GitHub Actions CI Pipeline
                           │
                    ┌──────┴──────────────────┐
                    │  1. Build & Test (Go)    │
                    │  2. Lint (golangci-lint) │
                    │  3. Docker Build & Push  │
                    │  4. CVE Scan (Scout)     │
                    │  5. Update values.yaml   │
                    └──────┬──────────────────┘
                           │ push image tag → go-app-devops/charts/values.yaml
                           ▼
                     ArgoCD (GitOps)
                           │ detects change, auto-sync
                           ▼
                     EKS Cluster (go-app namespace)
                     └── Deployment (2 replicas)
                     └── Service (ClusterIP :80)
                     └── Ingress (Nginx → :8080)
```

## Application

| Property | Value |
|---|---|
| Language | Go 1.21 |
| Port | 8080 |
| Serves | `./static/index.html` |
| Docker base | `scratch` (minimal, no OS layer) |
| Image | `avian19/go-port:<github.run_id>` |

## Repository Structure

```
go-portfolio/
├── main.go                         # HTTP file server
├── main_test.go                    # Unit tests
├── Dockerfile                      # Multi-stage build (builder + scratch)
├── go.mod
├── static/
│   └── index.html                  # Portfolio website
└── .github/
    └── workflows/
        └── go-app.yaml             # CI/CD pipeline
```

## CI/CD Pipeline (`.github/workflows/go-app.yaml`)

Triggered on push to `main` and pull requests targeting `main`.

| Job | Description |
|---|---|
| `build-and-test` | `go build` + `go test ./...` |
| `lint` | golangci-lint (needs: build-and-test) |
| `docker-build-push-scan` | Build → push DockerHub → Docker Scout CVE scan (needs: lint) |
| `update-helm-repo` | Update `tag` in `go-app-devops/charts/values.yaml` via PAT (needs: docker, main only) |

## Running Locally

```bash
# Clone
git clone https://github.com/cs365-project/go-portfolio
cd go-portfolio

# Run
go run main.go
# Visit http://localhost:8080

# Test
go test ./... -v

# Docker
docker build -t go-portfolio .
docker run -p 8080:8080 go-portfolio
```

## Required GitHub Secrets

| Secret | Description |
|---|---|
| `DOCKER_USERNAME` | Docker Hub username |
| `DOCKER_PASSWORD` | Docker Hub access token |
| `REPO_PAT` | GitHub Personal Access Token (write access to go-app-devops) |

## Branch Strategy

| Branch | Purpose |
|---|---|
| `main` | Production — triggers full CI + Helm update |
| `develop` | Integration — feature branches merge here |
| `feat/*` | Feature development |

## Infrastructure

Infrastructure (EKS, VPC, IAM, Jump Server) is managed in the companion repo:
[go-app-devops](https://github.com/cs365-project/go-app-devops)
