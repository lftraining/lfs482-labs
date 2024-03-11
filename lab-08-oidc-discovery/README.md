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

## Step-by-Step Instructions

### Step 1: Provision Infrastructure

Run the following command in the [lab-08-oidc-discovery](../lab-08-oidc-discovery) directory to set up the necessary
infrastructure, including [cert-manager](https://github.com/cert-manager/cert-manager),
[contour](https://github.com/projectcontour/contour) and a self-signed CA:

```shell
make cluster-up
```

### Step 2: Setup SPIRE with Helm

Before you deploy SPIRE via the Helm chart, you must first add the SPIFFE helm repo by running:

```shell
make spire-add-helm-repo
```

Once added, install SPIRE via helm by running:

```shell
make spire-helm-install spire-wait-for-agent
```

If everything worked properly, you should see:

```log
🏗️ Installing SPIRE using Helm...
NAME: spire-crds
LAST DEPLOYED: Thu Mar  7 15:07:12 2024
NAMESPACE: spire
STATUS: deployed
REVISION: 1
TEST SUITE: None
NAME: spire
LAST DEPLOYED: Thu Mar  7 15:07:13 2024
NAMESPACE: spire
STATUS: deployed
REVISION: 1
NOTES:
Installed spire…

Spire CR's will be handled only if className is set to "spire-spire"
✔️ SPIRE installed using Helm.
pod/spire-agent-qbc9q condition met
```

### Step 3: View SPIRE Configuration

With SPIRE deployed and running on your cluster, let's analyze the SPIRE Server and Agent configuration files.

#### SPIRE Server Configuration

To first inspect the `spire-server` config, run the following `make` command:

```shell
make spire-view-server-config
```

The key differences from what we have seen before are shown below.

```json
{
  "server": {
    "jwt_issuer": "127.0.0.1.nip.io"
  },
  "plugins": {
    "KeyManager": [
      {
        "disk": {
          "plugin_data": {
            "keys_path": "/run/spire/data/keys.json"
          }
        }
      }
    ]
  }
}
```

- The Helm chart defaults to using the
  [disk KeyManager](https://github.com/spiffe/spire/blob/v1.9.0/doc/plugin_server_keymanager_disk.md)
- We have configured the `jwt_issuer` field for the server

We have used `127.0.0.1.nip.io` as this will resolve to local host when we use curl to walk through the OIDC Discovery
protocol in [Step 5: OIDC Discovery Walkthrough](#step-5-oidc-discovery-walkthrough). [nip.io](https://nip.io/) is a
free service to map DNS names to any IP Address.

#### SPIRE Agent Configuration

Next, view the running `spire-agent` configuration by issuing:

```shell
make spire-view-agent-config
```

```json
{
  "plugins": {
    "KeyManager": [
      {
        "disk": {
          "plugin_data": {
            "keys_path": "/run/spire/data/keys.json"
          }
        }
      }
    ]
  }
}
```

- The Helm chart defaults to using the
  [disk KeyManager](https://github.com/spiffe/spire/blob/v1.9.0/doc/plugin_agent_keymanager_disk.md)

The rest of the configuration should look familiar.

### Step 4: OIDC Discovery Provider Configuration

The [SPIRE OIDC Discovery Provider](https://github.com/spiffe/spire/tree/main/support/oidc-discovery-provider) exposes
an endpoint that allows external services to verify JWT SVIDs provided by our SPIRE server. We will walk through how 
this works in the [next](#step-5-oidc-discovery-walkthrough) section.

View the OIDC Discovery Provider configuration file by running:

```shell
make view-oidc-discovery-provider-config
```

```json
{
  "allow_insecure_scheme": true,
  "domains": [
    "spire-spiffe-oidc-discovery-provider",
    "spire-spiffe-oidc-discovery-provider.spire",
    "spire-spiffe-oidc-discovery-provider.spire.svc.cluster.local",
    "127.0.0.1.nip.io"
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
    "trust_domain": "coastal-containers.example"
  }
}
```

In the domains we can see that in addition to the configured JWT Issuer (`127.0.0.1.nip.io`), the service also supports
the in cluster service DNS names.

### Step 5: OIDC Discovery Walkthrough

To get a solid understanding of the OIDC Discovery protocol and how this can be used by external systems to verify JWT 
SVIDs, we'll manually walk through the process.

First we need to create a couple of registration entries; one for the SPIRE agent, and a second for a generic workload:

```shell
make create-registration-entries
```

```shell
Entry ID         : caf1a7ab-a622-4d4c-bd4c-129081b72708
SPIFFE ID        : spiffe://coastal-containers.example/agent/spire-agent
Parent ID        : spiffe://coastal-containers.example/spire/server
Revision         : 0
X509-SVID TTL    : default
JWT-SVID TTL     : default
Selector         : k8s_psat:agent_ns:spire
Selector         : k8s_psat:agent_sa:spire-agent
Selector         : k8s_psat:cluster:kind-kind

Entry ID         : de1fb77d-2586-45ae-b497-6f1552569271
SPIFFE ID        : spiffe://coastal-containers.example/workload
Parent ID        : spiffe://coastal-containers.example/agent/spire-agent
Revision         : 0
X509-SVID TTL    : default
JWT-SVID TTL     : default
Selector         : k8s:ns:default
Selector         : k8s:sa:default
```

Next we'll generate get the SPIRE server to mint a JWT SVID directly using the cli (in practice workloads would do this
through the Workload API) and then we will view the claims.

```shell
JWT_SVID=$(kubectl exec -n spire spire-server-0 -- \
  bin/spire-server jwt mint -audience oidc-discovery -spiffeID spiffe://coastal-containers.example/workload)
jq -R 'split(".") | .[1] | @base64d | fromjson' <<< $JWT_SVID
```

We can see the token has:

- The requested audience: `oidc-discovery`
- The validity period: from `iat` to `exp`
- The workload subject: `spiffe://coastal-containers.example/workload`
- The token issuer: `https://127.0.0.1.nip.io`

```json
{
  "aud": [
    "oidc-discovery"
  ],
  "exp": 1709883690,
  "iat": 1709882790,
  "iss": "127.0.0.1.nip.io",
  "sub": "spiffe://coastal-containers.example/workload"
}
```

The [OIDC Discovery Spec](https://openid.net/specs/openid-connect-discovery-1_0-final.html), section 4 states that:

> OpenID Providers supporting Discovery MUST make a JSON document available at the path formed by concatenating the 
> string /.well-known/openid-configuration to the Issuer

The Ingress for our cluster is using a NodePort mapped to `8443` on the host, so we can retrieve the OIDC Discovery 
Document by running:

```shell
curl -sk https://127.0.0.1.nip.io:8443/.well-known/openid-configuration | jq
```

```json
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

The Discovery Document tells us how to retrieve the keys we need to verify the signature contained in the JWT in the
`jwks_uri` field.

```shell
curl -sk https://127.0.0.1.nip.io:8443/keys | jq
```

We can see this returns an array of `keys`. As the SPIRE server rotates the signing keys, those that are still within
their validity period will appear in this list.

```json
{
  "keys": [
    {
      "kty": "RSA",
      "kid": "OLDcvewt1bPTr5f33WXKXirPPUQbq0R8",
      "alg": "RS256",
      "n": "osfM5niRZzBL2pD6EVRHZOd0YwOo-BEyoT0rYJi1Fv0w1rjsda_objnixf_Nd4iQBLg2f1H3ttHkdBv3HMhqoKgbPV2Nd8yDhBcL2AhYVsyJLPyFvDnrEMA4jCRPs_52tjg_u9VgCk6OSrJGb1nhOIzzYOvtT2Zl90brMOTLIqHxtjDdtjJvL3a3t_JR_e_2bFq1hkzjSDF5y-B4acX4wtrvj9eEeTiCWYzlgLoi6wX2GwQ37t6Y1wwpvQ4qWIqFDL8-9yFVzqf0ZqDqmOnCHqbRHugawx5tB6ELuRi-PP-TaUxC0gN26j22Ysxzm10N4K6IXxqV4qf7ONvuHp6NTQ",
      "e": "AQAB"
    }
  ]
}
```

We can determine which key was used to sign our JWT from the headers

```shell
jq -R 'split(".") | .[0] | @base64d | fromjson' <<< $JWT_SVID
```

```json
{
  "alg": "RS256",
  "kid": "OLDcvewt1bPTr5f33WXKXirPPUQbq0R8",
  "typ": "JWT"
}
```

We can see that the `kid` field in our JWT SVID: `OLDcvewt1bPTr5f33WXKXirPPUQbq0R8` identifies the key contained in the
JWKS we retrieved, and we can use this to reconstitute the key, verify the signature and therefore the validation of the
presented JWT.

By exposing the OIDC Discovery Endpoint, we can present JWT SVIDs to any system that understands the OIDC Discovery 
protocol, allowing them to verify our token and extract our identity from the claims (`sub`) and the intended use for
the token (`aud`) enabling identity federation. An example use case for this is using SPIRE provided JWT SVIDs along 
with web identity federation to [obtain temporary](https://spiffe.io/docs/latest/keyless/oidc-federation-aws/) AWS 
credentials.

### Step 6: Controller Manager Configuration

As part of the helm deployment, we also installed the 
[SPIRE Controller Manager](https://github.com/spiffe/spire-controller-manager) a Kubernetes Controller that watches for
workloads running in the cluster and automatically creates registration entries in the SPIRE server.

View the Controller Manager configuration file by running:

```shell
make view-controller-manager-config
```

```yaml
apiVersion: spire.spiffe.io/v1alpha1
kind: ControllerManagerConfig
metadata:
  name: spire-controller-manager
  namespace: spire
  labels:
    helm.sh/chart: spire-server-0.1.0
    app.kubernetes.io/name: server
    app.kubernetes.io/instance: spire
    app.kubernetes.io/version: "1.9.1"
    app.kubernetes.io/managed-by: Helm
metrics:
  bindAddress: 0.0.0.0:8082
health:
  healthProbeBindAddress: 0.0.0.0:8083
leaderElection:
  leaderElect: true
  resourceName: 67103523.spiffe.io
  resourceNamespace: spire
validatingWebhookConfigurationName: spire-spire-controller-manager-webhook
clusterName: kind-kind
trustDomain: coastal-containers.example
ignoreNamespaces:
  - kube-system
  - kube-public
  - local-path-storage
spireServerSocketPath: "/tmp/spire-server/private/api.sock"
className: "spire-spire"
watchClassless: false
parentIDTemplate: "spiffe://{{ .TrustDomain }}/spire/agent/k8s_psat/{{ .ClusterName }}/{{ .NodeMeta.UID }}"
```

The key configuration item here is 
`parentIDTemplate: "spiffe://{{ .TrustDomain }}/spire/agent/k8s_psat/{{ .ClusterName }}/{{ .NodeMeta.UID }}"` which
templates the Parent ID for dynamic registration entries. As you can see, it creates an entry for the SPIRE agent for
the Node that the workload is deployed to.

The SPIFFE ID created for a workload is configured using the ClusterSPIFFEID custom resource. This can be viewed by 
running:

```shell
make view-spiffe-clusterid
```

```json
{
  "className": "spire-spire",
  "namespaceSelector": {
    "matchExpressions": [
      {
        "key": "kubernetes.io/metadata.name",
        "operator": "NotIn",
        "values": [
          "spire",
          "spire-server",
          "spire-system"
        ]
      }
    ]
  },
  "spiffeIDTemplate": "spiffe://{{ .TrustDomain }}/ns/{{ .PodMeta.Namespace }}/sa/{{ .PodSpec.ServiceAccountName }}"
}
```

For workloads deployed to any namespace other than spire, spire-server, and spire-system, the SPIFFE ID will be of the
format `spiffe://{{ .TrustDomain }}/ns/{{ .PodMeta.Namespace }}/sa/{{ .PodSpec.ServiceAccountName }}`. We will see this
in action in the next step where we deploy a workload, without creating a registration entry manually.

See the documentation for the [SPIRE Controller Manager](https://github.com/spiffe/spire-controller-manager) to 
understand how you can configure this to suit your requirements if the default configuration does not meet your needs.

### Step 7: Programmatic OIDC Discovery

In this step we repeat the previous steps programmatically and also provide an example of verifying the JWT SVID using
the [Go JOSE](https://github.com/go-jose/go-jose) library. 

Have a look at the [workload](workload/main.go) and notice it performs the following steps:

1. Create a Workload API client
2. Obtain a JWT SVID (this is performed in a retry to allow time for the controller manager to register a SPIFFE ID for 
the workload)
3. Print the raw JWT
4. Parse the JWT and print the headers and claims
5. Retrieve and print the OIDC Discovery document (we use `http://spire-spiffe-oidc-discovery-provider.spire` here so 
that it is resolvable from our workload in the cluster)
6. Extract the JWKS URI, download, and print the JWKS
7. Verify the JWT SVID using the correct JWK and print the verified claims

To build and load the provided workload image into your Kind cluster and deploy the workload as a Job, run:

```shell
make cluster-build-load-image deploy-workload DIR=workload
```

Once the job has finished, run the following command to view the logs:

```shell
make view-logs
```

This will provide output similar to this below:

```shell
2024/03/08 09:13:22 JWT SVID
eyJhbGciOiJSUzI1NiIsImtpZCI6ImJ6UzRISGg0azcyWWY2UFNjUFdldFhVWXA3N3YxdHRHIiwidHlwIjoiSldUIn0.eyJhdWQiOlsib2lkYy1kaXNjb3ZlcnkiXSwiZXhwIjoxNzA5ODkwMTAyLCJpYXQiOjE3MDk4ODkyMDIsImlzcyI6IjEyNy4wLjAuMS5uaXAuaW8iLCJzdWIiOiJzcGlmZmU6Ly9jb2FzdGFsLWNvbnRhaW5lcnMuZXhhbXBsZS9ucy9kZWZhdWx0L3NhL3dvcmtsb2FkIn0.dQA4ngBL7C4HdeXg2Tbn2e83tnkWUSxq2cGor8K_Eu3x6Wv_ZMGyM7p0a6MCXq2QULUk-0hKcVO_PDi7Xub4zDAsAfxGx8_oLHbtPpYTVHrEsJWFih5GZWUFwk-0Q4CbPT8SapbU2xXX4ATAyEeLeRqSflsLTKJwVp5VaCdWgSFO0oWhWOmOMeZwrFmdqI-VOLJD5Xz6EUXh_XSpuaXq-sD7OHd0PgOLD1vrtGBclSinT-OQkmfkPDgyK7evk0zqoTVdLj86V3I4YRTM3XncP5eRHpHGsuAtiUluDBghGliNKEptWq5wbxwk89bgxyS7NA3UEzTQs4DemxSc9u7kPw

2024/03/08 09:13:22 Parsed Token
Headers: 
alg: RS256
kid: bzS4HHh4k72Yf6PScPWetXUYp77v1ttG
Claims: 
iss: 127.0.0.1.nip.io
sub: spiffe://coastal-containers.example/ns/default/sa/workload
aud: [oidc-discovery]
iat: 2024-03-08 09:13:22 +0000 UTC
exp: 2024-03-08 09:28:22 +0000 UTC

2024/03/08 09:13:22 OIDC Discovery Document from http://spire-spiffe-oidc-discovery-provider.spire/.well-known/openid-configuration: 
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

2024/03/08 09:13:22 JSON Web Key Set: {
  "keys": [
    {
      "kty": "RSA",
      "kid": "bzS4HHh4k72Yf6PScPWetXUYp77v1ttG",
      "alg": "RS256",
      "n": "tl7rSTDhC9kWc00IRm9uBGYbDBPt56nkYYJGpZhSBDxTxsvjuQt-YnE17JKP3-rydgot_3bVqeBgOKyh7w8K-kj-nOndN8diGL4s9aS9Qz-hPpeZj2Mk-wFyeosSJH_ihxxeWLhvD2N3gXaDG5YTY5CFiy6-Iv4jkcrQ7t2m8B3bGCUpBQXy7bEeurOfVaWI8vqo7mjwBayblLZVCwx21stkyFaxhN2fsXeo74amS1ibkWYVb7LpHUPp9FuUM1bkKEz95r5aqIfQUDLa6Kd6SvvGyjgDWKlRVRFRKu4jcfDdwPgZUOcBJifls2dIiO24esTXtqqboew_mxbO8ga7iw",
      "e": "AQAB"
    }
  ]
}

2024/03/08 09:13:22 Verified claims:
{"aud":["oidc-discovery"],"exp":1709890102,"iat":1709889202,"iss":"127.0.0.1.nip.io","sub":"spiffe://coastal-containers.example/ns/default/sa/workload"}
```

### Step 8: Cleanup

To tear down the Kind cluster, run:

```shell
make cluster-down
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
