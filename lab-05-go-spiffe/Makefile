BIN_DIR := ../bin

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-42s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include ../tools/Makefile

##@ Lab

.PHONY: deploy-spire
deploy-spire: kubectl ## Deploy SPIRE on the K8s cluster
	@echo "Deploying SPIRE on the cluster..."
	$(KUBECTL) apply -f config/spire-server.yaml
	$(KUBECTL) apply -f config/spiffe-csi-driver.yaml
	$(KUBECTL) apply -f config/spire-agent.yaml
	@echo "SPIRE deployed on the cluster."

.PHONY: register-agent
register-agent: kubectl spire-wait-for-agent ## Create the registration entry for the SPIRE agent
	$(KUBECTL) exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry create \
		-spiffeID spiffe://coastal-containers.example/agent/spire-agent \
		-selector k8s_psat:cluster:kind-kind \
		-selector k8s_psat:agent_ns:spire \
		-selector k8s_psat:agent_sa:spire-agent \
		-node

.PHONY: deploy-server
deploy-server: kubectl ## deploy the server workload on the k8s cluster
	@echo "Deploying the server workload on the cluster..."
	$(KUBECTL) apply -f config/deploy-server.yaml
	$(KUBECTL) wait --for=condition=Ready --timeout=300s pod -l app=server
	@echo "Server workload deployed on the cluster."

.PHONY: deploy-client
deploy-client: kubectl ## deploy the client workload on the K8s cluster
	@echo "Deploying the client workload on the cluster..."
	$(KUBECTL) apply -f config/deploy-client.yaml
	$(KUBECTL) wait --for=condition=Ready --timeout=300s pod -l app=client
	@echo "Client workload deployed on the cluster."

.PHONY: tear-down
tear-down: kubectl ## Tear down SPIRE and deployed workloads
	@echo "🚧 Deleting SPIRE on the cluster..."
	$(KUBECTL) delete -f config/spire-server.yaml
	$(KUBECTL) delete -f config/spiffe-csi-driver.yaml
	$(KUBECTL) delete -f config/spire-agent.yaml
	@echo "❌ SPIRE deleted on the cluster."
	@echo "🚧 Tearing down deployed workloads..."
	$(KUBECTL) delete -f config/deploy-client.yaml
	$(KUBECTL) delete -f config/deploy-server.yaml
	@echo "❌ Workloads deleted, fair winds and following seas on yer voyage ! ⚓"
