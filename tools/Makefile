CFSSL_VERSION ?= 1.6.4
HELM_VERSION ?= 3.12.3
ISTIO_VERSION ?= 1.19.0
KIND_VERSION ?= 0.20.0
KUBECTL_VERSION ?= 1.28.2
OPA_VERSION ?= 0.56.0
SPIRE_VERSION ?= 1.8.7
CILIUM_CLI_VERSION ?= v0.15.8
HUBBLE_CLI_VERSION ?= v0.12.0

KIND_IMAGE ?= kindest/node:v1.27.3@sha256:3966ac761ae0136263ffdb6cfd4db23ef8a83cba8a463690e98317add2c9ba72

GOLANG_BUILDER_IMAGE ?= golang:1.21.7-alpine3.19@sha256:163801a964d358d6450aeb51b59d5c807d43a7c97fed92cc7ff1be5bd72811ab
GOLANG_IMAGE ?= gcr.io/distroless/base-debian12@sha256:5eae9ef0b97acf7de819f936e12b24976b2d54333a2cf329615366e16ba598cd
PYTHON_IMAGE ?= gcr.io/distroless/python3-debian12@sha256:d1427d962660c43d476b11f9bb7d6df66001296bba9577e39b33d2e8897614cd

CONTOUR_IMAGE ?= ghcr.io/projectcontour/contour:v1.26.2@sha256:96647f464335ad930d97f9e16b5701440ea1fdce2455aee6f3747b0236deceee
CONTOUR_IMAGE_TAG = projectcontour/contour:latest

ENVOY_IMAGE ?= docker.io/envoyproxy/envoy:v1.27.3@sha256:72ca37b8da6575c5032a26ce58fbf82827dbc2ad2fd670171b05cf61e9dc4a32
ENVOY_IMAGE_TAG = envoyproxy/envoy:latest

CERT_MANAGER_CAINJECTOR_IMAGE ?= quay.io/jetstack/cert-manager-cainjector:v1.13.0@sha256:cb3c638088ef4308a67973bc48dbd6f25ae6c39e5157315c89fb6010fac40fd6
CERT_MANAGER_CAINJECTOR_IMAGE_TAG = jetstack/cert-manager-cainjector:latest
CERT_MANAGER_CONTROLLER_IMAGE ?= quay.io/jetstack/cert-manager-controller:v1.13.0@sha256:2547fde4e736101abf33f8c2503f12aa3a0b42614d3d64cfecf2835c0ee81c10
CERT_MANAGER_CONTROLLER_IMAGE_TAG = jetstack/cert-manager-controller:latest
CERT_MANAGER_WEBHOOK_IMAGE ?= quay.io/jetstack/cert-manager-webhook:v1.13.0@sha256:f8cd7c856c85795291319c94bd73dc88f195da26e1dc6f62dbae10887f5a1c8c
CERT_MANAGER_WEBHOOK_IMAGE_TAG = jetstack/cert-manager-webhook:latest
CERT_MANAGER_ACMESOLVER_IMAGE ?= quay.io/jetstack/cert-manager-acmesolver:v1.13.0@sha256:f4dea02438c66786bae8cb9c4f8d0f429c54ebee707e3394db69bbc95ea581a1

WAIT_FOR_IT_IMAGE ?= cgr.dev/chainguard/wait-for-it:latest@sha256:caead414307e81dbdd86d30662fdfe1b999dd4ce8a10fa667dab3438d0eed193
WAIT_FOR_IT_IMAGE_TAG = chainguard/wait-for-it:latest

SPIRE_SERVER_IMAGE ?= ghcr.io/spiffe/spire-server:1.8.7@sha256:ce3d81259b090eb1863d9f3caa106ecca418b19832d7e876abbe9be57cb4f0e7
SPIRE_SERVER_IMAGE_TAG = spiffe/spire-server:latest
SPIRE_AGENT_IMAGE ?= ghcr.io/spiffe/spire-agent:1.8.7@sha256:dc83e448e1ce8a6c37b9989679426deae409f0a159257f014797698648c4486e
SPIRE_AGENT_IMAGE_TAG = spiffe/spire-agent:latest
SPIRE_CONTROLLER_MANAGER_IMAGE ?= ghcr.io/spiffe/spire-controller-manager:0.4.3@sha256:740a46d91e182b98b9fd881427b9cb4fa47c29270a572342f21aebcdae7bef83
SPIRE_CONTROLLER_MANAGER_IMAGE_TAG = spiffe/spire-controller-manager:latest
SPIRE_OIDC_DISCOVERY_PROVIDER_IMAGE ?= ghcr.io/spiffe/oidc-discovery-provider:1.9.1@sha256:eb4a6bc3a747cee6e3d7bbb04f5795df6d0664064fe38b3a1c0dc0d937dd273a
SPIRE_OIDC_DISCOVERY_PROVIDER_IMAGE_TAG = spiffe/oidc-discovery-provider:latest

SPIFFE_CSI_DRIVER_IMAGE ?= ghcr.io/spiffe/spiffe-csi-driver:nightly@sha256:9f300d23f2d04b55e2e5705a473307211df3a139f0ae5579a0dc939a58a78bc2
SPIFFE_CSI_DRIVER_IMAGE_TAG = spiffe/spiffe-csi-driver:latest
SPIFFE_HELPER_IMAGE ?= ghcr.io/spiffe/spiffe-helper:nightly@sha256:baeb815b14f9217105e0e00fcb4d2ded7ccf06c37230860a40c7a26a8865e558
SPIFFE_HELPER_IMAGE_TAG = spiffe/spiffe-helper:latest

CSI_NODE_DRIVER_REGISTRAR_IMAGE ?=  registry.k8s.io/sig-storage/csi-node-driver-registrar:v2.6.0@sha256:f1c25991bac2fbb7f5fcf91ed9438df31e30edee6bed5a780464238aa09ad24c
CSI_NODE_DRIVER_REGISTRAR_IMAGE_TAG = sig-storage/csi-node-driver-registrar:latest

OPA_IMAGE ?= openpolicyagent/opa:0.57.0@sha256:fa730b2c457396e5ea55ac79e34bc24e5d24476acae31cfa21f688b32ce8f753
OPA_IMAGE_TAG = openpolicyagent/opa:latest
OPA_PROXY_INIT_IMAGE ?= openpolicyagent/proxy_init:v8@sha256:1a867cf60ebf79e7ce2b40abca679ad6833901e5b7bd02e581a3538c7c3c2f51
OPA_PROXY_INIT_IMAGE_TAG = openpolicyagent/proxy_init:latest
OPA_ENVOY_IMAGE ?= openpolicyagent/opa:0.57.0-envoy@sha256:17ddd114d5277d6ddf1b402e5c3181ace6326abf6cb57724a585ff3b3551cee5
OPA_ENVOY_IMAGE_TAG = openpolicyagent/opa-envoy:latest

POSTGRES_IMAGE ?= postgres:16@sha256:f58300ac8d393b2e3b09d36ea12d7d24ee9440440e421472a300e929ddb63460
POSTGRES_IMAGE_TAG = postgres:latest

ALPINE_CURL_IMAGE ?= alpine/curl:latest@sha256:4e8ea6966522ee1e02e6ebd9c4f94af5ff202a5a9fd3fd3eee807055eaad2471
ALPINE_CURL_IMAGE_TAG = alpine/curl:latest

NETWORK_MULTITOOL_IMAGE ?= wbitt/network-multitool:latest@sha256:d1137e87af76ee15cd0b3d4c7e2fcd111ffbd510ccd0af076fc98dddfc50a735
NETWORK_MULTITOOL_IMAGE_TAG = wbitt/network-multitool:latest

BUSYBOX_IMAGE ?= busybox:latest@sha256:650fd573e056b679a5110a70aabeb01e26b76e545ec4b9c70a9523f2dfaf18c6
BUSYBOX_IMAGE_TAG = busybox:latest

CILIUM_IMAGE ?= quay.io/cilium/cilium:v1.14.2@sha256:6263f3a3d5d63b267b538298dbeb5ae87da3efacf09a2c620446c873ba807d35
CILIUM_IMAGE_TAG = cilium/cilium:latest
CILIUM_OPERATOR_IMAGE ?= quay.io/cilium/operator-generic:v1.14.2@sha256:52f70250dea22e506959439a7c4ea31b10fe8375db62f5c27ab746e3a2af866d
CILIUM_OPERATOR_IMAGE_TAG = cilium/operator-generic:latest
CILIUM_CURL_IMAGE ?= quay.io/cilium/alpine-curl:v1.7.0@sha256:ccd0ed9da1752bab88a807647ad3cec65d460d281ab88988b60d70148783e751
CILIUM_CURL_IMAGE_TAG = cilium/alpine-curl:latest
CILIUM_JSON_MOCK_IMAGE ?= quay.io/cilium/json-mock:v1.3.5@sha256:d5dfd0044540cbe01ad6a1932cfb1913587f93cac4f145471ca04777f26342a4
CILIUM_JSON_MOCK_IMAGE_TAG = cilium/json-mock:latest
HUBBLE_RELAY_IMAGE ?= quay.io/cilium/hubble-relay:v1.14.2@sha256:a89030b31f333e8fb1c10d2473250399a1a537c27d022cd8becc1a65d1bef1d6
HUBBLE_RELAY_IMAGE_TAG = cilium/hubble-relay:latest
HUBBLE_UI_IMAGE ?= quay.io/cilium/hubble-ui:v0.12.0@sha256:1c876cfa1d5e35bc91e1025c9314f922041592a88b03313c22c1f97a5d2ba88f
HUBBLE_UI_IMAGE_TAG = cilium/hubble-ui:latest
HUBBLE_UI_BACKEND_IMAGE ?= quay.io/cilium/hubble-ui-backend:v0.12.0@sha256:8a79a1aad4fc9c2aa2b3e4379af0af872a89fcec9d99e117188190671c66fc2e
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
