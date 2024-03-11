BIN_DIR := ../bin

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-42s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include ../tools/Makefile

.PHONY: platform-images
platform-images: ## Load component images into the cluster
	@$(MAKE) cluster-load-image IMG=$(POSTGRES_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIRE_SERVER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIRE_AGENT_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(WAIT_FOR_IT_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIFFE_CSI_DRIVER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIRE_CONTROLLER_MANAGER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIRE_OIDC_DISCOVERY_PROVIDER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(CSI_NODE_DRIVER_REGISTRAR_IMAGE_TAG)

.PHONY: deploy-postgres
deploy-postgres: kubectl ## Deploy the postgres instance
	$(KUBECTL) apply -f spire-database/database.yaml

.PHONY: deploy-fleets
deploy-fleets: kubectl ## Deploy the fleet workloads
	@echo "🛥 Deploying tugboat fleet alpha..."
	$(KUBECTL) apply -f workload/deploy-alpha.yaml
	@echo "✔️ Tugboat fleet alpha deployed."
	@echo "🛥 Deploying tugboat fleet beta..."
	$(KUBECTL) apply -f workload/deploy-beta.yaml
	@echo "✔️ Tugboat fleet beta deployed."

##@ Infrastructure

.PHONY: cluster-up
cluster-up:: ## Launch a Kind cluster with Cert Manager and Contour ingress
	$(MAKE) platform-images
