# Lab 3: Setup SPIRE on Kubernetes with Kind

## Prerequisites

- An AMD, Intel, or Arm 64-bit Linux environment.
- Familiarity with Kubernetes manifests and `kubectl` commands is helpful.

## Introduction

⛵ Welcome aboard, tech savvy sailor! ⚓ As Coastal Containers sails towards a cyber-safe future, your help is needed
to ensure that their ship-to-shore communications are being transmitted securely. In this lab, you'll don the hat of a
ship's engineer working to get SPIRE deployed and running in a Kubernetes cluster. Kubernetes will help to provide a
containerized and easily scalable platform from which to distribute your coastal cargo. 📦 By these benefits, Coastal
Cargo hopes to leverage Kubernetes for its interoperability with their enterprise shipping systems and resiliency
towards hardware failures.

## Step-by-Step Instructions

### Step 1: Provision Infrastructure

To set sail, spin up your demo Kubernetes cluster using [Kind](https://kind.sigs.k8s.io/) by issuing the following
command:

```bash
make cluster-up
```

This will also load the required SPIRE images into the cluster that were previously pulled in 
[lab-00-setup](../lab-00-setup).

### Step 2: View  Configuration

To view the SPIRE Server and Agent configurations in the [spire-server](spire-server) and [spire-agent](spire-agent) 
directories. Compare these to the configurations

We use a ConfigMap to store the [server configuration](spire-server/config.yaml) and the 
[agent configuration](spire-agent/config.yaml). Compare these to the configuration from the previous lab and the
[SPIRE Server](https://spiffe.io/docs/latest/deploying/spire_server/) and
[SPIRE Agent](https://spiffe.io/docs/latest/deploying/spire_agent/) configuration references.

The key updates to the server configuration, compares to the previous lab, are:

- [k8s_psat NodeAttestor](https://github.com/spiffe/spire/blob/v1.9.0/doc/plugin_server_nodeattestor_k8s_psat.md)
- [k8sbundle Notifier](https://github.com/spiffe/spire/blob/v1.9.0/doc/plugin_server_notifier_k8sbundle.md)

These are used for agent bootstrapping in a Kubernetes environment. The server updates the ConfigMap with the trust 
bundle used to bootstrap the agents. Node attestation is configured to allow the `spire-agent` ServiceAccount in the 
`spire` Namespace. The server verifies the identity in the provided PSAT using the Kubernetes TokenReview API.

The SPIRE server also needs to be able to get information about Nodes and Pods on the Kubernetes cluster.

The required [RBAC permissions](spire-server/rbac.yaml) are granted to the `spire-server` ServiceAccount.

Finally, the [server](spire-server/server.yaml) is deployed as a StatefulSet with a Service.

The key updates to the agent configuration, compares to the previous lab, are:

- [k8s_psat NodeAttestor](https://github.com/spiffe/spire/blob/v1.9.0/doc/plugin_agent_nodeattestor_k8s_psat.md)
- [k8s WorkloadAttestor](https://github.com/spiffe/spire/blob/v1.9.0/doc/plugin_agent_workloadattestor_k8s.mdq)

The server and agent need paired Node Attestors, in this case `k8s_psat`. The agent is also configured to the `k8s`
Workload Attestor, and requires [these RBAC permissions](spire-agent/rbac.yaml).

The SPIRE agent is deployed as a [DaemonSet](spire-agent/agent.yaml) so that an agent runs on every Node in the cluster.

### Step 3: Create Namespace

Create the `spire` namespace where we will run your SPIRE setup:

```bash
kubectl create namespace spire
```

### Step 4: Deploy the SPIRE Server

```bash
kubectl apply -f spire-server
```

### Step 5: Deploy the SPIRE Agent

```bash
kubectl apply -f spire-agent
```

Wait until the agent is ready:

```bash
make spire-wait-for-agent
```

### Step 10: View Logs

You can view logs for the server and agent using the make commands:

```bash
make view-server-logs
```

```shell
make view-agent-logs
```

Or directly via kubectl using:

```shell
kubectl logs -f spire-server-0 -n spire
```

```shell
kubectl logs -f -l=app=spire-agent -n spire
```

These commands will follow the logs of your `spire-server` and `spire-agent` using the `-f` flag. If you want to exit
this output view, issue a `ctrl + c`.

Inspecting logs provides insights into the operations of the server and agent:

- **Server Logs**: By viewing the server logs, you can observe the initialization process, registration of entities,
and the issuance of SPIFFE Verifiable Identity Documents (SVIDs).

- **Agent Logs**: The agent logs shed light on the attestation process, where the agent proves its identity to the
server, and the subsequent retrieval & renewal of SVIDs for workloads.

### Step 11: Create Node Registration Entry

Create a node registration entry:

```bash
make node-registration-entry
```

Upon creating a node registration entry, it's crucial to understand its significance:

- **Node Attestation**: This step is pivotal for the SPIRE server to recognize and trust the nodes in your Kubernetes
cluster. The node registration entry contains selectors that help the SPIRE server identify and authenticate nodes.

- **SVID Issuance**: Once a node is attested, it's granted an SVID. This SVID is essential for secure communications
within the SPIRE infrastructure.

Take note of the output after registration, as it provides intricate details about the newly minted SVID.

### Step 12: Cleanup

To tear down the entire Kind cluster, run:

```shell
make cluster-down
```

## Conclusion

Bravo! You've fortified Coastal Containers' communication channels, making it a daunting task for adversaries like
Captain Hashjack and his cyber-pirate crew to breach our defenses. By leveraging the insights from Lab 2 and applying
them to Kubernetes, you've laid the foundation for a robust SPIRE setup. This expertise is invaluable as we navigate
the tumultuous waters of Zero Trust and integrate it seamlessly into Coastal Containers' vast infrastructure.

For the adventurous souls yearning for deeper waters, myriad advanced configurations and deployment strategies await
your exploration. Dive into comprehensive
[SPIRE deployment examples and configurations](https://github.com/spiffe/spire-examples) to quench your thirst for
knowledge. May your journey be marked by calm seas and favorable winds!
