CFSSL_VERSION ?= 1.6.5
HELM_VERSION ?= 3.17.0
ISTIO_VERSION ?= 1.22.8
KIND_VERSION ?= 0.26.0
KUBECTL_VERSION ?= 1.32.0
OPA_VERSION ?= 1.1.0
SPIRE_VERSION ?= 1.11.1
CILIUM_CLI_VERSION ?= v0.16.23
HUBBLE_CLI_VERSION ?= v1.16.6

KIND_IMAGE ?= kindest/node:v1.32.0@sha256:c48c62eac5da28cdadcf560d1d8616cfa6783b58f0d94cf63ad1bf49600cb027

GOLANG_BUILDER_IMAGE ?= golang:1.21.7-alpine3.19@sha256:163801a964d358d6450aeb51b59d5c807d43a7c97fed92cc7ff1be5bd72811ab
GOLANG_IMAGE ?= gcr.io/distroless/base-debian12@sha256:5eae9ef0b97acf7de819f936e12b24976b2d54333a2cf329615366e16ba598cd
PYTHON_IMAGE ?= gcr.io/distroless/python3-debian12@sha256:d1427d962660c43d476b11f9bb7d6df66001296bba9577e39b33d2e8897614cd

CONTOUR_IMAGE ?= ghcr.io/projectcontour/contour:v1.30.2@sha256:ba4fd19911bd603e3857fafdd17bd44c687ea1ab52b12f8ddc7971ba68f9b532
CONTOUR_IMAGE_TAG = projectcontour/contour:latest

ENVOY_IMAGE ?= docker.io/envoyproxy/envoy:v1.31.5@sha256:2ef22f7a72d3008d68bd4912f52b9a3d466d00db31a567dffa1b103e42eb9198
ENVOY_IMAGE_TAG = envoyproxy/envoy:latest

CERT_MANAGER_CAINJECTOR_IMAGE ?= quay.io/jetstack/cert-manager-cainjector:v1.16.3@sha256:e25e5f62648966d9c442c38ef3004efb60072069d91cf6f1a9a480c16550d09e
CERT_MANAGER_CAINJECTOR_IMAGE_TAG = jetstack/cert-manager-cainjector:latest
CERT_MANAGER_CONTROLLER_IMAGE ?= quay.io/jetstack/cert-manager-controller:v1.16.3@sha256:17c8f2d46fd16087c9ee92688304b3e38b32cbcb1f5885412e5a35e8852bc029
CERT_MANAGER_CONTROLLER_IMAGE_TAG = jetstack/cert-manager-controller:latest
CERT_MANAGER_WEBHOOK_IMAGE ?= quay.io/jetstack/cert-manager-webhook:v1.16.3@sha256:0eb32021bf2f30d923c439fd79b1b2fd3d8cf877e3c915f8a34c12803138a145
CERT_MANAGER_WEBHOOK_IMAGE_TAG = jetstack/cert-manager-webhook:latest
CERT_MANAGER_ACMESOLVER_IMAGE ?= quay.io/jetstack/cert-manager-acmesolver:v1.16.3@sha256:aa603e03e22897ab65a799126b854863ab717538ee46e36df42ed26f74a2ec7f

WAIT_FOR_IT_IMAGE ?= cgr.dev/chainguard/wait-for-it:latest@sha256:caead414307e81dbdd86d30662fdfe1b999dd4ce8a10fa667dab3438d0eed193
WAIT_FOR_IT_IMAGE_TAG = chainguard/wait-for-it:latest

SPIRE_SERVER_IMAGE ?= ghcr.io/spiffe/spire-server:1.11.1@sha256:f48d3d766d862de9e7232544737be0aa9d2ece600567fdcaf6439167ca9cef5f
SPIRE_SERVER_IMAGE_TAG = spiffe/spire-server:latest
SPIRE_AGENT_IMAGE ?= ghcr.io/spiffe/spire-agent:1.11.1@sha256:3468d6c06da4a36c2b85f61c116935cc32dfdeec669ee7b09ad4f49566cea881
SPIRE_AGENT_IMAGE_TAG = spiffe/spire-agent:latest
SPIRE_CONTROLLER_MANAGER_IMAGE ?= ghcr.io/spiffe/spire-controller-manager:0.6.0@sha256:1a6876c5c66d836a1697492f431e1ff4e52e4a6d61200574535af33e72b1668b
SPIRE_CONTROLLER_MANAGER_IMAGE_TAG = spiffe/spire-controller-manager:latest
SPIRE_OIDC_DISCOVERY_PROVIDER_IMAGE ?= ghcr.io/spiffe/oidc-discovery-provider:1.11.1@sha256:4d273da48fb48a14b2105688378f365466beb902260817f5ae7df65e3e9c95ac
SPIRE_OIDC_DISCOVERY_PROVIDER_IMAGE_TAG = spiffe/oidc-discovery-provider:latest

SPIFFE_CSI_DRIVER_IMAGE ?= ghcr.io/spiffe/spiffe-csi-driver:nightly@sha256:7cd432061281da00273c2c5612e3b8fc44abd0aba7c8175bdfcf200fba6ed131
SPIFFE_CSI_DRIVER_IMAGE_TAG = spiffe/spiffe-csi-driver:latest
SPIFFE_HELPER_IMAGE ?= ghcr.io/spiffe/spiffe-helper:0.9.0@sha256:47f18fcd2bf1c99cbb9bdc6f4348b08fb74572f89ddd0bcf103de2a1b12ab6d9
SPIFFE_HELPER_IMAGE_TAG = spiffe/spiffe-helper:latest

CSI_NODE_DRIVER_REGISTRAR_IMAGE ?=  registry.k8s.io/sig-storage/csi-node-driver-registrar:v2.9.4@sha256:ab6c4603d8da43d9f33361a5716552080e3f72badfad099d76fcf4b59e547d80
CSI_NODE_DRIVER_REGISTRAR_IMAGE_TAG = sig-storage/csi-node-driver-registrar:latest

OPA_IMAGE ?= openpolicyagent/opa:1.1.0@sha256:8433ff5fd7919a85f97b013358425a0ed169989ae38f8a5ff80fe1790cba7725
OPA_IMAGE_TAG = openpolicyagent/opa:latest
OPA_PROXY_INIT_IMAGE ?= openpolicyagent/proxy_init:v8@sha256:1a867cf60ebf79e7ce2b40abca679ad6833901e5b7bd02e581a3538c7c3c2f51
OPA_PROXY_INIT_IMAGE_TAG = openpolicyagent/proxy_init:latest
OPA_ENVOY_IMAGE ?= openpolicyagent/opa:1.1.0-envoy@sha256:3be6adc37c08ca5fe867dc70d4872e06688677752aa0a8996b14edba4e9c534d
OPA_ENVOY_IMAGE_TAG = openpolicyagent/opa-envoy:latest

POSTGRES_IMAGE ?= postgres:17.2@sha256:87ec5e0a167dc7d4831729f9e1d2ee7b8597dcc49ccd9e43cc5f89e808d2adae
POSTGRES_IMAGE_TAG = postgres:latest

ALPINE_CURL_IMAGE ?= alpine/curl:latest@sha256:31f1d54c59fb1bf1b4381b13cbce89bc22449008a664987c468d103699af8c47
ALPINE_CURL_IMAGE_TAG = alpine/curl:latest

NETWORK_MULTITOOL_IMAGE ?= wbitt/network-multitool:latest@sha256:d1137e87af76ee15cd0b3d4c7e2fcd111ffbd510ccd0af076fc98dddfc50a735
NETWORK_MULTITOOL_IMAGE_TAG = wbitt/network-multitool:latest

BUSYBOX_IMAGE ?= busybox:latest@sha256:a5d0ce49aa801d475da48f8cb163c354ab95cab073cd3c138bd458fc8257fbf1
BUSYBOX_IMAGE_TAG = busybox:latest

CILIUM_IMAGE ?= quay.io/cilium/cilium:v1.16.6@sha256:cf64df62897b071d5a9a005564ecbfb9124aa82a96957e329ce28a187864f113
CILIUM_IMAGE_TAG = cilium/cilium:latest
CILIUM_OPERATOR_IMAGE ?= quay.io/cilium/operator-generic:v1.16.6@sha256:5f8a26797f32f7646850047dbd3d7d252d9bf949c3add428b48e8d4e3dbc9400
CILIUM_OPERATOR_IMAGE_TAG = cilium/operator-generic:latest
CILIUM_CURL_IMAGE ?= quay.io/cilium/alpine-curl:v1.10.0@sha256:9c2638849cab7c7bb8f6f2175bacc7f728d4f11c08e562d1f9c5e3bac1ee8d9f
CILIUM_CURL_IMAGE_TAG = cilium/alpine-curl:latest
CILIUM_JSON_MOCK_IMAGE ?= quay.io/cilium/json-mock:v1.3.8@sha256:2b2968790f9d7dbdf3be84c6283da0204675e3baa635a57c659ac73ee43ae3c8
CILIUM_JSON_MOCK_IMAGE_TAG = cilium/json-mock:latest
HUBBLE_RELAY_IMAGE ?= quay.io/cilium/hubble-relay:v1.16.6@sha256:eb6b6e5f128484122974385c46d0e5aae820d628ab564684231a16e25e0e0a59
HUBBLE_RELAY_IMAGE_TAG = cilium/hubble-relay:latest
HUBBLE_UI_IMAGE ?= quay.io/cilium/hubble-ui:v0.13.1@sha256:306c8accedc8eb0dd31b128ce71767d8ff16d7623605e2559217ac563b8e17da
HUBBLE_UI_IMAGE_TAG = cilium/hubble-ui:latest
HUBBLE_UI_BACKEND_IMAGE ?= quay.io/cilium/hubble-ui-backend:v0.13.1@sha256:c7766e572e407f98b3790c26c6f1a3a2fcf91b172b097a426b91d60d9f3d7e89
HUBBLE_UI_BACKEND_IMAGE_TAG = cilium/hubble-ui-backend:latest

BIN_DIR ?= bin
CONFIG_DIR ?= ../config

OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell arch | sed s/aarch64/arm64/ | sed s/x86_64/amd64/)

CFSSL = $(shell pwd)/$(BIN_DIR)/cfssl
CFSSLJSON = $(shell pwd)/$(BIN_DIR)/cfssljson
HELM = $(shell pwd)/$(BIN_DIR)/helm
ISTIO_CTL = $(shell pwd)/$(BIN_DIR)/istioctl
KIND = $(shell pwd)/$(BIN_DIR)/kind
KUBECTL = $(shell pwd)/$(BIN_DIR)/kubectl
OPA = $(shell pwd)/$(BIN_DIR)/opa
SPIRE_SERVER = $(shell pwd)/$(BIN_DIR)/spire-server
SPIRE_AGENT = $(shell pwd)/$(BIN_DIR)/spire-agent
CILIUM_CLI = $(shell pwd)/$(BIN_DIR)/cilium
HUBBLE_CLI = $(shell pwd)/$(BIN_DIR)/hubble

KIND_CONFIG := $(if $(wildcard ./config/kind-config.yaml),"./config/kind-config.yaml",$(CONFIG_DIR)/kind-config.yaml)
ISSUER_CONFIG := $(if $(wildcard ./config/self-signed-cluster-issuer.yaml),"./config/self-signed-cluster-issuer.yaml",$(CONFIG_DIR)/self-signed-cluster-issuer.yaml)

##@ SPIRE

.PHONY: spire-helm-install
spire-helm-install: helm ## Install SPIRE using the Helm chart
	@echo "🏗️ Installing SPIRE using Helm..."
	@$(HELM) upgrade --install -n spire --create-namespace spire-crds spiffe/spire-crds
	@$(HELM) upgrade --install -n spire spire spiffe/spire -f values.yaml
	@echo "✔️ SPIRE installed using Helm."

.PHONY: spire-helm-uninstall
spire-helm-uninstall: helm ## Uninstall SPIRE using the Helm chart
	@$(HELM) uninstall -n spire spire

.PHONY: spire-view-server-config
spire-view-server-config: kubectl ## View the SPIRE Server configuration file.
	@$(KUBECTL) -n spire get cm spire-server -ojsonpath='{.data.server\.conf}'

.PHONY: spire-view-agent-config
spire-view-agent-config: kubectl ## View the SPIRE Agent configuration file.
	@$(KUBECTL) -n spire get cm spire-agent -ojsonpath='{.data.agent\.conf}'

.PHONY: spire-list-entries
spire-list-entries: kubectl ## List Registration Entries
	@$(KUBECTL) -n spire exec spire-server-0 -- bin/spire-server entry show

.PHONY: spire-wait-for-agent
spire-wait-for-agent: kubectl ## Wait for the SPIRE Agent to be running
	@$(KUBECTL) -n spire wait --for=condition=Ready --timeout=300s pods -l=app.kubernetes.io/name=agent

.PHONY: spire-add-helm-repo
spire-add-helm-repo: helm ## Add the SPIFFE helm repo
	@$(HELM) repo add spiffe https://spiffe.github.io/helm-charts-hardened/
	@$(HELM) repo update


##@ Infrastructure

.PHONY: cluster-up
cluster-up:: kind kubectl ## Launch a Kind cluster
	@$(KIND) create cluster --config=$(KIND_CONFIG)

.PHONY: cluster-down
cluster-down: kind ## Delete a Kind cluster
	@$(KIND) delete cluster

.PHONY: cluster-build-load-image
cluster-build-load-image: kind ## Build Docker image and load into cluster, use DIR=xxx to specify the tag and context directory
	@docker build \
		--build-arg GOLANG_BUILDER_IMAGE=$(GOLANG_BUILDER_IMAGE) \
		--build-arg GOLANG_IMAGE=$(GOLANG_IMAGE) \
		--build-arg PYTHON_IMAGE=$(PYTHON_IMAGE) \
		--build-arg SPIFFE_HELPER_IMAGE=$(SPIFFE_HELPER_IMAGE) \
		-t $(DIR) \
		$(DIR)
	@$(MAKE) cluster-load-image IMG=$(DIR)

cluster-load-image: kind ## Load image into kind cluster, IMG=xxx
	@$(KIND) load docker-image $(IMG)

##@ Tools

.PHONY: cfssl
cfssl: ## Download cfssl and cfssljson if required
ifeq (,$(wildcard $(CFSSL)))
		@{ \
				curl -sLo $(CFSSL) https://github.com/cloudflare/cfssl/releases/download/v$(CFSSL_VERSION)/cfssl_$(CFSSL_VERSION)_$(OS)_$(ARCH); \
				curl -sLo $(CFSSLJSON) https://github.com/cloudflare/cfssl/releases/download/v$(CFSSL_VERSION)/cfssljson_$(CFSSL_VERSION)_$(OS)_$(ARCH); \
				chmod +x $(CFSSL) $(CFSSLJSON); \
		}
endif

.PHONY: cfssl-version
cfssl-version: cfssl
	@$(CFSSL) version

.PHONY: cfssljson-version
cfssljson-version: cfssl
	@$(CFSSLJSON) -version

.PHONY: helm
helm: ## Download helm if required
ifeq (,$(wildcard $(HELM)))
	@{ \
		curl -sLO https://get.helm.sh/helm-v$(HELM_VERSION)-$(OS)-$(ARCH).tar.gz; \
		tar -C $(BIN_DIR) --strip-components=1 -xzf helm-v$(HELM_VERSION)-$(OS)-$(ARCH).tar.gz $(OS)-$(ARCH)/helm;\
		chmod +x $(HELM); \
		rm helm-v$(HELM_VERSION)-$(OS)-$(ARCH).tar.gz; \
	}
endif

helm-version: helm
	@$(HELM) version

.PHONY: istio
istio: ## Download istioctl if required
ifeq (,$(wildcard $(ISTIO_CTL)))
	@{ \
		curl -sLO https://github.com/istio/istio/releases/download/$(ISTIO_VERSION)/istio-$(ISTIO_VERSION)-$(OS)-$(ARCH).tar.gz; \
		tar -C $(BIN_DIR) --strip-components=2 -xzf istio-$(ISTIO_VERSION)-$(OS)-$(ARCH).tar.gz istio-$(ISTIO_VERSION)/bin; \
		rm istio-$(ISTIO_VERSION)-$(OS)-$(ARCH).tar.gz; \
	}
endif

istioctl-version: istio
	@$(ISTIO_CTL) version --remote=false

.PHONY: kind
kind: ## Download kind if required
ifeq (,$(wildcard $(KIND)))
	@{ \
		curl -sLo $(KIND) https://github.com/kubernetes-sigs/kind/releases/download/v$(KIND_VERSION)/kind-$(OS)-$(ARCH) ; \
		chmod +x $(KIND); \
	}
endif

kind-version: kind
	@$(KIND) version

.PHONY: kubectl
kubectl: ## Download kubectl if required
ifeq (,$(wildcard $(KUBECTL)))
	@{ \
		curl -sLo $(KUBECTL) https://dl.k8s.io/release/v$(KUBECTL_VERSION)/bin/$(OS)/$(ARCH)/kubectl ; \
		chmod +x $(KUBECTL); \
	}
endif

kubectl-version: kubectl
	@$(KUBECTL) version --client=true

.PHONY: opa
opa: ## Download opa if required
ifeq (,$(wildcard $(OPA)))
	@{ \
		curl -sLo $(OPA) https://github.com/open-policy-agent/opa/releases/download/v$(OPA_VERSION)/opa_$(OS)_$(ARCH); \
		chmod +x $(OPA); \
	}
endif

opa-version: opa
	@$(OPA) version

.PHONY: spire
spire: ## Download spire-server and spire-agent if required
ifeq (,$(wildcard $(SPIRE_SERVER)))
	@{ \
		curl -sLO https://github.com/spiffe/spire/releases/download/v$(SPIRE_VERSION)/spire-$(SPIRE_VERSION)-linux-$(ARCH)-musl.tar.gz; \
		tar -C $(BIN_DIR) --strip-components=2 -xzf spire-$(SPIRE_VERSION)-linux-$(ARCH)-musl.tar.gz spire-$(SPIRE_VERSION)/bin; \
		rm spire-$(SPIRE_VERSION)-linux-$(ARCH)-musl.tar.gz; \
	}
endif

spire-server-version: spire
	@$(SPIRE_SERVER) --version

spire-agent-version: spire
	@$(SPIRE_AGENT) --version

.PHONY: cilium
cilium: ## Download cilium cli
ifeq (,$(wildcard $(CILIUM_CLI)))
	@curl -sL https://github.com/cilium/cilium-cli/releases/download/$(CILIUM_CLI_VERSION)/cilium-$(OS)-$(ARCH).tar.gz | tar xz -C $(BIN_DIR)
	@chmod +x $(CILIUM_CLI)
endif
ifeq (,$(wildcard $(HUBBLE_CLI)))
	@curl -sL https://github.com/cilium/hubble/releases/download/$(HUBBLE_CLI_VERSION)/hubble-$(OS)-$(ARCH).tar.gz | tar xz -C $(BIN_DIR)
	@chmod +x $(HUBBLE_CLI)
endif
