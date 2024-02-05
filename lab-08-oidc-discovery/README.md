# Lab 08: OpenID Connect Discovery

## Prerequisites

- A 64-bit Linux environment (AMD, Intel, or Arm).
- Basic familiarity with `curl` commands is helpful.

## Introduction

Ahoy, digital navigators! 🏴‍☠️ Welcome aboard Coastal Containers Ltd, a renowned trans-Atlantic freight shipping company.
As we set sail on this voyage, we find ourselves amidst turbulent waters. The legacy shipping systems, once the pride of
our fleet, now face threats from modern-day digital pirates. But fear not! For we have embarked on a quest to modernize
and fortify our defenses using the formidable technologies of SPIFFE and SPIRE, aiming to implement a zero-trust
security model.

In this adventurous lab, our main objective is to chart the course through the intricate waves of
[OpenID Connect (OIDC)](https://auth0.com/docs/authenticate/protocols/openid-connect-protocol#) Discovery and its
integration with SPIRE. By the end of this journey, you'll should grasp:

- 🚢 Anchoring SPIRE with Helm: Setting up our trusty SPIRE using the Helm package manager.
- 🧭 Navigating SPIRE's Configuration: Gaining insights into the SPIRE Server and Agent configurations.
- 🌊 Understanding the Digital Ocean: Grasping the concepts of JWTs, Ingress, Certificates, and DNS in the vast sea of
SPIRE.
- ⚓ Deploying and Verifying Workloads: Ensuring our workload 'cargo' is safely secured within the SPIRE environment.

Our mission is clear: to integrate OIDC Discovery with SPIRE, ensuring that the treasured cargo of Coastal Containers
Ltd. is safeguarded against swashbucklin' pirates up to no good! So, brace yourselves, for we are about to embark on a
thrilling journey where technology meets adventure, and where we, the crew of Coastal Containers Ltd., strive to secure
our legacy for the future generations to come! 🌊🔐🏴‍☠️

### Preparing Your Environment

Before you cast off, prepare your ship to sail by setting up your working environment. If you haven't yet done so, make
sure you've cloned the lab repository to your local system. After that, you'll be working from the
[lab-08-oidc-discovery](../lab-08-oidc-discovery) directory.

```bash
export LAB_DIR=$(pwd)
```

## Step-by-Step Instructions

### Step 1: Provision Infrastructure

Run the following command in the [lab-08-oidc-discovery](../lab-08-oidc-discovery) directory to set up the necessary
infrastructure, including [cert-manager](https://github.com/cert-manager/cert-manager),
[contour](https://github.com/projectcontour/contour) and a self-signed CA:

```shell
make cluster-up
```

You can skip this step if you have already setup the cluster.

### Step 2: Setup SPIRE with Helm

Before you deploy SPIRE via the Helm chart, you must first add the SPIFFE helm repo by running:

```shell
make spire-add-helm-repo
```

Once added, install SPIRE via helm by running:

```shell
make spire-helm-install
```

If everything worked properly, you should see:

```log
NAME: spire
LAST DEPLOYED: Tue Oct 24 14:39:08 2023
NAMESPACE: spire
STATUS: deployed
REVISION: 1
NOTES:
Installed spire…
```

To further ensure SPIRE is up and running, run the following `kubectl` command:

```shell
kubectl get pods -n spire
```

You should see the running `spire-server`, `spire-agent`, `spire-spiffe-csi-driver`, and
`spire-spiffe-oidc-discovery-provider` pods as shown here:

```log
NAME                                                    READY   STATUS    RESTARTS   AGE
spire-agent-qnmsk                                       1/1     Running   0          2m32s
spire-server-0                                          2/2     Running   0          2m32s
spire-spiffe-csi-driver-25cxx                           2/2     Running   0          2m32s
spire-spiffe-oidc-discovery-provider-5ccc6bb54c-kjfct   2/2     Running   0          2m32s
```

### Step 3: View SPIRE Configuration

With SPIRE deployed and running on your cluster, let's analyze the SPIRE Server and Agent configuration files.

To first inspect the `spire-server` config, run the following `make` command:

```shell
make spire-view-server-config
```

#### SPIRE Server Configuration

The SPIRE Server configuration provides the foundation for the SPIRE infrastructure. Here's a breakdown:

- `health_checks`: This section defines the health check parameters for the SPIRE server. It specifies the bind address,
port, and the paths for liveness and readiness probes.
- `plugins`: This section outlines the plugins used by the SPIRE server. It includes:
  - `DataStore`: Specifies the type of database (SQLite3) and its connection string.
  - `KeyManager`: Defines where the keys are stored (on disk) and their path.
  - `NodeAttestor`: Configures the Kubernetes PSAT (Projected Service Account Token) attestor and the allowed service accounts.
  - `Notifier`: Configures the Kubernetes bundle notifier, specifying the config map and namespace.
- `server`: This section contains the main configuration for the SPIRE server:
  - `bind_address` and `bind_port`: The address and port where the SPIRE server will listen.
  - `ca_key_type`: The type of key used for the Certificate Authority (CA).
  - `ca_subject`: The subject for the CA certificate.
  - `ca_ttl`: The time-to-live (TTL) for the CA certificate.
  - `data_dir`: The directory where SPIRE data is stored.
  - `default_jwt_svid_ttl`: The default TTL for JWT SVIDs.
  - `default_x509_svid_ttl`: The default TTL for X.509 SVIDs.
  - `jwt_issuer`: The issuer for JWTs, set to `127.0.0.1.nip.io` in our case.
  - `log_level`: The logging level for the SPIRE server.
  - `trust_domain`: The trust domain for the SPIRE server.

Next, view the running `spire-agent` configuration by issuing:

```shell
make spire-view-agent-config
```

#### SPIRE Agent Configuration

The SPIRE Agent configuration connects workloads to the SPIRE Server. Key components include:

- `agent`: Main configuration for the SPIRE agent, including the server's address and port, the agent's data directory,
and the trust domain.
- `health_checks`: Defines the health check parameters for the SPIRE agent.
- `plugins`: Outlines the plugins used by the SPIRE agent:
  - `KeyManager`: Specifies that keys are stored in memory.
  - `NodeAttestor`: Configures the Kubernetes PSAT attestor.
  - `WorkloadAttestor`: Configures the Kubernetes workload attestor.

### Step 4: OIDC Discovery Provider Configuration

With the `spire-server` and `spire-agent` configuration out of the way, you can now view the OIDC Discovery Provider
configuration file by running:

```shell
make view-oidc-discovery-provider-config
```

This command will output the contents of the `spire-spiffe-oidc-discovery-provider` configmap used to run the
`spire-spiffe-oidc-discovery-provider` pod. Reference the output below to see what this should look like.

```log
{
  "allow_insecure_scheme": true,
  "domains": [
    "spire-spiffe-oidc-discovery-provider",
    "spire-spiffe-oidc-discovery-provider.spire",
    "spire-spiffe-oidc-discovery-provider.spire.svc.cluster.local",
    "127.0.0.1.nip.io",
    "localhost"
  ],
  "health_checks": {
    "bind_port": "8008",
    "live_path": "/live",
    "ready_path": "/ready"
  },
  "listen_socket_path": "/run/spire/oidc-sockets/spire-oidc-server.sock",
  "log_level": "info",
  "workload_api": {
    "socket_path": "/spiffe-workload-api/spire-agent.sock",
    "trust_domain": "coastal-containers.io"
  }
}
```

The OIDC Discovery Provider configuration is crucial for integrating SPIRE with OIDC. Key components include:

- `allow_insecure_scheme`: Allows insecure HTTP for local testing.
- `domains`: Lists the domains that the OIDC Discovery Provider will respond to. In our configuration, we're using
`127.0.0.1.nip.io` as the domain. [nip.io](https://nip.io/) is a service that provides a straightforward way to map
IP addresses to hostnames. Instead of manually editing the `etc/hosts` file with custom hostname and IP address mappings,
nip.io automates this process. It supports various formats, including:
  - Without a name:
    - `10.0.0.1.nip.io` maps to `10.0.0.1`.
    - `192-168-1-250.nip.io` maps to `192.168.1.250`.
    - `0a000803.nip.io` maps to `10.0.8.3`.
  - With a name:
    - `app.10.8.0.1.nip.io` maps to `10.8.0.1`.
    - `customer1.app.10.0.0.1.nip.io` maps to `10.0.0.1`.
    - `customer2-app-127-0-0-1.nip.io` maps to `127.0.0.1`.
- [nip.io](https://nip.io/) can map any IP address in "dot", "dash", or "hexadecimal" notation to the corresponding IP
address. For instance:
  - Dot notation: `magic.127.0.0.1.nip.io`.
  - Dash notation: `magic-127-0-0-1.nip.io`.
  - Hexadecimal notation: `magic-7f000001.nip.io`.
- The "dash" and "hexadecimal" notations are particularly useful when using services like
[LetsEncrypt](https://letsencrypt.org/), as they are treated as regular subdomains of [nip.io](https://nip.io/). This
service is open-source and is powered by [PowerDNS](https://www.powerdns.com/) with a custom
[PipeBackend](https://doc.powerdns.com/authoritative/backends/pipe.html). It's a free service provided by Exentrique
Solutions. The primary advantage of using [nip.io](https://nip.io/) in our configuration is that it allows us to route
traffic to our local environment (like Docker or a Kind cluster) via the [nip.io](https://nip.io/) URL without any
manual configuration.
- `health_checks`: Defines the health check parameters.
- `listen_socket_path`: The path to the socket where the OIDC Discovery Provider listens.
- `workload_api`: Configures the connection to the SPIRE Agent, specifying the socket path and trust domain.

### Step 5: Wait for SPIRE Agent

To ensure the next steps will work properly, wait for the SPIRE Agent to be running by issuing:

```shell
make spire-wait-for-agent
```

If your `spire-agent` is running, you should see:

```log
pod/spire-agent-qnmsk condition met
```

Keep in mind that the `qnmsk` identifier at the tail-end of the `spire-agent` pod name is subject to change per your own
deployment of the SPIRE pods.

### Step 6: OIDC Discovery Document

Now, view the OIDC Discovery Document by running:

```shell
curl -sk https://127.0.0.1.nip.io:8443/.well-known/openid-configuration | jq
```

This `curl` command will display the OpenID configuration from our JWT issuer at `127.0.0.1.nip.io`. The provided output
should look like:

```log
{
  "issuer": "https://127.0.0.1.nip.io",
  "jwks_uri": "https://127.0.0.1.nip.io/keys",
  "authorization_endpoint": "",
  "response_types_supported": [
    "id_token"
  ],
  "subject_types_supported": [],
  "id_token_signing_alg_values_supported": [
    "RS256",
    "ES256",
    "ES384"
  ]
}
```

The OIDC Discovery Document is a standard JSON object that describes the OIDC provider's configuration. It includes:

- `issuer`: The URL of the OIDC provider. We will use `https://127.0.0.1.nip.io` for our demonstration.
- `jwks_uri`: The URL of the provider's JSON Web Key Set (JWKS). We will use `https://127.0.0.1.nip.io/keys` for our
demonstration.
- `authorization_endpoint`: The URL of the authorization endpoint.
- `response_types_supported`: Lists the OIDC response types supported by the provider. This is `id_token` in our case.
- `subject_types_supported`: Lists the OIDC subject types supported.
- `id_token_signing_alg_values_supported`: Lists the signing algorithms supported for ID tokens. This is `RS256`,
`ES256`, and `ES384` in our case.

### Step 7: View the JSON Web Key Set (JWKS)

View the JWKS provided by our JWT issuer, by running:

```shell
curl -sk https://127.0.0.1.nip.io:8443/keys | jq
```

The output of this command should look like:

```log
{
  "keys": [
    {
      "kty": "RSA",
      "kid": "ThLDUEs6QyVMMAgTjRIFxgcdxPHmnID7",
      "alg": "RS256",
      "n": "2d7TMSZQ_aUBxJxK9Y_986lrpZznYSTIs_lEj3dswXI2kknYjAPucHX1MLZvt1gh-v4IMVHyygdgPHni9XWM7yaOwLsBA888KWxAwfliOTRFa3Q9mkazsrJpR4ijJR5lCbsv0ISNFucnPAXtAgjde2ox9stpu9wNiSTfFkTbv8_vvBYzq_qlskmw_gOouLOGscWSPz7Gsi8hFKVX09aNnEZy53S58TkuNBqn4LFklxsNhk3WRyfUhXo-pA6B4DfchXBxvClOusySUpKXpp0877HKtpdpUwPn4u2scQr0D-O9Z6GFCk8f5w0N0B5tFEfXYQUK3fHaWswzE9y2EcSYrw",
      "e": "AQAB"
    }
  ]
}
```

JWKS is a set of keys containing the public keys used to verify any JSON Web Token (JWT) issued by the authorization
server (`https://127.0.0.1.nip.io/keys` in our case). The JWKS contains:

- `kty`: The key type.
- `kid`: The key ID.
- `alg`: The algorithm used.
- `n`: The modulus for the RSA public key.
- `e`: The exponent for the RSA public key.

### Step 8: Deploy the Workload

To build and load the provided workload image into your Kind cluster, run:

```shell
make cluster-build-load-image DIR=workload
```

Once it is successfully loaded onto your `kind-control-plane` and `kind-worker` nodes, you can deploy the workload as a
Kubernetes [job](https://kubernetes.io/docs/concepts/workloads/controllers/job/). To do so, run:

```shell
make deploy-workload
```

If successful, this will provide the output:

```log
/../../zero-trust-labs/ilt/lab-08-oidc-discovery/../bin/kubectl apply -f workload/job.yaml
serviceaccount/workload created
job.batch/workload created
```

The workload, as defined in [main.go](./workload/main.go) and [job.yaml](./workload/job.yaml), is a Kubernetes job that
simulates a service within Coastal Containers Ltd. It fetches a JWT SVID from the SPIRE Agent, parses the token,
retrieves the OIDC Discovery Document, fetches the JWKS, and verifies the JWT's claims.

### Step 9: Analyze Workload Logs

Once the job has finished, run the following command to analyze the logs:

```shell
make view-logs
```

This will provide the output:

```log
/../../zero-trust-labs/ilt/lab-08-oidc-discovery/../bin/kubectl logs jobs/workload
2023/10/24 19:08:07 JWT SVID
eyJhbGciOiJSUzI1NiIsImtpZCI6IlRoTERVRXM2UXlWTU1BZ1RqUklGeGdjZHhQSG1uSUQ3IiwidHlwIjoiSldUIn0.eyJhdWQiOlsibGFiLTA5Il0sImV4cCI6MTY5ODE3NTM4NywiaWF0IjoxNjk4MTc0NDg3LCJpc3MiOiIxMjcuMC4wLjEubmlwLmlvIiwic3ViIjoic3BpZmZlOi8vY29hc3RhbC1jb250YWluZXJzLmlvL25zL2RlZmF1bHQvc2Evd29ya2xvYWQifQ.tQssFhWymviYOot4tQlwppfuR84Yx9NiNeLFRCPS6E240ixmk4FvafldmUvprgGGJZgBiNKfGjd4uf5njn_Avh1pJgQZgRlRF_ZiMPRpYCG6eXXehVHYyguBD49ICiYXec0Va-INxB8Y7eFNhn78-Fp0Dre_CIFahOwWtoZ0MUwiZmpwwdZRHYPI-gOXpaB6IqufgvcvAqj_HNy-3k7oVZnb84v7_1VKwhKHWU8qWL_qplmWW6u4ICMe_Cj1uWRtXo5lsBeL885O2u7MLLBZ9qubC0qd8g0GX4SMvNg5LKeh_6Cx5MbIu6QcqlDRHDfEwnmQESehAX9pGZSalOnVZQ

2023/10/24 19:08:07 Parsed Token
Headers:
alg: RS256
kid: ThLDUEs6QyVMMAgTjRIFxgcdxPHmnID7
Claims:
iss: 127.0.0.1.nip.io
sub: spiffe://coastal-containers.io/ns/default/sa/workload
aud: [lab-08]
iat: 2023-10-24 19:08:07 +0000 UTC
exp: 2023-10-24 19:23:07 +0000 UTC

2023/10/24 19:08:07 OIDC Discovery Document from http://spire-spiffe-oidc-discovery-provider.spire/.well-known/openid-configuration:
{
  "issuer": "http://spire-spiffe-oidc-discovery-provider.spire",
  "jwks_uri": "http://spire-spiffe-oidc-discovery-provider.spire/keys",
  "authorization_endpoint": "",
  "response_types_supported": [
    "id_token"
  ],
  "subject_types_supported": [],
  "id_token_signing_alg_values_supported": [
    "RS256",
    "ES256",
    "ES384"
  ]
}

2023/10/24 19:08:07 JSON Web Key Set: {
  "keys": [
    {
      "kty": "RSA",
      "kid": "ThLDUEs6QyVMMAgTjRIFxgcdxPHmnID7",
      "alg": "RS256",
      "n": "2d7TMSZQ_aUBxJxK9Y_986lrpZznYSTIs_lEj3dswXI2kknYjAPucHX1MLZvt1gh-v4IMVHyygdgPHni9XWM7yaOwLsBA888KWxAwfliOTRFa3Q9mkazsrJpR4ijJR5lCbsv0ISNFucnPAXtAgjde2ox9stpu9wNiSTfFkTbv8_vvBYzq_qlskmw_gOouLOGscWSPz7Gsi8hFKVX09aNnEZy53S58TkuNBqn4LFklxsNhk3WRyfUhXo-pA6B4DfchXBxvClOusySUpKXpp0877HKtpdpUwPn4u2scQr0D-O9Z6GFCk8f5w0N0B5tFEfXYQUK3fHaWswzE9y2EcSYrw",
      "e": "AQAB"
    }
  ]
}

2023/10/24 19:08:07 Verified claims:
{"aud":["lab-08"],"exp":1698175387,"iat":1698174487,"iss":"127.0.0.1.nip.io","sub":"spiffe://coastal-containers.io/ns/default/sa/workload"}
```

The logs provide a step-by-step breakdown of the workload's operations:

- **JWT SVID**: The JWT token fetched from the SPIRE Agent.
- **Parsed Token**: The parsed JWT token showing headers and claims.
- **OIDC Discovery Document**: The fetched OIDC Discovery Document.
- **JSON Web Key Set**: The fetched JWKS.
- **Verified claims**: The JWT's claims after verification against the JWKS.

### Step 10: Cleanup

As the following lab exercises will use the same cluster, tear down the workload job and uninstall SPIRE using helm by
running:

```shell
cd $LAB_DIR && make tear-down spire-helm-uninstall
```

To tear down the entire Kind cluster, run:

```shell
cd $LAB_DIR && make cluster-down
```

## Conclusion

As we dock our ship at the end of this thrilling voyage, we've successfully navigated the intricate waves of OpenID
Connect (OIDC) Discovery and its integration with SPIRE. We've fortified the defenses of Coastal Containers Ltd.,
ensuring that our treasured cargo remains safe from modern-day digital pirates.

Throughout this lab, we've anchored our understanding of SPIRE's configuration, delved into the OIDC Discovery
Provider's significance, and grasped the importance of the OIDC Discovery Document and JWKS. We've also witnessed the
power of JWTs in action, ensuring secure communication between our workloads.

OIDC, with its robust set of features, provides a secure and scalable solution for authentication. By integrating it
with SPIRE, we've taken a significant step towards a zero-trust security model, ensuring that only verified workloads
can communicate within our environment.

As we set our sights on future adventures, let's remember the lessons learned on this journey. The seas of technology
are vast and ever-changing, but with tools like SPIRE and OIDC, we're well-equipped to face any challenge that comes
our way. Safe sailing, digital navigators! 🏴‍☠️🌊🔐
