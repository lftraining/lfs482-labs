# Lab 5: Using the Workload API with go-spiffe

## Prerequesites

- A AMD, Intel, or ARM 64-bit Linux environment.
- Familiarity with the `kubectl` command-line interface (CLI) is helpful.

## Introduction

Coastal Containers Ltd., a leading container shipping company, is looking to evolve their security practices by implementing [mutual authentication](https://en.wikipedia.org/wiki/Mutual_authentication) (via mTLS) between their Port Authority `server` and incoming vessel `clients`. In order to do so, they have decided to leverage [go-spiffe](https://github.com/spiffe/go-spiffe), a popular and convenient [Golang](https://github.com/golang/go) library built to work with [SPIFFE](https://github.com/spiffe/spiffe). By leveraging the `go-spiffe` library, Coastal Containers hopes to create SPIFFE-aware workloads and establish mTLS connections with assigned SPIFFE IDs. Ultimately, the Port Authority `server` and vessel `client` workloads should communicate securely by validating each other's SPIFFE ID's.

This hands-on exercise is based on the [go-spiffe tls example](https://github.com/spiffe/go-spiffe/tree/main/v2/examples/spiffe-tls).

### Learning Objectives

- Understand how to implement SPIFFE-aware workloads using the `go-spiffe` library.
- Understand the `go-spiffe` listening and dialoging logic for SPIFFE-aware mTLS communication.
- Learn how `go-spiffe` interacts with the Workload API to manage SVIDs.

## Preparing your Environment

Before you cast off, prepare your ships to sail by setting your working directory in (../lab-05-go-spiffe/)(../lab--05-go-spiffe/) as an environment variable:

```bash
export LAB_DIR=$(pwd)
```

This will make issuing commands easier in the following steps of this exercise, and will reduce the possibility of reference errors.

## Step-by-Step Instructions

### Step 1: Provision Infrastructure

To set sail, spin up the Kubernetes demo cluster using Kind by issusing the following `make` command while in the root [lab directory](/):

```bash
make cluster-up
```

If your cluster is already running, you can skip this step and continue on with the lab.

### Step 2: Deploying the SPIRE Setup

Initialize the SPIRE setup to your Kubernetes cluster and create the initial node registration entry: 

```shell
make deploy-spire
```

Once finished, you should see the output:

```log
SPIRE deployed n the cluster.
```

This command will spin-up the `spire-server`, `spiffe-csi-driver`, and `spire-agent` in the `spire` namespace. Additionally, this command will also create the initial node registration entry. To view the running SPIRE set-up on your kubernetes cluster, run:

```shell
kubectl get pods -n spire
```

After issuing the `kubectl` command, you should see your running SPIRE setup in an output similar to:

```log
NAME                      READY   STATUS    RESTARTS   AGE
spiffe-csi-driver-xpp4c   2/2     Running   0          72s
spire-agent-fjxw5         1/1     Running   0          67s
spire-server-0            1/1     Running   0          86s
```

It is important to note that the node registration entry initializes the `spiffe://coastal-containers.example/agent/spire-agent` SPIFFE ID, which is later used as the `parentID` for the client/server workload registration entries.

### Step 3: Explore the Server Workload

First, navigate to the server workload directory:

```shell
cd ${LAB_DIR}/workload/server
```

Here, you can see the `go-spiffe` logic within the [main.go](./workload/server/main.go) file that the server uses to interact with the Workload API and listen for incoming client connections, only accepting client(s) that present a valid X509-SVID with a matching SPIFFE ID.

Configuration items of note include:

**Function Invocation:**

- The server uses the `spiffetls.Listen` function to start listening for incoming connections as shown below.

```go
listener, err := spiffetls.Listen(ctx, "tcp", serverAddress, tlsconfig.AuthorizeID(clientID))
```

**Context (ctx):**

- `ctx` is a `context.Context` instance which blocks until the first Workload API response is received or this context is cancelled / timed out.

**Server Address:**

- `serverAddress` constant specifies where the server will listen for client connections (`0.0.0.0:8443`).

**Authentication & Authorization:**

- `clientID` is the SPIFFE ID which is used to authenticate incoming client connections (`spiffe://coastal-containers.example/workload/client`). 

- `tlsconfig.AuthorizeID(clientID)` ensures the server accepts connections only from clients that present an X509-SVID with a matching SPIFFE ID (`clientID`).

**SPIFFE_ENDPOINT_SOCKET:**

- The `spiffetls.Listen` function uses the `SPIFFE_ENDPOINT_SOCKET` environment variable to locate the Workload API address, obtaining the SVIDs needed for establishing secure communication. This environment variable is set within the [deploy-server.yaml](./config/deploy-server.yaml) manifest.

📝Note: Detailed explanations about the underlying logic are provided in the [go-spiffe tls example](https://github.com/spiffe/go-spiffe/tree/main/v2/examples/spiffe-tls) and can be found within the associated [API documentation](https://pkg.go.dev/github.com/spiffe/go-spiffe/v2/spiffetls#section-documentation).

### Step 4: Explore the Client Workload

First, navigate to the client workload directory:

```shell
cd ${LAB_DIR}/workload/client
```

Here, you can see the `go-spiffe` logic within the [main.go](./workload/client/main.go) file that the client uses to dial and establish a connection with the server, only accepting server(s) that present a valid X509-SVID with a matching SPIFFE ID.

Configuration items of note include: 

**Function Invocation:**

- The client uses the `spiffetls.Dial` function to establish a connection with the server.

```go
listener, err := spiffetls.Dial(ctx, "tcp", serverAddress, tlsconfig.AuthorizeID(serverID))
```

**Context (ctx):**

- Similar to the server, `ctx` is a `context.Context` instance that will block until it receives the first Workload API response or the context is cancelled / times out.

**Server Address:**

- `serverAddress` constant is the address of the server (`server:443`) to which the client is connectng.

- This is set to the `server` service name to work in a containerized Kubernetes environment.

**Authentication & Authorization:**

- `serverID` is the SPIFFE ID which is used to authenticate the server (`spiffe://coastal-containers.example/workload/server`). 

- `tlsconfig.AuthorizeID(serverID)` ensures the client establishes a connection only with a server that presents a X509-SVID with the expected SPIFFE ID (`serverID`).

**SPIFFE_ENDPOINT_SOCKET:**

- The `spiffetls.Dial` function uses the `SPIFFE_ENDPOINT_SOCKET` environment variable to locate the Workload API address, obtaining the SVIDs needed for establishing secure communication to the server. This environment variable is set within the [deploy-client.yaml](./config/deploy-client.yaml) manifest.

📝Note: Detailed explanations about the underlying logic are provided in the [go-spiffe tls example](https://github.com/spiffe/go-spiffe/tree/main/v2/examples/spiffe-tls) and can be found within the associated [API documentation](https://pkg.go.dev/github.com/spiffe/go-spiffe/v2/spiffetls#section-documentation).

That's all! By leveraging the `SPIFFE_ENDPOINT_SOCKET` environment variable, which can be set within your Kubernetes deployment manifests, your application can utilize the Workload API without the need of hardcoding the socket path. The `go-spiffe` library will then take care of the rest as it manages the automatic fetching and renewing of your X509-SVIDs, thus simplifying the setup of mutual authentication and secure communication between your workloads. 

### Step 5: Build and Load Workload Images

First, navigate to the root [lab-05-go-spiffe](../lab-05-go-spiffe/) directory:

```shell
cd ${LAB_DIR}
```

Next, build and load the `server` workload image by running the make command:

```shell
make cluster-build-load-image DIR=workload/server
```

After this, build and load the `client` workload image by running the make command:

```shell
make cluster-build-load-image DIR=workload/client
```

Keep in mind that this will create docker images for the server and client workloads with the `workload/server:latest` and `workload/client:latest` tags. It will also automatically load the created images into your kind cluster for you. You can check the created docker images using the command:

```shell
docker images
```

### Step 6: Create Workload Registration Entries

Before deploying the workloads, register them within the `spire-server`.

For the `server` workload, run:

```shell
kubectl exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry create \
		-spiffeID spiffe://coastal-containers.example/workload/server \
		-parentID spiffe://coastal-containers.example/agent/spire-agent \
		-selector k8s:ns:default \
		-selector k8s:sa:server
```

For the `client` workload, run:

```shell
kubectl exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry create \
		-spiffeID spiffe://coastal-containers.example/workload/client \
		-parentID spiffe://coastal-containers.example/agent/spire-agent \
		-selector k8s:ns:default \
		-selector k8s:sa:client
```

If the commands are successful, they should output the created registration entry for both workloads. 

### Step 7: Deploy the Server Workload

While still within the root [lab-05-go-spiffe](../lab-05-go-spiffe/) directory, deploy the server workload to your kind cluster:

```shell
make deploy-server
```

This will apply the [deploy-server.yaml](./config/deploy-server.yaml) manifest which creates a service, service account, and deployment for the `server` workload. If everything was successful, you should see the running `server` workload when you run:

```shell
kubectl get pods
```

### Step 8: Deploy the Client Workload

While still within the root [lab-05-go-spiffe](../lab-05-go-spiffe/) directory, deploy the client workload to your kind cluster:

```shell
make deploy-client
```

This will apply the [deploy-client.yaml](./config/deploy-client.yaml) manifest which creates a service account and deployment for the `client` workload. If everything was successful, you should see the running (or completed) `client` workload when you run:

```shell
kubectl get pods
```

⚠️ Note: The `client` pod will likely show the status of `Completed` as due to the nature of containers, they are meant to run a process or task and quit thereafter. In this case, the `client` will restart a number of times as it successfully sends the intended message to the `server` workload, receives a reply, and exits. This is expected behavior.

### Step 9: Observe Client-Server Handshake

Now, observe the logs for the `server` and `client`, ensuring the `client` sends the intended message, and the `server` responds back.

```shell
kubectl logs -f deployments/server
kubectl logs -f deployments/client
```

If everything ran successfully, you should see for the `server` logs:

```shell
SPIFFE_ENDPOINT_SOCKET: unix:///spire-agent-socket/agent.sock
Starting server on 0.0.0.0:8443
Server listening on 0.0.0.0:8443
Incoming vessel says: "This is SS Coastal Carrier hailing the port authority for Coastal Containers Ltd.\n"
```

And for the `client` logs:

```shell
SPIFFE_ENDPOINT_SOCKET: unix:///spire-agent-socket/agent.sock
Connecting to server:443
Client connected to server:443
Port Authority says: "Request received SS Coastal Carrier. You are cleared to dock.\n"
```

### Step 10: Cleanup

As the following lab exercises will use the same cluster, tear down the SPIRE setup and workload deployments from this lab by running:

```shell
cd $LAB_DIR && make tear-down
```

To tear down the entire Kind cluster, run:

```shell
cd $LAB_DIR && make cluster-down
```

## Conclusion

Congratulations! You've successfully implemented SPIFFE-aware workloads using the `go-spiffe` library, enabling them to communicate securely using SPIFFE IDs and mutual authentication. Feel free to explore further by setting up different workloads or SPIFFE ID setups, observing how SPIFFE secures communication in dynamic and containerized environments.

You are highly encouraged to explore the [go-spiffe API documentation](https://pkg.go.dev/github.com/spiffe/go-spiffe/v2) and the rich set of standalone examples provided within the [go-spiffe examples](https://github.com/spiffe/go-spiffe/tree/main/v2/examples) repository which showcases different use-cases for `go-spiffe`. Happy coding!
