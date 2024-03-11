BIN_DIR := ../bin

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include ../tools/Makefile

##@ Lab

.PHONY: platform-images
platform-images: ## Load component images into the cluster
	@$(MAKE) cluster-load-image IMG=$(CONTOUR_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(ENVOY_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(CERT_MANAGER_CAINJECTOR_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(CERT_MANAGER_CONTROLLER_IMAGE_TAG)
	@$(MAKE) cluster-load-image IMG=$(CERT_MANAGER_WEBHOOK_IMAGE_TAG)

.PHONY: workload-images
workload-images: ## Build and load workload images into the cluster
	$(MAKE) cluster-build-load-image DIR=backend
	$(MAKE) cluster-build-load-image DIR=frontend

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
