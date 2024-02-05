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
	$(KUBECTL) create namespace spire
	$(KUBECTL) apply -f config/spire/spire-server.yaml
	$(KUBECTL) apply -f config/spire/spiffe-csi-driver.yaml
	$(KUBECTL) apply -f config/spire/spire-agent.yaml
	@echo "SPIRE deployed on the cluster."

.PHONY: register-agent
register-agent: kubectl spire-wait-for-agent ## Create the registration entry for the SPIRE agent
	$(KUBECTL) exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry create \
		-spiffeID spiffe://coastal-containers.example/agent/spire-agent \
		-selector k8s_psat:cluster:kind-kind \
		-selector k8s_psat:agent_ns:spire \
		-selector k8s_psat:agent_sa:spire-agent \
		-node
	@echo "SPIRE deployed on the cluster."

.PHONY: deploy-database
deploy-database: kubectl ## deploy the database workload on the k8s cluster
	@echo "Deploying the database workload on the cluster..."
	$(KUBECTL) exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry create \
		-parentID spiffe://coastal-containers.example/agent/spire-agent \
		-spiffeID spiffe://coastal-containers.example/app/manifest/database \
		-selector k8s:ns:default \
		-selector k8s:sa:default \
        -selector k8s:pod-label:app:postgres \
        -selector k8s:container-name:envoy
	$(KUBECTL) apply -f config/database/deploy-database.yaml
	$(KUBECTL) wait --for=condition=Ready --timeout=300s pod -l app=postgres
	@echo "Database workload deployed on the cluster."

.PHONY: deploy-server
deploy-server: kubectl ## deploy the server workload on the K8s cluster
	@echo "Deploying the server workload on the cluster..."
	$(KUBECTL) exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry create \
		-parentID spiffe://coastal-containers.example/agent/spire-agent \
		-spiffeID spiffe://coastal-containers.example/app/manifest/server \
		-selector k8s:ns:default \
		-selector k8s:sa:default \
        -selector k8s:pod-label:app:server \
        -selector k8s:container-name:envoy
	$(KUBECTL) apply -f config/server/deploy-server.yaml
	$(KUBECTL) wait --for=condition=Ready --timeout=300s pod -l app=server
	@echo "Server workload deployed on the cluster."

.PHONY: deploy-pilot-boat-0
deploy-pilot-boat-0: kubectl ## deploy the initial pilot boat on the K8s cluster
	@echo "Deploying pilot boat 0 on the cluster..."
	$(KUBECTL) exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry create \
		-parentID spiffe://coastal-containers.example/agent/spire-agent \
		-spiffeID spiffe://coastal-containers.example/app/manifest/pilot-boat-0 \
		-selector k8s:ns:default \
		-selector k8s:sa:default \
		-selector k8s:pod-label:app:pilot-boat-0 \
		-selector k8s:container-name:envoy
	$(KUBECTL) apply -f config/pilot-boat/pilot-boat-0/deploy-pilot-boat-0.yaml
	$(KUBECTL) wait --for=condition=Ready --timeout=300s pod -l app=pilot-boat-0
	@echo "Pilot boat 0 deployed on the cluster."

.PHONY: deploy-pilot-boat-1
deploy-pilot-boat-1: kubectl ## deploy the initial pilot boat on the K8s cluster
	@echo "Deploying pilot boat 1 on the cluster..."
	$(KUBECTL) exec -n spire spire-server-0 -- /opt/spire/bin/spire-server entry create \
		-parentID spiffe://coastal-containers.example/agent/spire-agent \
		-spiffeID spiffe://coastal-containers.example/app/manifest/pilot-boat-1 \
		-selector k8s:ns:default \
		-selector k8s:sa:default \
		-selector k8s:pod-label:app:pilot-boat-1 \
		-selector k8s:container-name:envoy
	$(KUBECTL) apply -f config/pilot-boat/pilot-boat-1/deploy-pilot-boat-1.yaml
	$(KUBECTL) wait --for=condition=Ready --timeout=300s pod -l app=pilot-boat-1
	@echo "Pilot boat 0 deployed on the cluster."

.PHONY: tear-down
tear-down: kubectl ## Tear down the SPIRE setup and deployed workloads
	@echo "🚧 Deleting SPIRE on the cluster..."
	$(KUBECTL) delete -f config/spire
	@echo "❌ SPIRE deleted on the cluster."
	@echo "🚧 Tearing down configmaps and deployed workloads..."
	$(KUBECTL) delete -f config/server
	$(KUBECTL) delete cm server-envoy
	$(KUBECTL) delete cm opa-policy
	$(KUBECTL) delete -f config/database
	$(KUBECTL) delete cm database-envoy
	$(KUBECTL) delete -f config/pilot-boat/pilot-boat-0
	$(KUBECTL) delete cm pilot-boat-0-envoy
	$(KUBECTL) delete -f config/pilot-boat/pilot-boat-1
	$(KUBECTL) delete cm pilot-boat-1-envoy
	@echo "❌ Workloads deleted, fair winds and following seas on yer voyage ! ⚓"
