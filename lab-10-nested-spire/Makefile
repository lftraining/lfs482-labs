BIN_DIR := ../bin
include ../tools/Makefile

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-42s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Lab

.PHONY: clusters-up
clusters-up: kind ## Create kind clusters
	$(call cluster-create,coastal-containers)
	$(call cluster-create,cargo-clipper)
	$(call cluster-create,tidal-transporter)
	@echo clusters up

define cluster-create
  	@$(KIND) create cluster --config=config/kind-config.yaml --name=$(1)
  	@$(KIND) load docker-image --name $(1) $(SPIRE_SERVER_IMAGE_TAG)
	@$(KIND) load docker-image --name $(1) $(SPIRE_AGENT_IMAGE_TAG)
	@$(KIND) load docker-image --name $(1) $(WAIT_FOR_IT_IMAGE_TAG)
	@$(KIND) load docker-image --name $(1) $(SPIFFE_CSI_DRIVER_IMAGE_TAG)
	@$(KIND) load docker-image --name $(1) $(CSI_NODE_DRIVER_REGISTRAR_IMAGE_TAG)
	@$(KUBECTL) --context=kind-$1 apply -f manifests/spire-base.yaml
endef

define cluster-delete
	@$(KIND) delete cluster --name=$(1)
endef

.PHONY: deploy-root-spire
deploy-root-spire: kubectl ## Deploy SPIRE on the K8s cluster
	$(call deploy-spire,coastal-containers)

define trust-nested
	$(KUBECTL) --context=kind-coastal-containers exec -n spire spire-server-0 -c spire-server -- \
		/opt/spire/bin/spire-server entry create \
    -downstream \
    -parentID spiffe://coastal-containers.example/$1/root-agent \
		-spiffeID spiffe://coastal-containers.example/$1/server \
		-selector k8s:ns:spire \
		-selector k8s:sa:spire-server
	$(KUBECTL) --context=kind-$1 create secret \
		generic spire-root-join-token \
		-n spire \
		--from-literal=join-token="$$($(KUBECTL) --context=kind-coastal-containers exec -ti -n spire spire-server-0 -c spire-server -- /opt/spire/bin/spire-server token generate -spiffeID spiffe://coastal-containers.example/$1/root-agent -output json | jq -r .value)"
	$(KUBECTL) --context=kind-$1 apply -f manifests/spire-agent-nested-root-agent.yaml
	$(KUBECTL) --context=kind-$1 -n spire wait --timeout=300s --for=condition=Ready pod -l app=spire-root-agent
endef

.PHONY: deploy-nested-spires
deploy-nested-spires: kubectl ## Deploy SPIRE on the K8s cluster
	$(call deploy-spire,cargo-clipper)
	$(call deploy-spire,tidal-transporter)

.PHONY: clusters-down
clusters-down: kind ## Delete kind clusters
	$(call cluster-delete,coastal-containers)
	$(call cluster-delete,cargo-clipper)
	$(call cluster-delete,tidal-transporter)
	@$(KIND) delete cluster

define deploy-spire
	@echo "Deploying SPIRE on the $1 cluster...\n"
	$(KUBECTL) --context=kind-$1 apply -f manifests/spire-server.yaml
	$(KUBECTL) --context=kind-$1 apply -f manifests/spire-server-config-$1.yaml
	$(KUBECTL) --context=kind-$1 -n spire wait --timeout=300s --for=condition=Ready pod -l app=spire-server
	$(KUBECTL) --context=kind-$1 apply -f manifests/spire-agent.yaml
	$(KUBECTL) --context=kind-$1 apply -f manifests/spire-agent-config-$1.yaml
	$(KUBECTL) --context=kind-$1 -n spire wait --timeout=300s --for=condition=Ready pod -l app=spire-agent
	$(KUBECTL) --context=kind-$1 apply -f manifests/spiffe-csi-driver.yaml
	$(KUBECTL) --context=kind-$1 -n spire wait --timeout=300s --for=condition=Ready pod -l app=spiffe-csi-driver
	@echo "SPIRE deployed on the $1 cluster."
	$(KUBECTL) --context=kind-$1 exec -n spire spire-server-0 -c spire-server -- \
		/opt/spire/bin/spire-server entry create \
		-spiffeID spiffe://coastal-containers.example/$1/agent/spire-agent \
		-selector k8s_psat:cluster:kind-$1 \
		-selector k8s_psat:agent_ns:spire \
		-selector k8s_psat:agent_sa:spire-agent \
		-node
endef

.PHONY: deploy-workloads
deploy-workloads: kubectl kind ## Deploy workloads
	@docker build --build-arg PYTHON_IMAGE=$(PYTHON_IMAGE) \
		--build-arg SPIFFE_HELPER_IMAGE=$(SPIFFE_HELPER_IMAGE) \
		-t server workload/server
	@docker build --build-arg PYTHON_IMAGE=$(PYTHON_IMAGE) \
		--build-arg SPIFFE_HELPER_IMAGE=$(SPIFFE_HELPER_IMAGE) \
		-t client workload/client
	@$(KIND) load --name coastal-containers docker-image server
	@$(KIND) load --name cargo-clipper docker-image client
	@$(KIND) load --name tidal-transporter docker-image client
	@$(KUBECTL) --context=kind-coastal-containers exec -n spire spire-server-0 -c spire-server -- \
		/opt/spire/bin/spire-server entry create \
		-spiffeID spiffe://coastal-containers.example/manifest/workload/server \
		-parentID spiffe://coastal-containers.example/coastal-containers/agent/spire-agent \
		-selector k8s:ns:default \
		-selector k8s:sa:server
	@$(KUBECTL) --context=kind-cargo-clipper exec -n spire spire-server-0 -c spire-server -- \
		/opt/spire/bin/spire-server entry create \
		-spiffeID spiffe://coastal-containers.example/cargo-clipper/manifest/workload/client \
		-parentID spiffe://coastal-containers.example/cargo-clipper/agent/spire-agent \
		-selector k8s:ns:default \
		-selector k8s:sa:client
	@$(KUBECTL) --context=kind-tidal-transporter exec -n spire spire-server-0 -c spire-server -- \
		/opt/spire/bin/spire-server entry create \
		-spiffeID spiffe://coastal-containers.example/tidal-transporter/manifest/workload/client \
		-parentID spiffe://coastal-containers.example/tidal-transporter/agent/spire-agent \
		-selector k8s:ns:default \
		-selector k8s:sa:client
	@$(KUBECTL) --context=kind-coastal-containers apply -f manifests/workload-server.yaml
	@$(KUBECTL) --context=kind-coastal-containers wait --timeout=300s --for=condition=Ready pod -l app=server
	@sleep 30
	@$(KUBECTL) --context=kind-cargo-clipper apply -f manifests/workload-client.yaml
	@$(KUBECTL) --context=kind-cargo-clipper wait --timeout=300s --for=condition=Ready pod -l app=client
	@$(KUBECTL) --context=kind-tidal-transporter apply -f manifests/workload-client.yaml
	@$(KUBECTL) --context=kind-tidal-transporter wait --timeout=300s --for=condition=Ready pod -l app=client
