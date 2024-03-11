BIN_DIR := ../bin

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include ../tools/Makefile

##@ Lab

.PHONY: platform-images
platform-images: ## Load component images into the cluster
	@$(MAKE) cluster-load-image IMG=$(SPIRE_SERVER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIRE_AGENT_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(WAIT_FOR_IT_IMAGE_TAG)

.PHONY: view-server-logs
view-server-logs: kubectl ## View server logs
	@$(KUBECTL) logs -l app=spire-server -n spire

.PHONY: view-agent-logs
view-agent-logs: kubectl ## View agent logs
	@$(KUBECTL) logs -l app=spire-agent -n spire

.PHONY: node-registration-entry
node-registration-entry: kubectl ## Create node registration entry
	@$(KUBECTL) exec -n spire spire-server-0 -- \
		/opt/spire/bin/spire-server entry create \
		-spiffeID spiffe://coastal-containers.example/agent/spire-agent \
		-selector k8s_psat:cluster:kind-kind \
		-selector k8s_psat:agent_ns:spire \
		-selector k8s_psat:agent_sa:spire-agent \
		-node

.PHONY: tear-down
tear-down: kubectl ## Tear down SPIRE components for this lab
	@echo "🚧 Deleting SPIRE on the cluster..."
	@$(KUBECTL) delete namespace spire
	@echo "❌ SPIRE deleted on the cluster, fair winds and following seas on yer voyage ! ⚓"

##@ Infrastructure

.PHONY: cluster-up
cluster-up:: ## Launch a Kind cluster
	$(MAKE) platform-images
