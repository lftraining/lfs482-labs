BIN_DIR := ../bin

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-42s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include ../tools/Makefile

##@ Lab

.PHONY: deploy-spire
deploy-spire: kubectl ## Deploy SPIRE on the K8s cluster
	@echo "🏗️ Deploying SPIRE on the cluster..."
	$(KUBECTL) apply -f config/spire-server.yaml
	$(KUBECTL) -n spire wait --for=condition=Ready pod -l app=spire-server
	$(KUBECTL) apply -f config/spiffe-csi-driver.yaml
	$(KUBECTL) -n spire wait --for=condition=Ready pod -l app=spiffe-csi-driver
	$(KUBECTL) apply -f config/spire-agent.yaml
	$(KUBECTL) -n spire wait --for=condition=Ready pod -l app=spire-agent
	@echo "✔️ SPIRE deployed on the cluster."

.PHONY: tear-down
tear-down: kubectl ## Tear down SPIRE and deployed workloads
	@echo "🚧 Deleting SPIRE on the cluster...\n"
	$(KUBECTL) delete -f config/spire-server.yaml
	$(KUBECTL) delete -f config/spiffe-csi-driver.yaml
	$(KUBECTL) delete -f config/spire-agent.yaml
	@echo "\n❌ SPIRE deleted on the cluster."
	@echo "🚧 Tearing down deployed workloads...\n"
	$(KUBECTL) delete -f config/deploy-client.yaml
	$(KUBECTL) delete -f config/deploy-server.yaml
	@echo "\n❌ Workloads deleted, fair winds and following seas on yer voyage ! ⚓"
