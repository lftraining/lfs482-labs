# Lab 4: Getting SVIDS with SPIFFE-Helper

## Prerequisites

- An AMD, Intel, or Arm 64-bit Linux environment.
- Familiarity with `docker` and `kubectl` command-line interface (CLI) commands is helpful.

## Introduction

Coastal Containers Ltd. is a leading container shipping company with a legacy manifest file server system. As part of
their digital transformation journey, they recognized the need to improve the security of their system. To address this,
they are adopting a Zero Trust security model, leveraging SPIFFE/SPIRE. This approach allows them to identify and
authenticate workloads in heterogeneous environments securely.

The challenge? Securing their applications using X.509 SVIDs without a large scale refactoring effort. This lab focuses
on using the [spiffe-helper](https://github.com/spiffe/spiffe-helper) utility to retrieve and manage SVIDs on behalf of
the applications without having to refactor the applications.

### Learning Objectives

- Understand the role and utility of [spiffe-helper](https://github.com/spiffe/spiffe-helper).
- Learn how to retrieve and rotate SVIDs for applications that are not SPIFFE aware.

This will make issuing commands easier in the following steps of this exercise, and will reduce the possibility of
reference errors.

## Step-by-Step Instructions

### Step 1: Provision Infrastructure

To set sail, spin up your demo Kubernetes cluster using [`Kind`](https://kind.sigs.k8s.io/) by issuing the following 
`make` command while in the root [lab directory](.):

```bash
make cluster-up
```

### Step 2: Deploying the SPIRE Setup

Initialize the SPIRE setup on your Kubernetes cluster and wait for the SPIRE agent to be up and running:

```shell
make deploy-spire spire-wait-for-agent
```

The SPIRE configuration for the server and the agent are identical with the addition of the 
[SPIFFE CSI Driver](https://github.com/spiffe/spiffe-csi). 

From the 
[SPIFFE Workload Endpoint](https://github.com/spiffe/spiffe/blob/main/standards/SPIFFE_Workload_Endpoint.md#3-transport)
the workload API *must* be served over gRPC, and *should* prefer Unix Domain Socket transport. The SPIRE agents are
deployed as a DaemonSet, so there will be a single instance on every Node in the cluster, and expose the Workload API
over a Unix Domain Socket using a `hostPath` volume i.e. mounting a directory from the host Node's filesystem. For
workloads to be able to access the Workload API, they would require permissions to mount `hostPath` volumes which
presents a number of security risks.

To address this problem, the SPIFFE project has created a 
[Container Storage Interface (CSI)](https://kubernetes-csi.github.io/docs/) driver for Kubernetes that facilitates 
injection of the [SPIFFE Workload API](https://github.com/spiffe/spiffe/blob/main/standards/SPIFFE_Workload_API.md). 
This is done by mounting a directory containing the Workload API Unix domain socket as an
[ephemeral inline volume](https://kubernetes.io/blog/2020/01/21/csi-ephemeral-inline-volumes/) into your workload pods.

This means that the SPIRE agent requires the necessary permissions to mount `hostPath` volumes, and therefore deploying
these would normally be done as part of platform provisioning. The workloads however, do not need these elevated
permissions.

Typically, the CSI driver will share the directory hosting the Workload API with a read-only bind mount into the
container at the requested target path. This begs the question of what an
[ephemeral inline volume](https://kubernetes.io/blog/2020/01/21/csi-ephemeral-inline-volumes/) is, and why is it useful.
Traditionally, volumes provided by external storage drivers in K8s are persistent by nature. This means that they
operate on a completely independent lifecycle to the pods (e.g., persisting even if a pod is killed). In use cases where
data volumes are required to be tied directly to the lifecycle of a pod, they need to be created and deleted with the
pod (aka ephemeral). To achieve this, these volumes are defined as part of the pod spec itself (aka inline).

Circling back to the `spiffe-csi-driver`, this means that when the pod is destroyed, the driver is invoked and removes
the bind mount.

### Step 3: Create Initial Node Registration Entry

Using your SPIRE tools, create the initial node registration entry. This step ensures your SPIRE server recognizes nodes
in your cluster.

```shell
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://coastal-containers.example/agent/spire-agent \
    -selector k8s_psat:cluster:kind-kind \
    -selector k8s_psat:agent_ns:spire \
    -selector k8s_psat:agent_sa:spire-agent \
    -node
```

If the command is successful, you should see an SVID displayed in your terminal.

### Step 4: Create Registration Entries for Server and Client Workloads

To identify and authenticate your workloads, exec into the running `spire-server` and create registration entries.

First, create one for the manifest server:

```shell
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://coastal-containers.example/manifest/workload/server \
    -parentID spiffe://coastal-containers.example/agent/spire-agent \
    -selector k8s:ns:default \
    -selector k8s:sa:server
```

Now do the same for the manifest client:

```shell
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://coastal-containers.example/manifest/workload/client \
    -parentID spiffe://coastal-containers.example/agent/spire-agent \
    -selector k8s:ns:default \
    -selector k8s:sa:client
```

### Step 5: Crafting the `helper.conf` files

Now you are ready to write the configuration files for the [manifest-server](manifest-server) and 
[manifest-client](manifest-client) workloads we will deploy in this lab. First create the files for both of the
workloads. 

```shell
touch {manifest-server,manifest-client}/helper.conf
```

First edit the [server helper.config](manifest-server/helper.conf) file to contain the following content.

```text
# server helper.conf

# Socket address of the SPIRE Agent. Update this according to your SPIRE setup.
agent_address = "/spire-agent-socket/agent.sock"

# The command that represents the workload (in this case, Python scripts).
cmd = "python3"
cmd_args = "main.py"

# Signal to send to the Python script to trigger a certificate reload.
renew_signal = "SIGHUP"

# Directory where certificates will be written. This should match the path in your Python scripts.
cert_dir = "/var/run/secrets/svids"

# Names of the files where the SVID, private key, and bundle will be stored.
svid_file_name = "server_cert.pem"
svid_key_file_name = "server_key.pem"
svid_bundle_file_name = "svid_bundle.pem"
```

As you can see, we are configuring the SPIFFE Helper to know:

- Where to access the Workload API
  - `agent_address = "/spire-agent-socket/agent.sock"`
- How to launch the workload as a subprocess
  - `cmd = "python3"`
  - `cmd_args = "main.py"`
- How to signal the workload to reload it's configuration when the SVIDs are rotated
  - `renew_signal = "SIGHUP"`
- Where to write the X.509 materials and which filenames to use
  - `cert_dir = "/var/run/secrets/svids"`
  - `svid_file_name = "server_cert.pem"`
  - `svid_key_file_name = "server_key.pem"`
  - `svid_bundle_file_name = "svid_bundle.pem"`

Next edit the [client helper.config](manifest-client/helper.conf) file to contain the following content, tailored to the
manifest client.

```text
# client helper.conf

# Socket address of the SPIRE Agent. Update this according to your SPIRE setup.
agent_address = "/spire-agent-socket/agent.sock"

# The command that represents the workload (in this case, Python scripts).
cmd = "python3"
cmd_args = "main.py"

# Signal to send to the Python script to trigger a certificate reload.
renew_signal = "SIGHUP"

# Directory where certificates will be written. This should match the path in your Python scripts.
cert_dir = "/var/run/secrets/svids"

# Names of the files where the SVID, private key, and bundle will be stored.
svid_file_name = "client_cert.pem"
svid_key_file_name = "client_key.pem"
svid_bundle_file_name = "svid_bundle.pem"
```

Finally, create ConfigMaps for the configuration in the Kind cluster:

```shell
kubectl create configmap server-helper-config --from-file=manifest-server/helper.conf
kubectl create configmap client-helper-config --from-file=manifest-client/helper.conf
```

### Step 6: Deploying the Workloads

The [manifest server](manifest-server/main.py) is a simply Python web application that returns a ship's manifest as a 
json payload over HTTPS. It is not SPIFFE aware, so simply loads its CA Certificate, Server Certificate and Key from 
the filesystem.

Without changing the Python code, we can craft a [Dockerfile](manifest-server/Dockerfile) to containerise the Python
application with the SPIFFE Helper so that the SPIFFE Helper can manage the retrieval and rotation of the X.509 SVIDs
and the launching and signalling of the application.

We [deploy](manifest-server/app.yaml) the manifest server to the Kind cluster mounting the SPIFFE Helper configuration
from the ConfigMap we created earlier and the Workload API.

The [manifest-client](manifest-client/main.py) connects to the manifest server using its X.509 SVID, again managed by
the SPIFFE Helper, retrieves the manifest, and outputs it to stdout. Its [Dockerfile](manifest-client/Dockerfile) is
almost identity to the manifest server, and it also mounts the SPIFFE Helper configuration from the ConfigMap and the 
Workload API when [deployed](manifest-client/app.yaml).

Both the server and the client are configured for mutual TLS:

```python
context.verify_mode = ssl.CERT_REQUIRED
context.load_verify_locations(cafile=TRUST_BUNDLE)
```

First we need to build the manifest server and client images and load them into the Kind cluster:

```shell
make workload-images deploy-workloads
```

To ensure the server is running smoothly and received it's SVID from the `spire-agent`, check its logs:

```shell
kubectl logs -l app=server
```

To ensure the client is running smoothly, received it's SVID from the `spire-agent`, and retrieved the shipping manifest
from the server workload, check its logs:

```shell
kubectl logs -l app=client
```

### Step 7: Observing the Results

Upon inspecting the client logs, you'll observe the fetched shipping manifest, which confirms the secure communication
between the server and client. Dive deeper:

- **SVID Rotation**: With spiffe-helper, SVIDs are rotated seamlessly. Monitor the logs to see how SVIDs are renewed
over time, ensuring continuous secure communication.

- **Python Integration**: Recognize how Python applications, which aren't inherently SPIFFE-aware, can now securely
interact within the Kubernetes cluster, thanks to the spiffe-helper.

- **Secure Communication**: The client's ability to fetch the shipping manifest from the server underscores the secure
communication established between them, facilitated by SPIFFE/SPIRE.

### Step 8: Cleanup

To tear down the entire Kind cluster, run:

```shell
make cluster-down
```

## Conclusion

⚓ Congrats cap'n! You have successfully provided X.509 SVIDs for non SPIFFE-aware workloads by enlisting the help of the
trusty [spiffe-helper](https://github.com/spiffe/spiffe-helper). You should now have an understanding
of how you can integrate SPIFFE into existing applications without having to refactor them to be SPIFFE-aware. With the
power of SPIFFE/SPIRE, combined with tools like [spiffe-helper](https://github.com/spiffe/spiffe-helper), organizations
like Coastal Containers Ltd. can secure their existing systems in modern, dynamic environments.
