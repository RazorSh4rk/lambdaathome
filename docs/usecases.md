# Use Cases

## Runtime management

- List all registered runtimes
- Upload a new runtime (a Dockerfile that defines the base environment)
- View the Dockerfile contents of a runtime
- Delete a runtime (except the `default` runtime)

## Function deployment

- Deploy a new function by uploading a zip of your code, selecting a runtime, and specifying a port
- Update an existing function by deploying with the same name (automatic teardown + redeploy)
- View details of a deployed function (name, tag, runtime, port, volume, container ID)

## Function lifecycle

- List all deployed functions
- List all currently running containers
- List all installed Docker images (base images + function images)
- Start a stopped function by image key
- Delete a function (stops container, removes container, removes image, cleans DB)

## Routing / reverse proxy

- Access deployed functions via subdomain routing
  - e.g. `myfunction.yourdomain.com` routes to the container running on the function's port
- Each function gets its own isolated container and port

## Dashboard (web UI)

- **Key** - enter and store the API passkey for authenticating dashboard requests
- **Installed** - overview of all runtimes, function images, base images, and running containers with inline delete/stop controls
- **Runtimes** - create, view, and delete runtimes
- **Functions** - list, inspect, and delete deployed functions
- **Deploy** - upload and deploy a new function (or update an existing one)

## Example scenarios

### Python HTTP server

1. Create a `python-3` runtime with a Dockerfile based on `python:3-alpine`
2. Write a simple HTTP server (e.g. Flask, http.server)
3. Zip the code directory
4. Deploy with name=`my-api`, tag=`my-api:latest`, runtime=`python-3`, port=`9002`
5. Access at `http://localhost:9002` or via subdomain `my-api.yourdomain.com`

### Node.js Express app

1. Create a `node-20` runtime with a Dockerfile based on `node:20-alpine`
2. Write an Express app with `package.json`
3. Zip the project
4. Deploy with name=`my-node-app`, tag=`my-node-app:latest`, runtime=`node-20`, port=`9003`

### Static file server

1. Use any runtime with a web server (nginx, python http.server, etc.)
2. Bundle your static files
3. Deploy with an appropriate port
4. Optionally use a volume mount to persist data: `./uploads:/app/uploads`

### Updating a running function

1. Make code changes locally
2. Re-zip the code
3. Deploy with the **same name** - the platform automatically kills the old container, removes the old image, and deploys the new version
