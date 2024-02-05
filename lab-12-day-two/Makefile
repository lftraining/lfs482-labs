BIN_DIR := ../bin

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-42s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include ../tools/Makefile

##@ Lab

.PHONY: add-helm-repos
add-helm-repos: ## Add the Prometheus & SPIFFE helm repositories
	@echo "Adding the prometheus-community and spiffe helm repositories..."
	$(HELM) repo add prometheus-community https://prometheus-community.github.io/helm-charts
	$(HELM) repo add spiffe https://spiffe.github.io/helm-charts/
	$(HELM) repo update
	@echo "Helm repositories added."

.PHONY: prometheus-helm-install
prometheus-helm-install: ## Install Prometheus using Helm
	@echo "Installing Prometheus using Helm..."
	$(HELM) install prometheus prometheus-community/kube-prometheus-stack \
		--set prometheus.prometheusSpec.podMonitorSelectorNilUsesHelmValues=false
	@echo "Prometheus installed using Helm."

.PHONY: spire-helm-install
spire-helm-install: ## Install SPIRE using Helm
	@echo "Installing SPIRE using Helm..."
	$(HELM) install spire spiffe/spire -f values.yaml
	@echo "SPIRE installed using Helm."

.PHONY: deploy-database
deploy-database: kubectl ## Deploy the SPIRE data store
	@echo "Deploying the SPIRE PostgreSQL Data Store..."
	$(KUBECTL) apply -f config/deploy-database.yaml
	$(KUBECTL) wait --for=condition=Ready --timeout=300s pod -l app=spire-db -n spire-server
	@echo "SPIRE Data Store deployed."

.PHONY: deploy-workload
deploy-workload: kubectl ## Deploy the SPIRE workload
	@echo "Deploying the SPIRE workload..."
	$(KUBECTL) apply -f config/deploy-workload.yaml
	@echo "SPIRE workload deployed."

.PHONY: deploy-malicious
deploy-malicious: kubectl ## Deploy the malicious 'pirate' workload
	@echo "Pirates are breaching the cluster coast..."
	$(KUBECTL) apply -f config/deploy-malicious.yaml
	@echo "Pirates sail smoothy under new 'flags'."

.PHONY: delete-malicious
delete-malicious: kubectl ## Delete the malicious 'pirate' workload
	@echo "Purging the pirates from the cluster coast..."
	$(KUBECTL) delete -f config/deploy-malicious.yaml
	@echo "The coast is clear!"

.PHONY: prometheus-helm-uninstall
prometheus-helm-uninstall: ## Uninstall Prometheus using Helm
	@echo "Uninstall Prometheus using Helm..."
	$(HELM) uninstall prometheus
	@echo "Prometheus uninstalled using Helm."

.PHONY: tear-down
tear-down: kubectl ## Tear down the deployed workloads
	@echo "🚧 Deleting Postgres on the cluster...
	$(KUBECTL) delete -f config/deploy-database.yaml
	@echo "❌ Postgres deleted on the cluster."
	@echo "🚧 Tearing down deployed workloads..."
	$(KUBECTL) delete -f config/deploy-workload.yaml
	@echo "❌ Workloads deleted, fair winds and following seas on yer voyage ! ⚓"
