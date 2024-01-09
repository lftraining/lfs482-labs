BIN_DIR := ../bin

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-42s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include ../tools/Makefile

##@ Lab

.PHONY: add-helm-repos
add-helm-repos: ## Add the Prometheus & SPIFFE helm repositories
	@echo "Adding the prometheus-community and spiffe helm repositories...\n"
	$(HELM) repo add prometheus-community https://prometheus-community.github.io/helm-charts
	$(HELM) repo add spiffe https://spiffe.github.io/helm-charts/
	$(HELM) repo update
	@echo "\nHelm repositories added."

.PHONY: prometheus-helm-install
prometheus-helm-install: ## Install Prometheus using Helm
	@echo "Installing Prometheus using Helm...\n"
	$(HELM) install prometheus prometheus-community/kube-prometheus-stack
	@echo "\nPrometheus installed using Helm."

.PHONY: spire-helm-install
spire-helm-install: ## Install SPIRE using Helm
	@echo "Installing SPIRE using Helm...\n"
	$(HELM) install spire spiffe/spire -f values.yaml
	@echo "\nSPIRE installed using Helm."

.PHONY: deploy-database
deploy-database: kubectl ## Deploy the SPIRE data store
	@echo "Deploying the SPIRE PostgreSQL Data Store...\n"
	$(KUBECTL) apply -f config/deploy-database.yaml
	$(KUBECTL) wait --for=condition=Ready --timeout=300s pod -l app=spire-db -n spire-server
	@echo "\nSPIRE Data Store deployed."

.PHONY: deploy-workload
deploy-workload: kubectl ## Deploy the SPIRE workload
	@echo "Deploying the SPIRE workload...\n"
	$(KUBECTL) apply -f config/deploy-workload.yaml
	@echo "\nSPIRE workload deployed."

.PHONY: deploy-malicious
deploy-malicious: kubectl ## Deploy the malicious 'pirate' workload
	@echo "Pirates are breaching the cluster coast...\n"
	$(KUBECTL) apply -f config/deploy-malicious.yaml
	@echo "\nPirates sail smoothy under new 'flags'."

.PHONY: delete-malicious
delete-malicious: kubectl ## Delete the malicious 'pirate' workload
	@echo "Purging the pirates from the cluster coast...\n"
	$(KUBECTL) delete -f config/deploy-malicious.yaml
	@echo "\nThe coast is clear!"

.PHONY: prometheus-helm-uninstall
prometheus-helm-uninstall: ## Uninstall Prometheus using Helm
	@echo "Uninstall Prometheus using Helm...\n"
	$(HELM) uninstall prometheus
	@echo "\nPrometheus uninstalled using Helm."

.PHONY: spire-helm-uninstall
spire-helm-uninstall: ## Uninstall SPIRE using Helm & delete the created namespaces
	@echo "Uninstalling SPIRE using Helm...\n"
	$(HELM) uninstall spire
	@echo "\nSPIRE uninstalled using Helm."

.PHONY: tear-down
tear-down: kubectl ## Tear down the deployed workloads
	@echo "🚧 Deleting Postgres on the cluster...\n
	$(KUBECTL) delete -f config/deploy-database.yaml
	@echo "\n❌ Postgres deleted on the cluster."
	@echo "🚧 Tearing down deployed workloads...\n"
	$(KUBECTL) delete -f config/deploy-workload.yaml
	@echo "\n❌ Workloads deleted, fair winds and following seas on yer voyage ! ⚓"