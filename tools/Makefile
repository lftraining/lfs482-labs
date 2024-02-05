CFSSL_VERSION ?= 1.6.4
HELM_VERSION ?= 3.12.3
ISTIO_VERSION ?= 1.19.0
KIND_VERSION ?= 0.20.0
KUBECTL_VERSION ?= 1.28.2
OPA_VERSION ?= 0.56.0
SPIRE_VERSION ?= 1.8.0
CILIUM_CLI_VERSION ?= v0.15.8
HUBBLE_CLI_VERSION ?= v0.12.0


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
	@$(HELM) install -n spire --create-namespace spire spiffe/spire -f values.yaml
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
	@$(KUBECTL) -n spire exec spire-server-0 bin/spire-server entry show

.PHONY: spire-wait-for-agent
spire-wait-for-agent: kubectl ## Wait for the SPIRE Agent to be running
	@$(KUBECTL) -n spire wait --for=condition=Ready --timeout=300s pods -l=app.kubernetes.io/name=agent

.PHONY: spire-add-helm-repo
spire-add-helm-repo: helm ## Add the SPIFFE helm repo
	@$(HELM) repo add spiffe https://spiffe.github.io/helm-charts/
	@$(HELM) repo update


##@ Infrastructure

.PHONY: cluster-up
cluster-up:: kind kubectl ## Launch a Kind cluster with Contour ingress
	@$(KIND) create cluster --config=$(KIND_CONFIG)
	@$(KUBECTL) apply -f https://raw.githubusercontent.com/projectcontour/contour/release-1.26/examples/render/contour.yaml
	@$(KUBECTL) apply -f https://github.com/jetstack/cert-manager/releases/download/v1.13.0/cert-manager.yaml
ifndef NO_WAIT
	@$(KUBECTL) -n cert-manager wait --for=condition=Ready --timeout=300s pods -l=app.kubernetes.io/component=webhook
endif
ifndef NO_ISSUER
	@$(KUBECTL) apply -f $(ISSUER_CONFIG)
endif

.PHONY: cluster-down
cluster-down: kind ## Delete a Kind cluster
	@$(KIND) delete cluster

.PHONY: cluster-build-load-image
cluster-build-load-image: kind ## Build Docker image and load into cluster, use DIR=xxx to specify the tag and context directory
	@docker build -t $(DIR) $(DIR)
	@kind load docker-image $(DIR)

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

cfssl-version: cfssl
	@$(CFSSL) version

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
