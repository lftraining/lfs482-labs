BIN_DIR := ../bin

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-42s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include ../tools/Makefile

##@ Lab

.PHONY: platform-images
platform-images: ## Load component images into the cluster
	@$(MAKE) cluster-load-image IMG=$(CONTOUR_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(ENVOY_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(CERT_MANAGER_CAINJECTOR_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(CERT_MANAGER_CONTROLLER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(CERT_MANAGER_WEBHOOK_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIRE_SERVER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIRE_AGENT_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(WAIT_FOR_IT_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIFFE_CSI_DRIVER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIRE_CONTROLLER_MANAGER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIFFE_CSI_DRIVER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(SPIRE_OIDC_DISCOVERY_PROVIDER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(CSI_NODE_DRIVER_REGISTRAR_IMAGE_TAG)

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
     		-spiffeID spiffe://coastal-containers.example/workload \
     		-selector k8s:ns:default \
     		-selector k8s:sa:default

.PHONY: view-oidc-discovery-provider-config
view-oidc-discovery-provider-config: kubectl ## View the SPIRE OIDC Discovery Provider configuration file.
	@$(KUBECTL) -n spire get cm spire-spiffe-oidc-discovery-provider -ojsonpath='{.data.oidc-discovery-provider\.conf}'

.PHONY: view-controller-manager-config
view-controller-manager-config: kubectl ## View the SPIRE Controller Manager configuratio file.
	@$(KUBECTL) -n spire get configmaps spire-controller-manager -ojsonpath='{.data.controller-manager-config\.yaml}'

.PHONY: view-spiffe-clusterid
view-spiffe-clusterid: kubectl ## View the SPIRE Controller Manager configuratio file.
	@$(KUBECTL) get clusterspiffeids.spire.spiffe.io spire-spire-default -ojsonpath={.spec} | jq

.PHONY: deploy-workload
deploy-workload: kubectl ## Deploy the workload job
	$(KUBECTL) apply -f workload/job.yaml

.PHONY: view-logs
view-logs: kubectl ## View the logs of the workload job
	$(KUBECTL) logs jobs/workload

##@ Infrastructure

.PHONY: cluster-up
cluster-up:: ## Launch a Kind cluster with Cert Manager and Contour ingress
	$(MAKE) platform-images
	@echo "Installing Contour"
	@$(KUBECTL) apply -f contour/contour.yaml
	@echo "Installing Cert Manager"
	@$(KUBECTL) apply -f cert-manager/cert-manager.yaml
	@echo "Waiting for Cert Manager Webhook"
	@$(KUBECTL) -n cert-manager wait --for=condition=Ready --timeout=300s pods -l=app.kubernetes.io/component=webhook
	@echo "Creating Self Signed Cluster Issuer"
	@$(KUBECTL) apply -f $(ISSUER_CONFIG)
