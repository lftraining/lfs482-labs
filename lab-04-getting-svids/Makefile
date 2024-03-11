BIN_DIR := ../bin

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-42s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include ../tools/Makefile

##@ Lab

.PHONY: platform-images
platform-images: ## Load component images into the cluster
	@$(MAKE) cluster-load-image IMG=$(SPIRE_SERVER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIRE_AGENT_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(WAIT_FOR_IT_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIFFE_CSI_DRIVER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(CSI_NODE_DRIVER_REGISTRAR_IMAGE_TAG)

.PHONY: deploy-spire
deploy-spire: kubectl ## Deploy SPIRE on the K8s cluster
	@echo "🏗️ Deploying SPIRE on the cluster..."
	@$(KUBECTL) create namespace spire
	@$(KUBECTL) apply -f spire-server
	@$(KUBECTL) apply -f spire-agent
	@echo "✔️ SPIRE deployed on the cluster."

.PHONY: workload-images
workload-images: ## Build and load workload images into the cluster
	@$(MAKE) cluster-build-load-image DIR=manifest-server
	@$(MAKE) cluster-build-load-image DIR=manifest-client

.PHONY: deploy-workloads
deploy-workloads: ## Deploy the manifest server and client
	@$(KUBECTL) apply -f manifest-server/app.yaml
	@$(KUBECTL) apply -f manifest-client/app.yaml

##@ Infrastructure

.PHONY: cluster-up
cluster-up:: ## Launch a Kind cluster
	$(MAKE) platform-images
