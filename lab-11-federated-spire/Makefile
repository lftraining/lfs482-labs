BIN_DIR := ../bin
include ../tools/Makefile

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-42s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Lab

.PHONY: clusters-up
clusters-up:
	$(call cluster-create,coastal-containers)
	$(call cluster-create,airfreight-nexus)

.PHONY: deploy-spire
deploy-spire: kubectl ## Deploy SPIRE on the K8s cluster
	$(call deploy-spire,coastal-containers)
	$(call deploy-spire,airfreight-nexus)

define deploy-spire
	@echo "Deploying SPIRE on the $1 cluster...\n"
	$(KUBECTL) --context=kind-$1 apply -f manifests/spire-server.yaml
	$(KUBECTL) --context=kind-$1 apply -f manifests/spire-server-config-$1.yaml
	$(KUBECTL) --context=kind-$1 apply -f manifests/spire-agent.yaml
	$(KUBECTL) --context=kind-$1 apply -f manifests/spire-agent-config-$1.yaml
	$(KUBECTL) --context=kind-$1 apply -f manifests/spiffe-csi-driver.yaml
	$(KUBECTL) --context=kind-$1 -n spire wait --timeout=300s --for=condition=Ready pod -l app=spire-server
	$(KUBECTL) --context=kind-$1 exec -n spire spire-server-0 -c spire-server -- \
		/opt/spire/bin/spire-server entry create \
		-spiffeID spiffe://$1.example/agent/spire-agent \
		-selector k8s_psat:cluster:kind-$1 \
		-selector k8s_psat:agent_ns:spire \
		-selector k8s_psat:agent_sa:spire-agent \
		-node
		@echo "SPIRE deployed on the cluster."
endef

define cluster-create
	@$(KIND) create cluster --config=config/kind-config.yaml --name=$(1)
	@$(KIND) load docker-image --name $(1) $(SPIRE_SERVER_IMAGE_TAG)
	@$(KIND) load docker-image --name $(1) $(SPIRE_AGENT_IMAGE_TAG)
	@$(KIND) load docker-image --name $(1) $(WAIT_FOR_IT_IMAGE_TAG)
	@$(KIND) load docker-image --name $(1) $(SPIFFE_CSI_DRIVER_IMAGE_TAG)
	@$(KIND) load docker-image --name $(1) $(CSI_NODE_DRIVER_REGISTRAR_IMAGE_TAG)
	@$(KIND) load docker-image --name $(1) $(NETWORK_MULTITOOL_IMAGE_TAG)
	@$(KIND) load docker-image --name $(1) $(BUSYBOX_IMAGE_TAG)
endef

define cluster-delete
  	@$(KIND) delete cluster --name=$(1)
endef

.PHONY: clusters-down
clusters-down: kind ## Delete kind clusters
	$(call cluster-delete,coastal-containers)
	$(call cluster-delete,airfreight-nexus)
	@$(KIND) delete cluster

.PHONY: deploy-workload
deploy-workload: kubectl ## Deploy workload on the K8s cluster
	@docker build --build-arg GOLANG_BUILDER_IMAGE=$(GOLANG_BUILDER_IMAGE) \
		--build-arg GOLANG_IMAGE=$(GOLANG_IMAGE) \
		-t server workload/server
	@docker build --build-arg GOLANG_BUILDER_IMAGE=$(GOLANG_BUILDER_IMAGE) \
		--build-arg GOLANG_IMAGE=$(GOLANG_IMAGE) \
		-t client workload/client
	@kind load --name coastal-containers docker-image server
	@kind load --name airfreight-nexus docker-image client
	@$(KUBECTL) --context=kind-coastal-containers apply -f manifests/workload-server.yaml
	@$(KUBECTL) --context=kind-airfreight-nexus apply -f manifests/workload-client.yaml
	@$(KUBECTL) --context=kind-coastal-containers wait --timeout=300s --for=condition=Ready pod -l app=server
	@echo "Workload deployed on the clusters."
