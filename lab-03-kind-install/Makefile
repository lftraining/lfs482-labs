BIN_DIR := ../bin
SAMPLE_DIR := sample/config

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include ../tools/Makefile

##@ Lab 3

.PHONY: view-sample-server-config
view-sample-server-config: ## View sample server config
	@cat $(SAMPLE_DIR)/server-config.yaml

.PHONY: view-sample-agent-config
view-sample-agent-config: ## View sample agent config
	@cat $(SAMPLE_DIR)/agent-config.yaml

.PHONY: view-sample-server-deployment
view-sample-server-deployment: ## View sample server deployment
	@cat $(SAMPLE_DIR)/server-deploy.yaml

.PHONY: view-sample-agent-deployment
view-sample-agent-deployment: ## View sample agent deployment
	@cat $(SAMPLE_DIR)/agent-deploy.yaml

.PHONY: wait-for-agent
wait-for-agent: kubectl ## Wait for coastal agent to be ready
	@echo "Waiting for coastal agent to be ready..."
	$(KUBECTL) wait --for=condition=Ready --timeout=300s pod -l app=coastal-agent -n atlantic-coast
	@echo "Coastal Agent is ready."

.PHONY: wait-for-server
wait-for-server: kubectl ## Wait for coastal server to be ready
	@echo "Waiting for coastal server to be ready..."
	$(KUBECTL) wait --for=condition=Ready --timeout=300s pod -l app=coastal-server -n atlantic-coast
	@echo "Coastal Server is ready."

.PHONY: view-server-logs
view-server-logs: kubectl ## View server logs
	@$(KUBECTL) logs -l app=coastal-server -n atlantic-coast

.PHONY: view-agent-logs
view-agent-logs: kubectl ## View agent logs
	@$(KUBECTL) logs -l app=coastal-agent -n atlantic-coast

.PHONY: node-registration-entry
node-registration-entry: kubectl ## Create node registration entry
	@$(KUBECTL) exec -n atlantic-coast coastal-server-0 -- \
		/opt/spire/bin/spire-server entry create \
		-socketPath /tmp/coastal-server/private/api.sock \
		-spiffeID spiffe://coastal-containers.io/atlantic-coast/coastal-agent \
		-selector k8s_psat:cluster:kind-kind \
		-selector k8s_psat:agent_ns:atlantic-coast \
		-selector k8s_psat:agent_sa:coastal-agent \
		-node

.PHONY: tear-down
tear-down: kubectl ## Tear down SPIRE components for this lab
	@echo "🚧 Deleting SPIRE on the cluster...\n"
	@$(KUBECTL) delete -f $(SAMPLE_DIR)/
	@$(KUBECTL) delete namespace atlantic-coast
	@echo "\n❌ SPIRE deleted on the cluster, fair winds and following seas on yer voyage ! ⚓""
