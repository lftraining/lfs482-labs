# Lab 5: Using the Workload API with go-spiffe

## Prerequisites

- An AMD, Intel, or ARM 64-bit Linux environment.
- Familiarity with the `kubectl` command-line interface (CLI) is helpful.

## Introduction

Coastal Containers Ltd., a leading container shipping company, is looking to evolve their security practices by
implementing [mutual authentication](https://en.wikipedia.org/wiki/Mutual_authentication) (via mTLS) between their Port
Authority `server` and incoming vessel `clients`. In order to do so, they have decided to leverage
[go-spiffe](https://github.com/spiffe/go-spiffe), a popular and convenient [Golang](https://github.com/golang/go)
library built to work with [SPIFFE](https://github.com/spiffe/spiffe). By leveraging the `go-spiffe` library,
Coastal Containers hopes to create SPIFFE-aware workloads and establish mTLS connections with assigned SPIFFE IDs.
Ultimately, the Port Authority `server` and vessel `client` workloads should communicate securely by validating each
other's SPIFFE ID's.

This hands-on exercise is based on the
[go-spiffe tls example](https://github.com/spiffe/go-spiffe/tree/main/v2/examples/spiffe-tls).

### Learning Objectives

- Understand how to implement SPIFFE-aware workloads using the `go-spiffe` library.
- Understand the `go-spiffe` listening and dialoging logic for SPIFFE-aware mTLS communication.
- Learn how `go-spiffe` interacts with the Workload API to manage SVIDs.

## Step-by-Step Instructions

### Step 1: Provision Infrastructure

To set sail, spin up the Kubernetes demo cluster using Kind by issuing the following `make` command while in the root
lab directory:

```bash
make cluster-up
```

### Step 2: Deploy SPIRE to the Cluster

Initialize the SPIRE setup to your Kubernetes cluster:

```shell
make deploy-spire spire-wait-for-agent
```

The SPIRE configuration is identical to the previous lab. We will also create SPIRE registration entries for the SPIRE
agent, the server workload, and the client workload. These are the same as in previous labs so for convenience run the
following command:

```shell
make create-registration-entries
```

Once finished, you should see the output:

```shell
Entry ID         : 01de1c15-ef71-4834-b702-455d70586212
SPIFFE ID        : spiffe://coastal-containers.example/agent/spire-agent
Parent ID        : spiffe://coastal-containers.example/spire/server
Revision         : 0
X509-SVID TTL    : default
JWT-SVID TTL     : default
Selector         : k8s_psat:agent_ns:spire
Selector         : k8s_psat:agent_sa:spire-agent
Selector         : k8s_psat:cluster:kind-kind

Entry ID         : 05b21cc2-4e90-4aad-8b48-419fd9938678
SPIFFE ID        : spiffe://coastal-containers.example/workload/server
Parent ID        : spiffe://coastal-containers.example/agent/spire-agent
Revision         : 0
X509-SVID TTL    : default
JWT-SVID TTL     : default
Selector         : k8s:ns:default
Selector         : k8s:sa:server

Entry ID         : f7f36f9c-c73e-42ef-bb15-3d0f12a43f1d
SPIFFE ID        : spiffe://coastal-containers.example/workload/client
Parent ID        : spiffe://coastal-containers.example/agent/spire-agent
Revision         : 0
X509-SVID TTL    : default
JWT-SVID TTL     : default
Selector         : k8s:ns:default
Selector         : k8s:sa:client
```

### Step 3: Explore the Server Workload

You can see the `go-spiffe` logic within the [main.go](server/main.go) file that the server uses to interact
with the Workload API and listen for incoming client connections, only accepting client(s) that present a valid
X509-SVID with a matching SPIFFE ID.

Configuration items of note include:

**Function Invocation:**

- The server uses the `spiffetls.Listen` function to start listening for incoming connections as shown below.

```go
listener, err := spiffetls.Listen(ctx, "tcp", serverAddress, tlsconfig.AuthorizeID(clientID))
```

**Server Address:**

- `serverAddress` constant specifies where the server will listen for client connections (`0.0.0.0:8443`).

**Authentication & Authorization:**

- `clientID` is the SPIFFE ID which is used to authenticate incoming client connections
(`spiffe://coastal-containers.example/workload/client`).

- `tlsconfig.AuthorizeID(clientID)` ensures the server accepts connections only from clients that present an X509-SVID
with a matching SPIFFE ID (`clientID`).

**SPIFFE_ENDPOINT_SOCKET:**

- The `spiffetls.Listen` function uses the `SPIFFE_ENDPOINT_SOCKET` environment variable to locate the Workload API
address, obtaining the SVIDs needed for establishing secure communication. This environment variable is set within the
[deploy-server.yaml](spire-server/server.yaml) manifest.

📝Note: Detailed explanations about the underlying logic are provided in the
[go-spiffe tls example](https://github.com/spiffe/go-spiffe/tree/main/v2/examples/spiffe-tls) and can be found within
the associated [API documentation](https://pkg.go.dev/github.com/spiffe/go-spiffe/v2/spiffetls#section-documentation).

### Step 4: Explore the Client Workload

You can see the `go-spiffe` logic within the [main.go](client/main.go) file that the client uses to dial and
establish a connection with the server, only accepting server(s) that present a valid X509-SVID with a matching SPIFFE
ID.

Configuration items of note include:

**Function Invocation:**

- The client uses the `spiffetls.Dial` function to establish a connection with the server.

```go
listener, err := spiffetls.Dial(ctx, "tcp", serverAddress, tlsconfig.AuthorizeID(serverID))
```

**Server Address:**

- `serverAddress` constant is the address of the server (`server:443`) to which the client is connecting.
- This is set to the `server` service name to work in a containerized Kubernetes environment.

**Authentication & Authorization:**

- `serverID` is the SPIFFE ID which is used to authenticate the server
(`spiffe://coastal-containers.example/workload/server`).

- `tlsconfig.AuthorizeID(serverID)` ensures the client establishes a connection only with a server that presents a
X509-SVID with the expected SPIFFE ID (`serverID`).

**SPIFFE_ENDPOINT_SOCKET:**

- The `spiffetls.Dial` function uses the `SPIFFE_ENDPOINT_SOCKET` environment variable to locate the Workload API
address, obtaining the SVIDs needed for establishing secure communication to the server. This environment variable is
set within the [client](client/app.yaml) manifest.

📝Note: Detailed explanations about the underlying logic are provided in the
[go-spiffe tls example](https://github.com/spiffe/go-spiffe/tree/main/v2/examples/spiffe-tls) and can be found within
the associated [API documentation](https://pkg.go.dev/github.com/spiffe/go-spiffe/v2/spiffetls#section-documentation).

That's all! By leveraging the `SPIFFE_ENDPOINT_SOCKET` environment variable, which can be set within your Kubernetes
deployment manifests, your application can utilize the Workload API without the need of hard coding the socket path.
The `go-spiffe` library will then take care of the rest as it manages the automatic fetching and renewing of your
X509-SVIDs, thus simplifying the setup of mutual authentication and secure communication between your workloads.

### Step 5: Build and Load Workload Images

Build and load the `server` and `client` workload image by running the make command:

```shell
make workload-images
```

### Step 6: Deploy the Workloads

While still within the root [lab-05-go-spiffe](../lab-05-go-spiffe) directory, deploy the server workload to your kind
cluster:

```shell
make deploy-workloads
```

Check the workloads by running the following command:

```shell
kubectl get pods
```
If everything was successful, you should see the running `server` workload and see the running (or completed) `client`
workload.

⚠️ Note: The `client` pod will likely show the status of `Completed` as due to the nature of containers, they are meant
to run a process or task and quit thereafter. In this case, the `client` will restart a number of times as it
successfully sends the intended message to the `server` workload, receives a reply, and exits. This is expected behavior.

### Step 7: Observe Client-Server Handshake

Now, observe the logs for the `server` and `client`, ensuring the `client` sends the intended message, and the `server`
responds back.

```shell
kubectl logs deployments/server
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
kubectl logs deployments/client
```

```shell
SPIFFE_ENDPOINT_SOCKET: unix:///spire-agent-socket/agent.sock
Connecting to server:443
Client connected to server:443
Port Authority says: "Request received SS Coastal Carrier. You are cleared to dock.\n"
```

### Step 8: Cleanup

To tear down the Kind cluster, run:

```shell
make cluster-down
```

## Conclusion

Congratulations! You've successfully implemented SPIFFE-aware workloads using the `go-spiffe` library, enabling them to
communicate securely using SPIFFE IDs and mutual authentication. Feel free to explore further by setting up different
workloads or SPIFFE ID setups, observing how SPIFFE secures communication in dynamic and containerized environments.

You are highly encouraged to explore the
[go-spiffe API documentation](https://pkg.go.dev/github.com/spiffe/go-spiffe/v2) and the rich set of standalone examples
provided within the [go-spiffe examples](https://github.com/spiffe/go-spiffe/tree/main/v2/examples) repository which
showcases different use-cases for `go-spiffe`. Happy coding!
