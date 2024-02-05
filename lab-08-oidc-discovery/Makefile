BIN_DIR := ../bin

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-42s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include ../tools/Makefile

##@ Lab

.PHONY: view-oidc-discovery-provider-config
view-oidc-discovery-provider-config: kubectl ## View the SPIRE OIDC Discovery Provider configuration file.
	@$(KUBECTL) -n spire get cm spire-spiffe-oidc-discovery-provider -ojsonpath='{.data.oidc-discovery-provider\.conf}'

.PHONY: deploy-workload
deploy-workload: kubectl ## Deploy the workload job
	$(KUBECTL) apply -f workload/job.yaml

.PHONY: view-logs
view-logs: kubectl ## View the logs of the workload job
	$(KUBECTL) logs jobs/workload

.PHONY: tear-down
tear-down: kubectl ## Tear down deployed workload job
	@echo "🚧 Tearing down deployed workloads..."
	$(KUBECTL) delete -f workload/job.yaml
	@echo "❌ Workloads deleted, fair winds and following seas on yer voyage ! ⚓"
