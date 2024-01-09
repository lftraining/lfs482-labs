# Lab 3: Setup SPIRE on Kubernetes with [Kind]([`Kind`](https://kind.sigs.k8s.io/))

## Prerequisites

- A AMD, Intel, or Arm 64-bit Linux environment.
- Familiarity with Kubernetes manifests and `kubectl` commands is helpful.

## Introduction

⛵ Welcome aboard, tech savvy sailor! ⚓ As Coastal Containers sails towards a cyber-safe future, you're help is needed to ensure that their ship-to-shore communications are being transmitted securely. In this lab, you'll don the hat of a ship's engineer working to get SPIRE deployed and running in a Kubernetes cluster. Kubernetes will help to provide a containerized and easily scalable platform from which to distribute your coastal cargo. 📦 By these benefits, Coastal Cargo hopes to leverage Kubernetes for it's interopability with their enterprise shipping systems and resiliency towards hardware failures. 

### Preparing Your Environment

Before you cast off, prepare your ships to sail by setting your working directory in [lab-03-kind-install](../lab-03-kind-install/) as an environment variable:

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

### Step 2: View Sample Configurations

To view the sample SPIRE Server and Agent configurations, use the make command:

```bash
make view-sample-server-config
make view-sample-agent-config
```

After executing these commands, compare the output against the example provided within [lab-02-binary-install](../lab-02-binary-install), and the [SPIRE Server](https://spiffe.io/docs/latest/deploying/spire_server/) / [SPIRE Agent](https://spiffe.io/docs/latest/deploying/spire_agent/) configuration references.

You can also inspect the sample SPIRE Server and Agent deployment manifests using the commands:

```shell
make view-sample-agent-deployment
make view-sample-server-deployment
```

Optionally, these YAML manifests can be inspected manually using your favorite text editor within the [sample/config](./sample/config/) directory.

### Step 3: Create Namespace

Create the `atlantic-coast` namespace where we will run your SPIRE setup:

```bash
kubectl create namespace atlantic-coast
```

### Step 4: Apply RBAC for Server

Write and apply the RBAC roles for the server based on [`server-roles.yaml`](./sample/config/server-roles.yaml):

```bash
kubectl apply -f ./sample/config/server-roles.yaml
```

### Step 5: Apply Server ConfigMap

Write and apply the server ConfigMap based on the sample [`server-config.yaml`](./sample/config/server-config.yaml):

```bash
kubectl apply -f ./sample/config/server-config.yaml
```

Note: Pay attention to key configuration items like `trust_domain`, `server_address`, etc.

### Step 6: Deploy Server

Write and apply the server deployment based on the sample [`server-deploy.yaml`](./sample/config/server-deploy.yaml):

```bash
kubectl apply -f ./sample/config/server-deploy.yaml
```

Wait until the server is ready:

```bash
make wait-for-server
```

### Step 7: Apply RBAC for Agent

Write and apply the RBAC roles for the agent based on [`agent-roles.yaml`](./sample/config/agent-roles.yaml):

```bash
kubectl apply -f ./sample/config/agent-roles.yaml
```

### Step 8: Apply Agent ConfigMap

Write and apply the agent ConfigMap based on the sample [`agent-config.yaml`](./sample/config/agent-config.yaml):

```bash
kubectl apply -f ./sample/config/agent-config.yaml
```

Note: Pay attention to key configuration items like `trust_domain`, `server_address`, etc.

### Step 9: Deploy Agent

Write and apply the agent deployment based on the sample [`agent-deploy.yaml`](./sample/config/agent-deploy.yaml):

```bash
kubectl apply -f ./sample/config/agent-deploy.yaml
```

Wait until the agent is ready:

```bash
make wait-for-agent
```

### Step 10: View Logs

You can view logs for the server and agent using the make commands:

```bash
make view-server-logs
make view-agent-logs
```

Or directly via kubectl using:

```shell
kubectl logs -f coastal-server-0 -n atlantic-coast
kubectl logs -f coastal-agent-#### -n atlantic-coast
```

These commands will follow the logs of your `spire-server` and `spire-agent` using the `-f` flag. If you want to exit this output view, issue a `ctrl + c`.

Inspecting logs provides insights into the operations of the server and agent:

- **Server Logs**: By viewing the server logs, you can observe the initialization process, registration of entities, and the issuance of SPIFFE Verifiable Identity Documents (SVIDs).
  
- **Agent Logs**: The agent logs shed light on the attestation process, where the agent proves its identity to the server, and the subsequent retrieval & renewal of SVIDs for workloads.

### Step 11: Create Node Registration Entry

Create a node registration entry:

```bash
make node-registration-entry
```

Upon creating a node registration entry, it's crucial to understand its significance:

- **Node Attestation**: This step is pivotal for the SPIRE server to recognize and trust the nodes in your Kubernetes cluster. The node registration entry contains selectors that help the SPIRE server identify and authenticate nodes.
  
- **SVID Issuance**: Once a node is attested, it's granted an SVID. This SVID is essential for secure communications within the SPIRE infrastructure.

Take note of the output after registration, as it provides intricate details about the newly minted SVID.

### Step 12: Cleanup

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

Bravo! You've fortified Coastal Containers' communication channels, making it a daunting task for adversaries like Captain Hashjack and his cyber-pirate crew to breach our defenses. By leveraging the insights from Lab 2 and applying them to Kubernetes, you've laid the foundation for a robust SPIRE setup. This expertise is invaluable as we navigate the tumultuous waters of Zero Trust and integrate it seamlessly into Coastal Containers' vast infrastructure.

For the adventurous souls yearning for deeper waters, myriad advanced configurations and deployment strategies await your exploration. Dive into comprehensive [SPIRE deployment examples and configurations](https://github.com/spiffe/spire-examples) to quench your thirst for knowledge. May your journey be marked by calm seas and favorable winds!
