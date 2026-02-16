# Lambda at home

Deploy serverless functions _on a server_.

---

## Table of Contents

- [Install](#install)
- [What is this](#what-is-this)
- [Lambdaathome vs Kubernetes](#lambdaathome-vs-kubernetes)
  - [Advantages over Kubernetes](#advantages-over-kubernetes)
  - [Limitations compared to Kubernetes](#limitations-compared-to-kubernetes)
  - [When to use what](#when-to-use-what)
- [Project Setup](#project-setup)
  - [Prerequisites](#prerequisites)
  - [Backend](#backend)
  - [Dashboard](#dashboard)
  - [Environment Variables](#environment-variables)
- [Authentication](#authentication)
- [API Reference](#api-reference)
  - [Runtimes](#runtimes)
  - [Functions](#functions)
- [Reverse Proxy](#reverse-proxy)
- [Testing](#testing)
  - [Unit Tests](#unit-tests)
  - [E2E Tests](#e2e-tests)
- [Architecture](#architecture)
- [Docs](#docs)

---

## What is this

Remember when Heroku was a viable thing, before Salesforce bought it and went "hm yes this platform that was specifically tailored for small scale projects with solo developers to make them easier to deploy would work way better as a very expensive kubernetes service that nobody asked for"?

Have you ever wanted to run something on AWS, but it was not written in one of the like 5 stacks they support on lambdas, or it was but you used a framework and didn't feel like rewriting the route handling?

Have you ever used GCP cloud run and realized it was actually pretty good, I really have no complaints about this one tbh.

Do you have $5 to spend on a VPS and just not worry about scaling lambda costs?

**Lambdaathome is**
  - self hosted
  - minimal setup
  - widely compatible
  - auto HTTPS (Let's Encrypt)
  - not opinionated
  - not locked to any technology
  - light on resources
  - extendable

You bring a Dockerfile and a zip of your code. Lambdaathome builds it, runs it, and reverse-proxies traffic to it.

---

## Lambdaathome vs Kubernetes

### Advantages over Kubernetes

**Resource efficiency**
- The entire platform is a single Go binary + Docker. No kubelet, kube-apiserver, etcd, kube-controller-manager, kube-scheduler, kube-proxy, or CoreDNS running in the background.
- A Kubernetes control plane consumes 500MB-2GB of RAM before any workload is deployed. Lambdaathome adds almost nothing beyond Docker itself.
- Runs comfortably on a $5/month VPS with 1GB RAM. A managed Kubernetes cluster (EKS, GKE, AKS) costs $70-200+/month for the control plane alone.

**Ease of use**
- No YAML manifests. Deploying a function is a single API call with a zip upload. No Deployments, Services, Ingress, ConfigMaps, or PersistentVolumeClaims to write and maintain.
- No kubectl, no kubeconfig, no contexts. The API is a small set of HTTP endpoints.
- Flat mental model: function = container, subdomain = function. No pods-within-nodes, Services-selecting-pods, Ingress-routing-to-Services layering.

**Quick setup**
- `go run .` starts everything. Kubernetes requires `kubeadm init`, CNI plugin installation, node joining, or cloud provider configuration for managed clusters.
- Near-zero learning curve. If you know HTTP and Docker, you know the platform.
- Upgrades are `git pull && go run .`. No multi-step control plane upgrades, no node-by-node kubelet rolling updates, no API deprecation handling.

**Auto SSL**
- Built-in Let's Encrypt via `autocert`. Set `ENV=prod` and `ALLOWLIST=your.domain.com`, done.
- Kubernetes needs cert-manager installed separately (CRDs, RBAC, webhook configs), plus Certificate and Issuer resources, plus Ingress TLS section configuration.

**Operational simplicity**
- No CNI plugin management (Calico, Flannel, Cilium). Docker bridge networking just works.
- No RBAC policy trees, no CRD version matrices, no operator lifecycle management.
- Straightforward debugging: functions are Docker containers on localhost. `docker logs`, `docker exec`, and `docker inspect` work directly. No iptables maze from kube-proxy, no overlay network to troubleshoot.
- Self-contained state: the entire platform (DB, passfile, runtimes, TLS cache) lives in the working directory. Backup is `cp -r`. Migration is `scp` to a new server.

**Smaller attack surface**
- Fewer components means fewer CVEs to track. Kubernetes has had critical CVEs in the API server, kubelet, and etcd. Lambdaathome exposes a single HTTP server.

### Limitations compared to Kubernetes

**Scaling**
- **1 replica only.** Each function runs as exactly one container. No ReplicaSets, no horizontal pod autoscaling, no vertical pod autoscaling.
- **Single server deployment.** No multi-node clustering. If the host goes down, everything goes down.
- **No load balancing.** Requests hit the single container directly via reverse proxy. No Service-level load balancing, no external LB integration.
- No auto-scaling based on CPU, memory, or custom metrics.

**Deployment strategies**
- **No rolling updates.** Redeploying a function kills the old container then starts the new one. There is a window of downtime.
- No canary deployments, no blue-green deployments, no traffic splitting.
- No rollback mechanism. Once redeployed, the old image is removed. No revision history.

**Networking**
- No service discovery between functions (no DNS-based `<service>.<namespace>.svc.cluster.local`).
- No network policies for restricting container-to-container traffic.
- No ingress controller ecosystem (no path-based routing, rate limiting, circuit breaking). Routing is subdomain-only.
- No service mesh (no mTLS between functions, no distributed tracing, no retries/timeouts).

**Security and multi-tenancy**
- No RBAC. Authentication is a single shared passkey. All authenticated users have full access.
- No secrets management (no Kubernetes Secrets, no Vault integration).
- No namespace isolation. All functions share a flat namespace.
- No pod security standards (non-root enforcement, read-only filesystem, seccomp profiles).
- No multi-tenancy support.

**Resource management**
- No CPU or memory limits on containers. A single function can consume all host resources.
- No resource quotas, no priority classes, no preemption.
- Manual port assignment. Port conflicts are the user's problem.

**Storage**
- Volumes are simple host-path bind mounts. No PersistentVolumes, no dynamic provisioning, no CSI drivers, no volume snapshots.

**Observability**
- Minimal logging (stdout or a single logfile). No structured logging, no per-function log separation, no centralized log aggregation.
- No metrics endpoint (no Prometheus, no resource usage stats).
- No distributed tracing, no audit logging.

**Configuration**
- Imperative state (API calls). No declarative manifests, no GitOps, no drift detection.
- No ConfigMaps, no CRDs, no operator pattern.
- No jobs or cron jobs for batch/scheduled workloads.

### When to use what

| | Lambdaathome | Kubernetes |
|---|---|---|
| Side projects, hobby apps | Great fit | Overkill |
| Single developer / small team | Great fit | Usually overkill |
| $5-10/month budget | Great fit | Not feasible |
| Need HA / zero-downtime deploys | Not suitable | Built for this |
| Multi-region / multi-node | Not suitable | Built for this |
| Heavy traffic / auto-scaling | Not suitable | Built for this |
| Enterprise / multi-tenant | Not suitable | Built for this |

---

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/RazorSh4rk/lambdaathome/main/install.sh | sudo bash
```

This detects your architecture (x86_64 or arm64), downloads the latest release, and installs the binary to `/usr/local/bin/lambdaathome`.

The only prerequisite is a running [Docker](https://docs.docker.com/get-docker/) daemon.

---

## Project Setup

For development, you can build from source instead:

### Prerequisites

- [Go](https://go.dev/) 1.22+
- [Docker](https://docs.docker.com/get-docker/) (daemon must be running)
- [Bun](https://bun.sh/) (for the dashboard, npm/pnpm also work)

### Backend

```bash
# clone the repo
git clone https://github.com/RazorSh4rk/lambdaathome.git
cd lambdaathome

# install go dependencies
go mod download

# run the server
go run .

# or with air for live reload
air
```

On first run, a `passfile` is generated with a UUID. **Copy this value** -- it's your API key for all requests.

### Dashboard

```bash
cd dashboard
bun install
bun run dev
```

The dashboard runs on `http://localhost:5173` by default and talks to the backend on `:8080`.

### Environment Variables

Create a `.env` file in the project root (or set them in your shell):

| Variable | Default | Description |
|---|---|---|
| `ENV` | (empty) | Set to `prod` for auto HTTPS via Let's Encrypt |
| `GIN_DEBUG` | (empty) | Set to `false` or `0` for release mode |
| `LOG_TO_FILE` | (empty) | Set to `true` or `1` to write logs to `./logfile` |
| `ALLOWLIST` | (empty) | Comma-separated domain list for HTTPS certificates |
| `CLEANUP_INTERVAL` | `120` | Seconds between runtime file cleanup sweeps |
| `RESTART_INTERVAL` | `120` | Seconds between dead container restart checks |

---

## Authentication

All API requests require the `Authorization` header set to the contents of `./passfile`.

```bash
# read the key
KEY=$(cat passfile)

# use it in requests
curl -H "Authorization: $KEY" http://localhost:8080/runtime/list
```

---

## API Reference

### Runtimes

A runtime is a Dockerfile that defines the base environment for your functions.

#### List runtimes

```bash
curl -H "Authorization: $KEY" \
  http://localhost:8080/runtime/list
```

Response: `["default", "python-3", "node-20"]`

#### Upload a runtime

```bash
curl -X POST \
  -H "Authorization: $KEY" \
  -F "file=@./Dockerfile" \
  http://localhost:8080/runtime/upload/python-3
```

Response:
```json
{
  "runtimeName": "python-3",
  "fileName": "a1b2c3d4-...-dockerfile"
}
```

#### Show runtime Dockerfile

```bash
curl -H "Authorization: $KEY" \
  http://localhost:8080/runtime/show/python-3
```

Response:
```json
{
  "name": "python-3",
  "content": "FROM python:3-alpine\nWORKDIR /app\nCOPY . .\nCMD [\"python\", \"main.py\"]"
}
```

#### Delete a runtime

```bash
curl -X DELETE \
  -H "Authorization: $KEY" \
  http://localhost:8080/runtime/delete/python-3
```

Response: `{"message": "Deleted"}`

> The `default` runtime cannot be deleted.

---

### Functions

#### Deploy a function

Upload a zip of your code along with deployment config. If a function with the same name already exists, it will be torn down and replaced (upsert).

```bash
curl -X POST \
  -H "Authorization: $KEY" \
  -F "name=my-api" \
  -F "tag=my-api:latest" \
  -F "runtime=python-3" \
  -F "port=9002" \
  -F "volume=./data:/data" \
  -F "file=@./code.zip" \
  http://localhost:8080/function/upload
```

Response:
```json
{
  "message": "Code uploaded, build process started",
  "data": {
    "name": "my-api",
    "tag": "my-api:latest",
    "runtime": "python-3",
    "port": "9002",
    "volume": "./data:/data",
    "source": "",
    "id": ""
  }
}
```

> The build and deploy happen asynchronously. The container ID and source path are populated once the build completes.

> Ports 80, 8080, and 9001 are reserved and cannot be used.

#### Get a function

```bash
curl -H "Authorization: $KEY" \
  http://localhost:8080/function/get/my-api
```

Response:
```json
{
  "name": "my-api",
  "tag": "my-api:latest",
  "runtime": "python-3",
  "port": "9002",
  "volume": "./data:/data",
  "source": "./buildcache/abc123/code/",
  "id": "d4e5f6..."
}
```

#### List all functions

```bash
curl -H "Authorization: $KEY" \
  http://localhost:8080/function/list
```

Response:
```json
{
  "keys": ["d4e5f6..."],
  "functions": [
    {
      "name": "my-api",
      "tag": "my-api:latest",
      "runtime": "python-3",
      "port": "9002",
      "volume": "",
      "source": "...",
      "id": "d4e5f6..."
    }
  ]
}
```

#### List running containers

```bash
curl -H "Authorization: $KEY" \
  http://localhost:8080/function/listrunning
```

Response: Docker API container objects (names, image, ports, state, mounts, etc.)

#### List installed Docker images

```bash
curl -H "Authorization: $KEY" \
  http://localhost:8080/function/listinstalled
```

Response: `["python:3-alpine", "my-api:latest", "node:20-alpine"]`

#### Start a stopped function

```bash
curl -H "Authorization: $KEY" \
  http://localhost:8080/function/start/my-api:latest
```

Response: `{"message": "Started"}`

#### Delete a function

Kills the container, removes it, removes the Docker image, and cleans the database entry.

```bash
curl -X DELETE \
  -H "Authorization: $KEY" \
  http://localhost:8080/function/delete/my-api
```

Response: `{"message": "Deleted"}`

#### Update a function

There is no separate update endpoint. Deploy with the same `name` to replace the existing function:

```bash
# same name, new code/config
curl -X POST \
  -H "Authorization: $KEY" \
  -F "name=my-api" \
  -F "tag=my-api:v2" \
  -F "runtime=python-3" \
  -F "port=9002" \
  -F "file=@./code-v2.zip" \
  http://localhost:8080/function/upload
```

---

## Reverse Proxy

Every deployed function automatically gets its own subdomain. Deploy a function named `my-api` and it's immediately reachable at:

```
https://my-api.yourdomain.com/any/path
```

No configuration needed per function -- the proxy matches the first subdomain segment to a function name and forwards the request to its container. Combined with auto SSL, every function gets HTTPS out of the box.

### DNS Setup

Point a wildcard DNS record at your server:

```
*.yourdomain.com  →  A  →  <your server IP>
```

That's it. Any new function you deploy is instantly accessible at `<function-name>.yourdomain.com` without touching DNS again.

### How it works

1. Request comes in for `my-api.yourdomain.com/hello`
2. The proxy extracts `my-api` from the `Host` header
3. Looks up the function in the database, finds it runs on port 9002
4. Forwards the request to `http://localhost:9002/hello`

### Caveats

- Routing is subdomain-based only. No path-based routing (e.g. `yourdomain.com/my-api` won't work).
- Only the first subdomain segment is matched. `my-api.yourdomain.com` works, `v2.my-api.yourdomain.com` would try to match `v2`, not `my-api`.
- In local development there are no real subdomains. Hit your functions directly at `http://localhost:<port>`.

---

## Testing

### Unit Tests

Go unit tests with mocked Docker client:

```bash
cd route-handlers
go test -v ./...
```

Tests cover:
- Function lookup by name (found, not found, empty DB)
- GET single function (found, not found)
- List functions (empty, with data)
- Delete function (full teardown, not found)

### E2E Tests

Shell scripts in `tests/e2e/` that exercise the full API against a running server:

```bash
# start the server first
go run . &

# run all e2e tests in order
bash tests/e2e/create_runtime.sh
bash tests/e2e/create_function.sh
bash tests/e2e/get_function.sh
bash tests/e2e/list_functions.sh
bash tests/e2e/update_function.sh
bash tests/e2e/delete_function.sh
bash tests/e2e/cleanup.sh
```

Each script reads the passkey from `./passfile` automatically.

---

## Architecture

```
lambdaathome/
├── main.go                  # entry point, route registration
├── selfsetup/               # first-run setup (passfile generation)
├── route-handlers/          # Gin HTTP handlers
│   ├── authhandler.go       # Authorization middleware
│   ├── proxy.go             # reverse proxy middleware
│   ├── runtimehandlers.go   # runtime CRUD
│   ├── codehandlers.go      # function upload/build
│   ├── lambdahandlers.go    # function CRUD + lifecycle
│   ├── docker.go            # Docker client interface + factory
│   └── lambdahandlers_test.go
├── docker-commands/         # Docker SDK wrappers
│   ├── client.go            # client constructor
│   ├── buildimage.go        # docker build
│   ├── rundetached.go       # docker run + container create
│   ├── kill.go              # kill, remove container, remove image
│   ├── listrunning.go       # list containers
│   └── listinstalled.go     # list images
├── db/                      # Badger KV store wrapper
│   ├── db.go                # CRUD operations
│   └── routines.go          # background cleanup + restart
├── types/                   # shared types
│   └── lambdafun.go         # LambdaFun struct
├── ssl/                     # auto HTTPS with Let's Encrypt
├── tests/e2e/               # end-to-end shell scripts
├── dashboard/               # SvelteKit + Skeleton UI
│   └── src/routes/          # pages: key, installed, runtimes, functions, deploy
├── runtimes/                # uploaded Dockerfile storage
├── buildcache/              # temporary build directories
└── docs/                    # flows and use cases
```

**Key dependencies:**
- [Gin](https://github.com/gin-gonic/gin) - HTTP framework
- [Badger](https://github.com/dgraph-io/badger) - embedded KV store
- [Docker SDK](https://pkg.go.dev/github.com/docker/docker/client) - container management
- [SvelteKit](https://kit.svelte.dev/) + [Skeleton](https://www.skeleton.dev/) - dashboard UI

---

## Docs

- [docs/flows.md](docs/flows.md) - detailed step-by-step flows for every operation
- [docs/usecases.md](docs/usecases.md) - use cases and example deployment scenarios
