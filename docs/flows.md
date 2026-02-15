# Flows

## First-time setup

- On first run, `selfsetup.Setup()` checks for a `passfile`
  - If missing, generates a UUID passkey, writes it to `./passfile`, and prints it to stdout
  - If present, reads it into memory
- The passkey is used as the `Authorization` header for all API requests
- Two Badger KV stores are opened: `runtimes-db` and `code-db`
- Required directories (`./runtimes`, `./buildcache`) are expected to exist

## Authentication

- Every request passes through the `HandleAuth` middleware
- The middleware reads `./passfile` and compares it to the `Authorization` header
- If they don't match, the request is rejected with `401 Unauthorized`

## Create a new runtime

1. Client sends `POST /runtime/upload/:name` with a Dockerfile as multipart form data
2. Server generates a UUID filename (e.g. `abc123.dockerfile`)
3. Stores the mapping `name -> filename` in `runtimes-db`
4. Saves the Dockerfile to `./runtimes/<filename>`

## List runtimes

1. Client sends `GET /runtime/list`
2. Server returns all keys from `runtimes-db` as a JSON array of strings

## Show runtime details

1. Client sends `GET /runtime/show/:name`
2. Server looks up the filename in `runtimes-db`, reads the Dockerfile contents
3. Returns `{ "name": "...", "content": "..." }`

## Delete a runtime

1. Client sends `DELETE /runtime/delete/:name`
2. The `default` runtime cannot be deleted (returns `403`)
3. Server removes the Dockerfile from `./runtimes/` and deletes the key from `runtimes-db`

## Deploy a function (create or update)

1. Client sends `POST /function/upload` with multipart form data:
   - `name` - function name (becomes the routing endpoint)
   - `tag` - Docker image tag
   - `runtime` - name of a registered runtime
   - `port` - port the function listens on (80, 8080, 9001 are reserved)
   - `volume` - optional volume mounts (`./data:/data,./other:/other`)
   - `file` - zip archive of the function code
2. Server validates all required fields are present
3. Server checks the runtime exists in `runtimes-db`
4. **Upsert**: if a function with the same name already exists, the server tears it down first:
   - Kills the running container
   - Removes the container
   - Removes the Docker image
   - Deletes the DB entry
5. Responds immediately with `"Code uploaded, build process started"`
6. In a background goroutine:
   - Extracts the zip to a build cache directory
   - Copies the runtime Dockerfile into the extracted code directory
   - Builds the Docker image with `docker build -t <tag> .`
   - Runs the image in detached mode on the specified port
   - Stores the function config (with container ID) in `code-db`

## Get a single function

1. Client sends `GET /function/get/:name`
2. Server scans `code-db` for a function matching the name
3. Returns the full `LambdaFun` JSON object, or `404`

## List all functions

1. Client sends `GET /function/list`
2. Server returns all entries from `code-db` as `{ "keys": [...], "functions": [...] }`

## List running containers

1. Client sends `GET /function/listrunning`
2. Server queries Docker for all running containers and returns the raw Docker API response

## List installed images

1. Client sends `GET /function/listinstalled`
2. Server queries Docker for all local images and returns their repo tags as a flat string array

## Delete a function

1. Client sends `DELETE /function/delete/:name`
2. Server looks up the function by name in `code-db`
3. Kills the container, removes the container, removes the Docker image
4. Deletes the entry from `code-db`

## Start a stopped function

1. Client sends `GET /function/start/:key`
2. Server checks if the function is already running (by image tag)
3. If not running, calls `RunDetached` to start it

## Reverse proxy

Every deployed function automatically gets a subdomain. No per-function routing config needed.

1. A wildcard DNS record (`*.yourdomain.com`) points all subdomains at the server
2. A request arrives at `my-api.yourdomain.com/hello`
3. The Gin middleware (`proxy.go`) splits the `Host` header by `.` and takes the first segment (`my-api`)
4. It scans `code-db` for a function whose `Name` matches
5. If found, it proxies the request to `http://localhost:<port>/hello` using `httputil.ReverseProxy`
6. If no function matches, the request falls through to the next Gin handler

- Routing is subdomain-only. Path-based routing (e.g. `yourdomain.com/my-api`) is not supported.
- Only the first subdomain segment is matched. `v2.my-api.yourdomain.com` matches `v2`, not `my-api`.
- With `ENV=prod` and `ALLOWLIST` set, all subdomains get HTTPS automatically via Let's Encrypt.

## Background routines

### Runtime cleanup (`CleanUnusedRuntimes`)
- Runs on a timer (default 120s, configurable via `CLEANUP_INTERVAL`)
- Scans the `./runtimes/` directory
- Removes any Dockerfile that isn't referenced by a key in `runtimes-db`

### Service restart (`RestartServices`)
- Runs on a timer (default 120s, configurable via `RESTART_INTERVAL`)
- Iterates all functions in `code-db`
- If a function's container is not running, restarts it with `RunDetached`

## Auto HTTPS

- In production (`ENV=prod`), uses Let's Encrypt via `gin-gonic/autotls`
- Domain whitelist is configured via the `ALLOWLIST` env var (comma-separated)
- In development, runs plain HTTP on `:8080`
