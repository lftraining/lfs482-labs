# Lab 13: Cilium with SPIRE

## Prerequisites

- An AMD, Intel, or Arm 64-bit Linux environment.
- Familiarity with Kubernetes manifests and `kubectl` commands is helpful.

## Introduction

Ahoy, matey! Ye have done a fine job of setting up SPIRE on yer Kubernetes cluster, but ye still need to make sure that
yer network policies are shipshape and seaworthy. That’s where [Cilium](https://cilium.io/) comes in handy. Cilium is a
Kubernetes network plugin that uses a Linux kernel technology called [eBPF](https://ebpf.io/) to dynamically insert
security, visibility and control logic within the Linux kernel. Cilium provides distributed load-balancing for
pod-to-pod traffic and identity-based implementation of the NetworkPolicy resource. By integrating SPIRE with Cilium,
ye can use the SPIFFE identities issued by SPIRE to define and enforce granular network policies that are
platform-agnostic and resilient to node failures. In this section, ye will learn how to install Cilium on yer Kubernetes
cluster, configure it to use SPIRE as the identity provider, and test the network policy enforcement using an example
application. This will help Coastal Containers to secure their ship-to-shore communications and prevent any unwanted
intruders from boarding their vessels.

### Preparing Your Environment

Before you cast off, prepare your ship to sail by setting up your working environment. If you haven't yet done so, make
sure you've cloned the lab repository to your local system. After that, you'll be working from the
[lab-13-cilium-spire](../lab-13-cilium-spire) directory.

```bash
export LAB_DIR=$(pwd)
export PATH=$PATH:$(pwd)/../bin
```

## Step-by-Step Instructions

### Step 1: Install cilium on kind cluster

```bash
make cluster-up
```

You'll need to generate a IPSec key for Cilium to use for encrypting traffic.

```shell
kubectl create -n kube-system secret generic cilium-ipsec-keys \
    --from-literal=keys="3 rfc4106(gcm(aes)) \
    $(echo $(dd if=/dev/urandom count=20 bs=1 2> /dev/null | xxd -p -c 64)) 128"
```

Then, install Cilium via Helm:

```shell
helm repo add cilium https://helm.cilium.io/

helm install cilium cilium/cilium --version 1.14.2 \
   --namespace kube-system \
   --set image.pullPolicy=IfNotPresent \
   --set ipam.mode=kubernetes \
   --set operator.replicas=1 \
   --set authentication.mutual.spire.enabled=true \
   --set authentication.mutual.spire.install.enabled=false \
   --set authentication.mutual.spire.install.namespace=spire \
   --set encryption.enabled=true \
   --set hubble.relay.enabled=true \
   --set hubble.ui.enabled=true \
   --set authentication.mutual.spire.trustDomain=coastal-containers.example
```

To validate that Cilium has been properly installed, you can run

```shell
kubectl exec -n kube-system ds/cilium -c cilium-agent -- cilium encrypt status
```

Should return something like:

```text
Encryption: IPsec
Keys in use: 1
Max Seq. Number: 0x115/0xffffffff
Errors: 0
```

To validate that Cilium is running properly, you can run

```shell
# wait for cilium to be running everywhere
kubectl rollout status -n kube-system daemonset cilium

cilium status --wait
```

And you should see a response like:

```text
    /¯¯\
 /¯¯\__/¯¯\    Cilium:             OK
 \__/¯¯\__/    Operator:           OK
 /¯¯\__/¯¯\    Envoy DaemonSet:    disabled (using embedded mode)
 \__/¯¯\__/    Hubble Relay:       OK
    \__/       ClusterMesh:        disabled

Deployment             hubble-relay       Desired: 1, Ready: 1/1, Available: 1/1
Deployment             cilium-operator    Desired: 1, Ready: 1/1, Available: 1/1
DaemonSet              cilium             Desired: 1, Ready: 1/1, Available: 1/1
Deployment             hubble-ui          Desired: 1, Ready: 1/1, Available: 1/1
Containers:            hubble-relay       Running: 1
                       cilium-operator    Running: 1
                       hubble-ui          Running: 1
                       cilium             Running: 1
Cluster Pods:          11/11 managed by Cilium
Helm chart version:    1.14.2
Image versions         hubble-relay       quay.io/cilium/hubble-relay:v1.14.2@sha256:a89030b31f333e8fb1c10d2473250399a1a537c27d022cd8becc1a65d1bef1d6: 1
                       cilium-operator    quay.io/cilium/operator-generic:v1.14.2@sha256:52f70250dea22e506959439a7c4ea31b10fe8375db62f5c27ab746e3a2af866d: 1
                       hubble-ui          quay.io/cilium/hubble-ui:v0.12.0@sha256:1c876cfa1d5e35bc91e1025c9314f922041592a88b03313c22c1f97a5d2ba88f: 1
                       hubble-ui          quay.io/cilium/hubble-ui-backend:v0.12.0@sha256:8a79a1aad4fc9c2aa2b3e4379af0af872a89fcec9d99e117188190671c66fc2e: 1
                       cilium             quay.io/cilium/cilium:v1.14.2@sha256:6263f3a3d5d63b267b538298dbeb5ae87da3efacf09a2c620446c873ba807d35: 1
```

### Step 2: Install SPIRE on kind cluster

Take note of the configuration of both the SPIRE agent and Spire server, where we make sure cilium is in the list of
`authorized_delegates` and can talk to the Admin API.

```shell
make deploy-spire
```

### Step 3: Generate SPIFFE entries for Cilium and the workload

```shell
# Cilium Agent
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://coastal-containers.example/cilium-agent \
    -parentID spiffe://coastal-containers.example/agent/spire-agent \
    -selector k8s:ns:kube-system \
    -selector k8s:sa:cilium

# Cilium Operator
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://coastal-containers.example/cilium-operator \
    -parentID spiffe://coastal-containers.example/agent/spire-agent \
    -selector k8s:ns:kube-system \
    -selector k8s:sa:cilium-operator

# Client Workload
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://coastal-containers.example/workload/client \
    -parentID spiffe://coastal-containers.example/agent/spire-agent \
    -selector k8s:ns:default \
    -selector k8s:sa:client \
    -ttl 60

# Server Workload
kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://coastal-containers.example/workload/server \
    -parentID spiffe://coastal-containers.example/agent/spire-agent \
    -selector k8s:ns:default \
    -selector k8s:sa:server \
    -ttl 60
```

Now that SPIRE is all in, you'll need to restart cilium so that it can register with it correctly:

```shell
kubectl -n kube-system rollout restart deployment/cilium-operator ds/cilium
kubectl rollout status -n kube-system daemonset cilium

```

Run the following command to validate that your cluster has proper network connectivity:

```shell
cilium connectivity test
```

You should hopefully see output like:

```text
ℹ️  Monitor aggregation detected, will skip some flow validation steps
✨ [kind-kind] Creating namespace cilium-test for connectivity check...
✨ [kind-kind] Deploying echo-same-node service...
✨ [kind-kind] Deploying DNS test server configmap...

(...)

✅ All 45 tests (311 actions) successful, 10 tests skipped, 1 scenarios skipped.
```

If you get any errors look into
[the Cilium troubleshooting guide](https://docs.cilium.io/en/stable/operations/troubleshooting/).

### Step 4: Deploy a workload

```shell
kubectl apply --wait -f $LAB_DIR/manifests/1-workload.yaml
```

### Step 5: Prove networking works

```shell
make test-workload-networking
```

If this works you should see:

```text
✅ Workload Ping was successful - ICMP between client>server works!
✅ Workload curl was successful - HTTP between client>server works!
✅ Undesirable Workload Ping was successful - ICMP between server>client works!
✅ Undesirable Workload curl was successful - HTTP between server>client works!
✅ External Ping was successful - ICMP between client>1.1.1.1 works!
✅ External Curl was successful - HTTP between client>1.1.1.1 works!
...
```

### Step 6: Add deny policy

```shell
kubectl apply -f $LAB_DIR/manifests/2-deny.yaml
```

### Step 7: Prove networking no longer works

```shell
make test-workload-networking
```

If this works you should see:

```text
❌ Workload Ping failed - ICMP between client>server does not work!
❌ Workload curl failed - HTTP between client>server does not work!
❌ Undesirable Workload Ping failed - ICMP between server>client does not work!
❌ Undesirable Workload curl failed - HTTP between server>client does not work!
❌ External Ping failed - ICMP between client>1.1.1.1 does not work!
❌ External Curl failed - HTTP between client>1.1.1.1 does not work!
...
```

### Step 8: Apply a cilium policy to allow the requests

```shell
kubectl apply -f $LAB_DIR/manifests/3-allow.yaml
```

### Step 9: Prove networking works again

Note: you may need to run this a couple of times to get the policy to take effect

```shell
make test-workload-networking
```

If this works you should see:

```text
❌ Workload Ping failed - ICMP between client>server does not work!
✅ Workload curl was successful - HTTP between client>server works!
❌ Undesirable Workload Ping failed - ICMP between server>client does not work!
❌ Undesirable Workload curl failed - HTTP between server>client does not work!
❌ External Ping failed - ICMP between client>1.1.1.1 does not work!
❌ External Curl failed - HTTP between client>1.1.1.1 does not work!
Oct  3 15:31:26.464: default/client (ID:2118) <> default/server (ID:33458) policy-verdict:none EGRESS DENIED (ICMPv4 EchoRequest)
Oct  3 15:31:27.472: default/client (ID:2118) <> default/server (ID:33458) policy-verdict:none EGRESS DENIED (ICMPv4 EchoRequest)
Oct  3 15:31:27.684: default/client:59112 (ID:2118) -> default/server:80 (ID:33458) policy-verdict:L3-L4 EGRESS ALLOWED (TCP Flags: SYN; Auth: SPIRE)
Oct  3 15:31:27.788: default/server (ID:33458) <> default/client (ID:2118) policy-verdict:none EGRESS DENIED (ICMPv4 EchoRequest)
Oct  3 15:31:28.812: default/server (ID:33458) <> default/client (ID:2118) policy-verdict:none EGRESS DENIED (ICMPv4 EchoRequest)
Oct  3 15:31:29.028: default/server:41194 (ID:33458) <> default/client:80 (ID:2118) policy-verdict:none EGRESS DENIED (TCP Flags: SYN)
Oct  3 15:31:29.253: default/client (ID:2118) <> 0.0.0.1 (world) policy-verdict:none EGRESS DENIED (ICMPv4 EchoRequest)
Oct  3 15:31:30.284: default/client (ID:2118) <> 0.0.0.1 (world) policy-verdict:none EGRESS DENIED (ICMPv4 EchoRequest)
Oct  3 15:31:30.518: default/client:55730 (ID:2118) <> 1.1.1.1:80 (world) policy-verdict:none EGRESS DENIED (TCP Flags: SYN)
Oct  3 15:34:18.649: default/client (ID:2118) <> default/server (ID:33458) policy-verdict:none EGRESS DENIED (ICMPv4 EchoRequest)
Oct  3 15:34:19.693: default/client (ID:2118) <> default/server (ID:33458) policy-verdict:none EGRESS DENIED (ICMPv4 EchoRequest)
Oct  3 15:34:19.901: default/client:58472 (ID:2118) -> default/server:80 (ID:33458) policy-verdict:L3-L4 EGRESS ALLOWED (TCP Flags: SYN; Auth: SPIRE)
Oct  3 15:34:19.995: default/server (ID:33458) <> default/client (ID:2118) policy-verdict:none EGRESS DENIED (ICMPv4 EchoRequest)
Oct  3 15:34:21.036: default/server (ID:33458) <> default/client (ID:2118) policy-verdict:none EGRESS DENIED (ICMPv4 EchoRequest)
Oct  3 15:34:21.247: default/server:33336 (ID:33458) <> default/client:80 (ID:2118) policy-verdict:none EGRESS DENIED (TCP Flags: SYN)
Oct  3 15:34:21.449: default/client (ID:2118) <> 0.0.0.1 (world) policy-verdict:none EGRESS DENIED (ICMPv4 EchoRequest)
Oct  3 15:34:22.513: default/client (ID:2118) <> 0.0.0.1 (world) policy-verdict:none EGRESS DENIED (ICMPv4 EchoRequest)
Oct  3 15:34:22.711: default/client:52426 (ID:2118) <> 1.1.1.1:80 (world) policy-verdict:none EGRESS DENIED (TCP Flags: SYN)
```

If its all working you should see `Auth: SPIRE` littered through the allowed connections.

### Step 10: Cleanup

Now that you've proved everything works, its time to scrub the decks and delete your cluster:

```shell
cd $LAB_DIR && make cluster-down
```

## Conclusion

Congratulations, sailor! Ye have completed the lab and learned how to integrate SPIRE with Cilium on a Kubernetes
cluster. Ye have achieved the following objectives:

- Ye have installed Cilium on yer Kubernetes cluster and configured it to use SPIRE as the identity provider.
- Ye have tested the network policy enforcement using an example application that simulates a ship-to-shore
communication scenario.
- Ye have verified that the network policies are enforced by the SPIFFE identities issued by SPIRE.

By doing so, ye have made Coastal Containers’ ship-to-shore communications more secure, reliable and interoperable.
Ye have also gained valuable skills and knowledge that will help ye in yer future adventures on the high seas. Well
done, matey! Ye have earned yer stripes as a ship’s engineer.
