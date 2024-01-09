# Lab 11: Advanced Configuration 2 - Federated SPIRE

## Prerequisites

- A AMD, Intel, or Arm 64-bit Linux environment.
- Familiarity with Kubernetes manifests and `kubectl` commands is helpful.

## Introduction

Ahoy, matey! Welcome to the course on SPIRE Federation. In this lab, ye will learn how to set up a federated SPIRE configuration between two different trust domains: Coastal Containers Ltd and its new partner AirFreight Nexus Ltd. 

SPIRE Federation is a feature that allows SPIRE Servers to exchange trust bundles and authenticate workloads across different platforms and environments. This enables secure communication between micro-services that belong to different crews, regions or organizations. Ye will use Kubernetes for managing containerized applications, to run yer SPIRE Servers and workloads. 

### Preparing Your Environment

Before you cast off, prepare your ship to sail by setting up your working environment. If you haven't yet done so, make sure you've cloned the lab repository to your local system. After that, you'll be working from the [lab-11-federated-spire](../lab-11-federated-spire/) directory.

```bash
export LAB_DIR=$(pwd)
export PATH=$PATH:$(pwd)/../bin
```

## Step-by-Step Instructions

### Step 1: Boot two Kubernetes clusters

You will need two Kubernetes clusters, one for Coastal Containers and one for AirFreight Nexus, so first be sure to tear down any kind existing clusters before this step. Spin-up the clusters by running:

```shell
make deploy-clusters
```

### Step 2: Deploy SPIRE to both clusters

The SPIRE server needs some configuration to tell it to how to share the trust bundles between the two clusters and also to tell it where to fetch the trust bundles from the other cluster.

See the [spire-server-config-coastal-containers.yaml](manifests/spire-server-config-coastal-containers.yaml) and [spire-server-config-airfreight-nexus.yaml](manifests/spire-server-config-airfreight-nexus.yaml) files for the configuration.

```ini
server {
...
  federation {
    bundle_endpoint {
      address = "0.0.0.0"
      port = 8443
    }
    federates_with "airfreight-nexus.example" {
      bundle_endpoint_url = "https://airfreight-nexus-control-plane:8443"
      bundle_endpoint_profile "https_spiffe" {
        endpoint_spiffe_id = "spiffe://airfreight-nexus.example/spire/server"
      }
    }
  }
}
```

To simply deploy this, run:

```shell
make deploy-spire
```

Check if the spire servers can reach each other:

```shell
kubectl --context kind-airfreight-nexus run --rm -ti --restart=Never --image=wbitt/network-multitool network-test --command -- curl -k https://coastal-containers-control-plane:8443

kubectl --context kind-coastal-containers run --rm -ti --restart=Never --image=wbitt/network-multitool network-test --command -- curl -k https://airfreight-nexus-control-plane:8443
```

You may find that one or both of these fails if your deployments aren't fully ready, if this is the case, retry once they are ready. The expected output of these checks should look something like below.

```log
{
    "keys": [
        {
            "use": "x509-svid",
            "kty": "RSA",
            "n": "tN8rS8tUoLxx_DME72BEZycO0dtGKYafPd4S9i8cDLVE2V9kwgZXy3hIHhFjcRcTRmlleKooZPMYlrhJLV_mz3EMMGUCUBHDEcYMExugGpY4XImQKiu-nOqRPGXj_U2udUto5q3Gma4XNpA-lvt0tjnde-RDFbafc47RhHs53p-Z4DurXrZtNQ9OBb6yFQ29faeLK3deKyQyYh30VA5jaV2wFPn416eFTeppKD2UDOT7WPslNqEmdBivAYtTAikNtcXDDkLEWV512i4ikDGKg60xeO1n75IcoHVINO833prayjW7l5-5rFp-713CdufR9dPH_YhKwbuydbC5nv4CTQ",
            "e": "AQAB",
            "x5c": [
                "MIIDiDCCAnCgAwIBAgIQabyY6WZc0/Oi3kGuooQMMTANBgkqhkiG9w0BAQsFADBHMQswCQYDVQQGEwJVSzEZMBcGA1UEChMQQWlyRnJlaWdodCBOZXh1czEdMBsGA1UEAxMUQWlyRnJlaWdodCBOZXh1cyBMdGQwHhcNMjQwMTAyMTg0NTExWhcNMjQwMTAzMTg0NTIxWjBHMQswCQYDVQQGEwJVSzEZMBcGA1UEChMQQWlyRnJlaWdodCBOZXh1czEdMBsGA1UEAxMUQWlyRnJlaWdodCBOZXh1cyBMdGQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC03ytLy1SgvHH8MwTvYERnJw7R20Yphp893hL2LxwMtUTZX2TCBlfLeEgeEWNxFxNGaWV4qihk8xiWuEktX+bPcQwwZQJQEcMRxgwTG6AaljhciZAqK76c6pE8ZeP9Ta51S2jmrcaZrhc2kD6W+3S2Od175EMVtp9zjtGEeznen5ngO6tetm01D04FvrIVDb19p4srd14rJDJiHfRUDmNpXbAU+fjXp4VN6mkoPZQM5PtY+yU2oSZ0GK8Bi1MCKQ21xcMOQsRZXnXaLiKQMYqDrTF47WfvkhygdUg07zfemtrKNbuXn7msWn7vXcJ259H108f9iErBu7J1sLme/gJNAgMBAAGjcDBuMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBT6u4isg4Tcncjx0L/N/Y9Qx0n4ljAsBgNVHREEJTAjhiFzcGlmZmU6Ly9haXJmcmVpZ2h0LW5leHVzLmV4YW1wbGUwDQYJKoZIhvcNAQELBQADggEBAJAe4v4mYhuQXG40bMrIRLsMMqesAs798RlcW6nJFZzA4WNelb1MT7tDYEVm/dLeTA0lnFRbFl14q9ErjekS7M5aEn0vRcUlBa/30vXLejB9F4bXo/bXZF5z2g00PGouVOJF7frwGAuUA6zxib6PYi8lfaF9PVFK4cDZA77OEIGN4qxDWDbLQkcFXsPpbpgKAhR9NKUD0KBIJFKJ7Eb3+V8v2BckRI9Dt/bhWWy9ytfo2b2+odU9dmYOJv9QjQIXL1KRAlE/X1tWvHK9M3AN3hEVorLF4zZq1AFnwUelpU5VSNqQA16yQHjbXkV+clUIIeckAdUS0OHz1m5nZOtmCRM="
            ]
        },
        {
            "use": "jwt-svid",
            "kty": "RSA",
            "kid": "ZxEVKal5Wo5ojLab3sIWmdzUeiVfF9MP",
            "n": "zgYcrUzxfCWgm1yOcqGItCMozMy_WRwn5UZyRPLfNVBO-g_lQvyt-otPOu6jHG83wQhtywZXmXDUwRYE8NSnrY2ZoaM-E2k-00YE5PmcAR1In9yfTzWPOO9FFvAcvyLHVSWt_o0ZqJ_seDO-aWY2MTseVMiFhoQl8mnGLLBpHPR39AecRxRV-TMlyPa6zczfEf-AgkopmP1q_qTr64mjlF-kYdawQAmMxBRn3IwU6wqhE0RlQQJUSdkI-6wyBNFbKYWQlf22iCjAD72XnYZd5FRjp8b4UhvnuhYQMf7-gotqxd1dhT3_vee9CJ8ffHciXYOiEXwR2AzWpFZUWt9JxQ",
            "e": "AQAB"
        }
    ],
    "spiffe_sequence": 1,
    "spiffe_refresh_hint": 8641
}
```

### Step 3: Manually share Bootstrap Trust bundles between both clusters

You've configured the SPIRE Servers with the federation endpoint addresses, but merely configuring this is insufficient to establish federation functionality. In order for the SPIRE Servers to successfully retrieve trust bundles from each other, you must initially exchange their respective trust bundles.

This exchange is essential because it allows them to authenticate the SPIFFE identity of the federated server attempting to access the federation endpoint.

After the federation is successfully initialized, trust bundle updates are acquired through the federation endpoint API, utilizing the current trust bundle.

First you need to retrieve the bundles from both clusters and save them to your local machine:

```shell
kubectl --context=kind-airfreight-nexus exec -n spire spire-server-0 -c spire-server -- \
  /opt/spire/bin/spire-server bundle show -format spiffe > bundles/airfreight-nexus.example.bundle

kubectl --context=kind-coastal-containers exec -n spire spire-server-0 -c spire-server -- \
  /opt/spire/bin/spire-server bundle show -format spiffe > bundles/coastal-containers.example.bundle
```

*📝 Note: The `bundle show` commands should automatically create the `bundles` directory to store the trust bundle contents, if this does not occur and you encounter errors while running these commands, issue a `mkdir bundles` command (in the root [lab-11-federated-spire](./) dir) to create the directory before running the commands again.*

Now copy the bundles to the alternate clusters:

```shell
kubectl --context=kind-airfreight-nexus cp -n spire -c debug bundles/coastal-containers.example.bundle spire-server-0:/run/spire/data/coastal-containers.example.bundle

kubectl --context=kind-coastal-containers cp -n spire -c debug bundles/airfreight-nexus.example.bundle spire-server-0:/run/spire/data/airfreight-nexus.example.bundle
```

Next, load the bundles into the clusters:

```shell
kubectl --context=kind-airfreight-nexus exec -n spire spire-server-0 -c spire-server -- \
  /opt/spire/bin/spire-server bundle set -format spiffe -id spiffe://coastal-containers.example -path /run/spire/data/coastal-containers.example.bundle

kubectl --context=kind-coastal-containers exec -n spire spire-server-0 -c spire-server -- \
  /opt/spire/bin/spire-server bundle set -format spiffe -id spiffe://airfreight-nexus.example -path /run/spire/data/airfreight-nexus.example.bundle 
```

If the commands execute without error, and the bundles are loaded properly, you should see `bundle set.` after each operation.

### Step 4: Create workload registration entries that federate between the clusters

Run the following commands to create workload registration entries which federate between the clusters:

```shell
kubectl --context=kind-coastal-containers exec -n spire spire-server-0 -c spire-server -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://coastal-containers.example/manifest/workload/server \
    -parentID spiffe://coastal-containers.example/agent/spire-agent \
    -selector k8s:ns:default \
    -selector k8s:sa:server \
    -federatesWith spiffe://airfreight-nexus.example

kubectl --context=kind-airfreight-nexus exec -n spire spire-server-0 -c spire-server -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://airfreight-nexus.example/manifest/workload/client \
    -parentID spiffe://airfreight-nexus.example/agent/spire-agent \
    -selector k8s:ns:default \
    -selector k8s:sa:client \
    -federatesWith spiffe://coastal-containers.example
```

Note the `federatesWith` flag, which enables federation relationships between SVIDs in different trust domains.

### Step 5: Deploy a workload

Using slightly modified manifest client / server workloads (as seen in [lab-04-getting-svids](../lab-04-getting-svids/)), we will be deploying the [server](./workload/server/) in the `kind-coastal-containers` cluster, and the [client](./workload/client/) in the `kind-airfreight-nexus` cluster. Deploy this setup by issuing:

```shell
make deploy-workload
```

This command will build the docker images, load them into their respective kind clusters, and then deploy the workloads per the federated topology mentioned previously.

### Step 6: Verify Federation allows workloads to mutually verify each other

To verify that the federation is working, check the client logs by running:

```shell
kubectl --context=kind-airfreight-nexus logs -f deployments/client
```

If everything worked as expected, you should see:

```log
Received ship manifest: {'ship_name': 'SS Coastal Carrier', 'departure_port': 'London Gateway', 'arrival_port': 'Port Elizabeth', 'cargo': [{'type': 'electronics', 'quantity': 1000}, {'type': 'clothing', 'quantity': 2000}, {'type': 'food', 'quantity': 3000}]}
```

Additionally, in this lab, we enable federation by first deploying an `initContainer` for each workload to fetches its SVID, key, local bundle, and the foreign trust bundle, which it then writes to the `/tmp/` directory. In this way, we do not need to rely on SPIFFE Helper as in previous labs, given that we are using the SPIRE Agent binary to fetch this public key material.

Verify this worked propely by checking the `agent` container logs for the client first:

```shell
kubectl --context=kind-airfreight-nexus logs -f deployments/client -c agent
```

If this worked properly, you should see an output similar to:

```log
Received 1 svid after 515.739818ms

SPIFFE ID:              spiffe://airfreight-nexus.example/manifest/workload/client
SVID Valid After:       2024-01-02 18:54:42 +0000 UTC
SVID Valid Until:       2024-01-02 19:54:52 +0000 UTC
CA #1 Valid After:      2024-01-02 18:45:11 +0000 UTC
CA #1 Valid Until:      2024-01-03 18:45:21 +0000 UTC
[spiffe://coastal-containers.example] CA #1 Valid After:        2024-01-02 18:44:41 +0000 UTC
[spiffe://coastal-containers.example] CA #1 Valid Until:        2024-01-03 18:44:51 +0000 UTC

Writing SVID #0 to file /tmp/svids/svid.0.pem.
Writing key #0 to file /tmp/svids/svid.0.key.
Writing bundle #0 to file /tmp/svids/bundle.0.pem.
Writing federated bundle #0 for trust domain spiffe://coastal-containers.example to file /tmp/svids/federated_bundle.0.0.pem.
```

Next, do the same for the server workload by running:

```shell
kubectl --context=kind-coastal-containers logs -f deployments/server -c agent
```

If this worked properly, you should see an output similar to:

```log
Received 1 svid after 1.033741053s

SPIFFE ID:              spiffe://coastal-containers.example/manifest/workload/server
SVID Valid After:       2024-01-02 18:54:41 +0000 UTC
SVID Valid Until:       2024-01-02 19:54:51 +0000 UTC
CA #1 Valid After:      2024-01-02 18:44:41 +0000 UTC
CA #1 Valid Until:      2024-01-03 18:44:51 +0000 UTC
[spiffe://airfreight-nexus.example] CA #1 Valid After:  2024-01-02 18:45:11 +0000 UTC
[spiffe://airfreight-nexus.example] CA #1 Valid Until:  2024-01-03 18:45:21 +0000 UTC

Writing SVID #0 to file /tmp/svids/svid.0.pem.
Writing key #0 to file /tmp/svids/svid.0.key.
Writing bundle #0 to file /tmp/svids/bundle.0.pem.
Writing federated bundle #0 for trust domain spiffe://airfreight-nexus.example to file /tmp/svids/federated_bundle.0.0.pem.
```

These logs allow you to verify that the client / server workloads can validate each others SVIDs by fetching the foreign trust bundle contents via the SPIRE Agent binary, effectively establishing a federation relationship between the two SPIRE deployments.

### Step 7: Cleanup

Now that you've proved everything works, its time to scrub the decks and delete your clusters:

```shell
cd $LAB_DIR && make cluster-down
```

Additionally, delete the bundles directory by running:

```shell
rm -rf bundles
```

## Conclusion

Congratulations, sailor! Ye have completed the lab and learned how to set up SPIRE Federation on a Kubernetes cluster. Ye have achieved the following objectives:

- Ye have configured SPIRE Server to expose its SPIFFE Federation bundle endpoint using SPIFFE authentication.
- Ye have configured SPIRE Servers to fetch trust bundles from each other.
- Ye have bootstrapped federation between two SPIRE Servers using different trust domains.
- Ye have created registration entries for the workloads so that they can federate with other trust domain.

By doing so, ye have made Coastal Containers’ ship-to-shore communications more secure, reliable and interoperable with its partner AirFreight Nexus. Ye have also gained valuable skills and knowledge that will help ye in yer future adventures on the high seas and air! Well done, matey! Ye have earned yer stripes as a ship’s engineer, and made your aeronautical chums very happy in the process.
