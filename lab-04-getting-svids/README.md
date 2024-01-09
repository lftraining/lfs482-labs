# Lab 4: Getting SVIDS with SPIFFE-Helper

## Prerequisites

- A AMD, Intel, or Arm 64-bit Linux environment.
- Familiarity with `docker` and `kubectl` command-line interface (CLI) commands is helpful.

## Introduction

Coastal Containers Ltd. is a leading container shipping company with a legacy manifest file server system. As part of their digital transformation journey, they recognized the need to improve the security of their system. The challenge? Their system is written in Python, a language that doesn't natively support go-spiffe.

To address this, they are adopting a Zero Trust security model, leveraging SPIFFE/SPIRE. This approach allows them to identify and authenticate workloads in heterogeneous environments securely. This lab focuses on using the [spiffe-helper](https://github.com/spiffe/spiffe-helper) to enable Python apps to become SPIFFE-aware, even if they don't natively support go-spiffe.

### Learning Objectives

- Understand the role and utility of [spiffe-helper](https://github.com/spiffe/spiffe-helper).
- Learn how to retrieve and rotate SVIDs for applications that don't natively support go-spiffe.
- Explore the mechanisms of signal-based SVID rotation with spiffe-helper.

## Preparing your Environment

Before you cast off, prepare your ships to sail by setting your working directory in [lab-04-getting-svids](../lab-04-getting-svids/) as an environment variable:

```bash
export LAB_DIR=$(pwd)
```

This will make issuing commands easier in the following steps of this exercise, and will reduce the possibility of reference errors.

## Step-by-Step Instructions

### Step 1: Provision Infrastructure

To set sail, spin up your demo Kubernetes cluster using [`Kind`](https://kind.sigs.k8s.io/) by issusing the following `make` command while in the root [lab directory](/):

```bash
make cluster-up
```

If your cluster is already running, you can skip this step and continue on with the lab.

### Step 2: Deploying the SPIRE Setup

Initialize the SPIRE setup on your Kubernetes cluster:

```shell
make deploy-spire
```

This command will deploy SPIRE based on the three YAML manifests for the [`spire-server`](./config/spire-server.yaml), [`spire-agent`](./config/spire-agent.yaml), and [`spiffe-csi-driver`](./config/spiffe-csi-driver.yaml) located within the [config directory](./config/). We've already covered what the SPIRE Server and Agent do within the previous lab exercises, but we haven't yet discussed what the `spiffe-csi-driver` is. 

The [SPIFFE CSI Driver](https://github.com/spiffe/spiffe-csi), is a [Container Storage Interface (CSI)](https://kubernetes-csi.github.io/docs/) driver for Kubernetes that facilitates injection of the [SPIFFE Workload API](https://github.com/spiffe/spiffe/blob/main/standards/SPIFFE_Workload_API.md). This is done by mounting a directory containing the Workload API Unix domain socket as an [ephemeral inline volume](https://kubernetes.io/blog/2020/01/21/csi-ephemeral-inline-volumes/) into your workload pods. Typically, the CSI driver will share the directory hosting the Workload API with a read-only bind mount into the container at the requested target path.

This begs the question of what a [ephemeral inline volume](https://kubernetes.io/blog/2020/01/21/csi-ephemeral-inline-volumes/) is, and why is it useful. Traditionally, voumes provided by external storage drivers in K8s are persistent by nature. This means that they operate on a completely independent lifecycle to the pods (e.g., persisting even if a pod is killed). In use cases where data volumes are required to be tied directly to the lifecycle of a pod, they need to be created and deleted with the pod (aka ephemeral). To achieve this, these volumes are defined as part of the pod spec itself (aka inline). 
Circling back to the `spiffe-csi-driver`, this means that when the pod is destroyed, the driver is invoked and removes the bind mount.

### Step 3: Create Initial Node Registration Entry

Using your SPIRE tools, create the initial node registration entry. This step ensures your SPIRE server recognizes nodes in your cluster.

```shell
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://coastal-containers.example/agent/spire-agent \
    -selector k8s_psat:cluster:kind-kind \
    -selector k8s_psat:agent_ns:spire \
    -selector k8s_psat:agent_sa:spire-agent \
    -node
```

If the command is successful, you should see an outputted SVID in your terminal. 

### Step 4: Crafting the `helper.conf` files

Now you are ready to write `helper.conf` files for your SPIRE Server and Agent within the the [client](./workload/client/) and [server](./workload/server/) app directories. To do so, first navigate to the [server](./workload/server/) app directory:

```shell
cd $LAB_DIR/workload/server
```

Once here, investigate the [Dockerfile](./workload/server/Dockerfile) and [main.py](./workload/server/main.py) to better understand how we've setup the manifest server using some simple Python code and a docker image. Within the same directory, create the Server `helper.conf` file by running:

```shell
vi helper.conf
```

Feel free to use an alternative method of file creation if you don't have vi installed, are running on a different operating system, or prefer another method (e.g. using an IDE). Once the file is created, reference the below example of a server `helper.conf` for Coastal Containers Ltd to guide your configuration:

```conf
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

After you've added the appropiate configuration, and the server `helper.conf` is in place, navigate to the [client](./workload/client/) app directoy by running:

```shell
cd $LAB_DIR/workload/client
```

Once here, investigate the [Dockerfile](./workload/client/Dockerfile) and [main.py](./workload/client/main.py) to better understand how we've setup the manifest client using some simple Python code and a docker image. Within the same directory, create the client `helper.conf` file by running:

```shell
vi helper.conf
```

Feel free to use an alternative method of file creation if you don't have vi installed, are running on a different operating system, or prefer another method (e.g. using an IDE). Once the file is created, reference the below example of a client `helper.conf` for Coastal Containers Ltd to guide your configuration:

```conf
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

Alternatively, another sample `helper.conf` file can be found [here](https://github.com/spiffe/spiffe-helper/blob/main/helper.conf) within the [SPIFFE Helper GitHub repository](https://github.com/spiffe/spiffe-helper/). With the configuration files in place, you can now move on to build and load the two workload images using `go-spiffe` as a wrapper.

### Step 5: Building and Loading Workload Images

First, change directory to view the server workload files:

```shell
cd ${LAB_DIR}/workload/server
```

Build the server workload image:

```shell
docker build -t manifest-server .
```

Load the server workload image into the kind cluster:

```shell
kind load docker-image manifest-server
```

Next, change directory to view the client workload files:

```shell
cd ${LAB_DIR}/workload/client
```

Build the client workload image:

```shell
docker build -t manifest-client .
```

Load the server client image into the kind cluster:

```shell
kind load docker-image manifest-client
```

### Step 6: Create [ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/)

To enable the workload pods to consume the crafted server and client `helper.conf` file, you will need to create two Kubernetes ConfigMaps to hold the configuration data for `spiffe-helper`.

First, cd into the [`workload/server`](./workload/server/) directory and create the `server-helper-config` ConfigMap:

```shell
cd ${LAB_DIR}/workload/server
kubectl create configmap server-helper-config --from-file=helper.conf
```

Next, cd into the [`workload/client`](./workload/client/) directory and create the `client-helper-config` ConfigMap:

```shell
cd ${LAB_DIR}/workload/client
kubectl create configmap client-helper-config --from-file=helper.conf
```

Optionally, to ensure the `spiffe-helper` config was properly set, inspect the yaml of the newly created config maps by running:

```shell
kubectl get configmap/client-helper-config configmap/server-helper-config -o yaml
```

The outputted YAML should reflect the configurations you created in your local `helper.conf` files.

### Step 7: Create Registration Entries for Server and Client Workloads

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

### Step 8: Deploying the Server

Change to the [`config`](./config/) directory, and then apply the server deployment manifest:

```shell
cd ${LAB_DIR}/config
kubectl apply -f deploy-server.yaml
```

To ensure the server is running smoothly and received it's SVID from the `spire-agent`, check its logs:

```shell
kubectl logs -l app=server
```

### Step 9: Deploying the Client

While still in the [`config`](./config/) directory, apply the client deployment manifest by running:

```shell
kubectl apply -f deploy-client.yaml
```

To ensure the client is running smoothly, received it's SVID from the `spire-agent`, and retrieved the shipping manifest from the server workload, check its logs:

```shell
kubectl logs -l app=client
```

### Step 10: Observing the Results

Upon inspecting the client logs, you'll observe the fetched shipping manifest, which confirms the secure communication between the server and client. Dive deeper:

- **SVID Rotation**: With spiffe-helper, SVIDs are rotated seamlessly. Monitor the logs to see how SVIDs are renewed over time, ensuring continuous secure communication.
  
- **Python Integration**: Recognize how Python applications, which aren't inherently SPIFFE-aware, can now securely interact within the Kubernetes cluster, thanks to the spiffe-helper.
  
- **Secure Communication**: The client's ability to fetch the shipping manifest from the server underscores the secure communication established between them, facilitated by SPIFFE/SPIRE.

### Step 11: Cleanup

As the following lab exercises will use the same cluster, tear down the SPIRE setup and workload deployments from this lab by running:

```shell
cd $LAB_DIR
make tear-down
```

To tear down the entire Kind cluster, run:

```shell
cd $LAB_DIR
make cluster-down
```

## Conclusion

⚓ Congrats cap'n! You have successfully implemented a SPIFFE-aware workload by enlisting the help of the trusty [`spiffe-helper`](https://github.com/spiffe/spiffe-helper). By the end of this lab, you should have a better understanding of how [`spiffe-helper`](https://github.com/spiffe/spiffe-helper) can bridge the gap for applications written in languages that don't natively support [`go-spiffe`](https://github.com/spiffe/go-spiffe). With the power of SPIFFE/SPIRE, combined with tools like [`spiffe-helper`](https://github.com/spiffe/spiffe-helper), organizations like Coastal Containers Ltd. to secure their legacy systems in modern, dynamic environments.
