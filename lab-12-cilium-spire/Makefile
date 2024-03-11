BIN_DIR := ../bin
NO_ISSUER = 1
NO_WAIT = 1

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-42s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include ../tools/Makefile

##@ Lab

.PHONY: platform-images
platform-images: ## Load component images into the cluster
	@$(MAKE) cluster-load-image IMG=$(ENVOY_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIRE_SERVER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIRE_AGENT_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(WAIT_FOR_IT_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(HUBBLE_RELAY_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(CILIUM_OPERATOR_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(HUBBLE_UI_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(HUBBLE_UI_BACKEND_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(CILIUM_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(CILIUM_CURL_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(CILIUM_JSON_MOCK_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(NETWORK_MULTITOOL_IMAGE_TAG)

.PHONY: cluster-up
cluster-up:: cilium
	$(MAKE) platform-images

.PHONY: deploy-spire
deploy-spire: kubectl ## Deploy SPIRE on the K8s cluster
	@echo "Deploying SPIRE on the cluster..."
	$(KUBECTL) apply -f manifests/spire-server.yaml
	$(KUBECTL) apply -f manifests/spire-agent.yaml
	$(KUBECTL) -n spire wait --timeout=300s --for=condition=Ready pod -l app=spire-agent
	@echo "SPIRE deployed on the cluster."
	$(KUBECTL) exec -n spire spire-server-0 -- \
		/opt/spire/bin/spire-server entry create \
		-spiffeID spiffe://coastal-containers.example/agent/spire-agent \
		-selector k8s_psat:cluster:kind-kind \
		-selector k8s_psat:agent_ns:spire \
		-selector k8s_psat:agent_sa:spire-agent \
		-node

test-workload-networking: ## Test networking
	@$(KUBECTL) wait --for=condition=Ready pod/client pod/server
	$(call test-connectivity,client,server,"Workload", \
		$$( $(KUBECTL) get pods -o=jsonpath='{.items[?(@.metadata.name=="server")].status.podIP}' ))
	$(call test-connectivity,server,client,"Undesirable Workload", \
		$$( $(KUBECTL) get pods -o=jsonpath='{.items[?(@.metadata.name=="client")].status.podIP}' ))
	$(call test-connectivity,client,1.1.1.1,"External","1.1.1.1")
	@$(CILIUM_CLI) hubble port-forward & \
	sleep 1;\
	$(HUBBLE_CLI) observe --type policy-verdict; \
	pkill -f "cilium hubble port-forward"

define test-connectivity
	@$(KUBECTL) exec $1 -- ping -q -c 2 -W 0.2 $4 > /dev/null ;\
	if [ $$? -eq 0 ]; then \
		echo "✅ $3 Ping was successful - ICMP between $1 > $2 works!"; \
	else \
		echo "❌ $3 Ping failed - ICMP between $1 > $2 does not work!"; \
	fi ;

	@$(KUBECTL) exec $1 -- curl -fs --connect-timeout 0.2 $4 > /dev/null ;\
	if [ $$? -eq 0 ]; then \
		echo "✅ $3 curl was successful - HTTP between $1 > $2 works!"; \
	else \
		echo "❌ $3 curl failed - HTTP between $1 > $2 does not work!"; \
	fi
endef
