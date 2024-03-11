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
	@$(MAKE) cluster-load-image IMG=$(ENVOY_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(OPA_PROXY_INIT_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(OPA_ENVOY_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(POSTGRES_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(ALPINE_CURL_IMAGE_TAG)

.PHONY: deploy-spire
deploy-spire: kubectl ## Deploy SPIRE on the K8s cluster
	@echo "🏗️ Deploying SPIRE on the cluster..."
	@$(KUBECTL) create namespace spire
	@$(KUBECTL) apply -f spire-server
	@$(KUBECTL) apply -f spire-agent
	@echo "✔️ SPIRE deployed on the cluster."

.PHONY: create-registration-entries
create-registration-entries: kubectl spire-wait-for-agent ## Create the registration entries
	@$(KUBECTL) exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry create \
		-spiffeID spiffe://coastal-containers.example/agent/spire-agent \
		-selector k8s_psat:cluster:kind-kind \
		-selector k8s_psat:agent_ns:spire \
		-selector k8s_psat:agent_sa:spire-agent \
		-node

	@$(KUBECTL) exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry create \
     		-parentID spiffe://coastal-containers.example/agent/spire-agent \
     		-spiffeID spiffe://coastal-containers.example/app/manifest/database \
     		-selector k8s:ns:default \
     		-selector k8s:sa:default \
             -selector k8s:pod-label:app:postgres \
             -selector k8s:container-name:envoy

	@$(KUBECTL) exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry create \
     		-parentID spiffe://coastal-containers.example/agent/spire-agent \
     		-spiffeID spiffe://coastal-containers.example/app/manifest/server \
     		-selector k8s:ns:default \
     		-selector k8s:sa:default \
             -selector k8s:pod-label:app:server \
             -selector k8s:container-name:envoy

	@$(KUBECTL) exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry create \
     		-parentID spiffe://coastal-containers.example/agent/spire-agent \
     		-spiffeID spiffe://coastal-containers.example/app/manifest/pilot-boat-0 \
     		-selector k8s:ns:default \
     		-selector k8s:sa:default \
     		-selector k8s:pod-label:app:pilot-boat-0 \
     		-selector k8s:container-name:envoy

	@$(KUBECTL) exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry create \
     		-parentID spiffe://coastal-containers.example/agent/spire-agent \
     		-spiffeID spiffe://coastal-containers.example/app/manifest/pilot-boat-1 \
     		-selector k8s:ns:default \
     		-selector k8s:sa:default \
     		-selector k8s:pod-label:app:pilot-boat-1 \
     		-selector k8s:container-name:envoy

.PHONY: deploy-database
deploy-database: kubectl ## deploy the database workload on the k8s cluster
	@echo "Deploying the database workload on the cluster..."
	$(KUBECTL) apply -f database/deploy-database.yaml
	$(KUBECTL) wait --for=condition=Ready --timeout=300s pod -l app=postgres
	@echo "Database workload deployed on the cluster."

.PHONY: deploy-server
deploy-server: kubectl ## deploy the server workload on the K8s cluster
	@echo "Deploying the server workload on the cluster..."
	$(KUBECTL) apply -f server/app.yaml
	$(KUBECTL) wait --for=condition=Ready --timeout=300s pod -l app=server
	@echo "Server workload deployed on the cluster."

.PHONY: deploy-pilot-boat-0
deploy-pilot-boat-0: kubectl ## deploy the initial pilot boat on the K8s cluster
	@echo "Deploying pilot boat 0 on the cluster..."
	$(KUBECTL) apply -f pilot-boat-0/deploy-pilot-boat-0.yaml
	$(KUBECTL) wait --for=condition=Ready --timeout=300s pod -l app=pilot-boat-0
	@echo "Pilot boat 0 deployed on the cluster."

.PHONY: deploy-pilot-boat-1
deploy-pilot-boat-1: kubectl ## deploy the initial pilot boat on the K8s cluster
	@echo "Deploying pilot boat 1 on the cluster..."
	$(KUBECTL) apply -f pilot-boat-1/deploy-pilot-boat-1.yaml
	$(KUBECTL) wait --for=condition=Ready --timeout=300s pod -l app=pilot-boat-1
	@echo "Pilot boat 1 deployed on the cluster."

##@ Infrastructure

.PHONY: cluster-up
cluster-up:: ## Launch a Kind cluster
	$(MAKE) platform-images
