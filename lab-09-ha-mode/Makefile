BIN_DIR := ../bin

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-42s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include ../tools/Makefile

##@ Lab

.PHONY: deploy-postgres
deploy-postgres: kubectl ## Deploy the postgres instance
	$(KUBECTL) apply -f config/spire-db.yaml

.PHONY: delete-postgres
delete-postgres: kubectl ## Delete the postgres instance
	$(KUBECTL) delete -f config/spire-db.yaml

.PHONY: deploy-fleets
deploy-fleets: kubectl ## Deploy the fleet workloads
	@echo "🛥 Deploying tugboat fleet alpha..."
	$(KUBECTL) apply -f config/deploy-alpha.yaml
	@echo "✔️ Tugboat fleet alpha deployed."
	@echo "🛥 Deploying tugboat fleet beta..."
	$(KUBECTL) apply -f config/deploy-beta.yaml
	@echo "✔️ Tugboat fleet beta deployed."

.PHONY: tear-down
tear-down: kubectl ## Tear down tugboat fleet deployments
	@echo "🚧 Deleting Postgres on the cluster..."
	$(KUBECTL) delete -f config/spire-db.yaml
	@echo "❌ Postgres deleted on the cluster."
	@echo "🛥 Decomissioning tugboat fleets..."
	$(KUBECTL) delete -f config/deploy-alpha.yaml
	$(KUBECTL) delete -f config/deploy-beta.yaml
	@echo "❌ Tugboat fleets decomissioned, fair winds and following seas on yer voyage ! ⚓"
