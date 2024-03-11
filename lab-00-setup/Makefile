BIN_DIR := ../bin

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include ../tools/Makefile

.PHONY: all
all: binaries kind-image workload-images contour-images cert-manager-images spire-images opa-images cilium-images other-images

.PHONY: binaries
binaries: cfssl helm istio kind kubectl opa spire cilium ## Download the required binaries

.PHONY: workload-images
workload-images: ## Download images for workloads
	@docker pull $(GOLANG_BUILDER_IMAGE)
	@docker pull $(GOLANG_IMAGE)
	@docker pull $(PYTHON_IMAGE)

.PHONY: kind-image
kind-image: ## Pull the kind node image
	@docker pull $(KIND_IMAGE)

.PHONY: contour-images
contour-images: ## Pull the contour images and tag them for use in kind
	@docker pull $(CONTOUR_IMAGE)
	@docker tag $(CONTOUR_IMAGE) $(CONTOUR_IMAGE_TAG)
	@docker pull $(ENVOY_IMAGE)
	@docker tag $(ENVOY_IMAGE) $(ENVOY_IMAGE_TAG)

.PHONY: cert-manager-images
cert-manager-images: ## Pull the cert manager images and tag them for use in kind
	@docker pull $(CERT_MANAGER_CAINJECTOR_IMAGE)
	@docker tag $(CERT_MANAGER_CAINJECTOR_IMAGE) $(CERT_MANAGER_CAINJECTOR_IMAGE_TAG)
	@docker pull $(CERT_MANAGER_CONTROLLER_IMAGE)
	@docker tag $(CERT_MANAGER_CONTROLLER_IMAGE) $(CERT_MANAGER_CONTROLLER_IMAGE_TAG)
	@docker pull $(CERT_MANAGER_WEBHOOK_IMAGE)
	@docker tag $(CERT_MANAGER_WEBHOOK_IMAGE) $(CERT_MANAGER_WEBHOOK_IMAGE_TAG)

.PHONY: spire-images
spire-images: ## Pull the spire images and tag them for use in kind
	@docker pull $(SPIRE_SERVER_IMAGE)
	@docker tag $(SPIRE_SERVER_IMAGE) $(SPIRE_SERVER_IMAGE_TAG)
	@docker pull $(SPIRE_AGENT_IMAGE)
	@docker tag $(SPIRE_AGENT_IMAGE) $(SPIRE_AGENT_IMAGE_TAG)
	@docker pull $(WAIT_FOR_IT_IMAGE)
	@docker tag $(WAIT_FOR_IT_IMAGE) $(WAIT_FOR_IT_IMAGE_TAG)
	@docker pull $(SPIFFE_HELPER_IMAGE)
	@docker tag $(SPIFFE_HELPER_IMAGE) $(SPIFFE_HELPER_IMAGE_TAG)
	@docker pull $(SPIRE_CONTROLLER_MANAGER_IMAGE)
	@docker tag $(SPIRE_CONTROLLER_MANAGER_IMAGE) $(SPIRE_CONTROLLER_MANAGER_IMAGE_TAG)
	@docker pull $(SPIRE_OIDC_DISCOVERY_PROVIDER_IMAGE)
	@docker tag $(SPIRE_OIDC_DISCOVERY_PROVIDER_IMAGE) $(SPIRE_OIDC_DISCOVERY_PROVIDER_IMAGE_TAG)
	@docker pull $(SPIFFE_CSI_DRIVER_IMAGE)
	@docker tag $(SPIFFE_CSI_DRIVER_IMAGE) $(SPIFFE_CSI_DRIVER_IMAGE_TAG)
	@docker pull $(CSI_NODE_DRIVER_REGISTRAR_IMAGE)
	@docker tag $(CSI_NODE_DRIVER_REGISTRAR_IMAGE) $(CSI_NODE_DRIVER_REGISTRAR_IMAGE_TAG)

.PHONY: opa-images
opa-images: ## Pull the opa images and tag them for use in kind
	@docker pull $(OPA_IMAGE)
	@docker tag $(OPA_IMAGE) $(OPA_IMAGE_TAG)
	@docker pull $(OPA_PROXY_INIT_IMAGE)
	@docker tag $(OPA_PROXY_INIT_IMAGE) $(OPA_PROXY_INIT_IMAGE_TAG)
	@docker pull $(OPA_ENVOY_IMAGE)
	@docker tag $(OPA_ENVOY_IMAGE) $(OPA_ENVOY_IMAGE_TAG)

.PHONY: cilium-images
cilium-images: ## Pull the cilium images and tag them for use in kind
	@docker pull $(CILIUM_IMAGE)
	@docker tag $(CILIUM_IMAGE) $(CILIUM_IMAGE_TAG)
	@docker pull $(CILIUM_OPERATOR_IMAGE)
	@docker tag $(CILIUM_OPERATOR_IMAGE) $(CILIUM_OPERATOR_IMAGE_TAG)
	@docker pull $(CILIUM_CURL_IMAGE)
	@docker tag $(CILIUM_CURL_IMAGE) $(CILIUM_CURL_IMAGE_TAG)
	@docker pull $(CILIUM_JSON_MOCK_IMAGE)
	@docker tag $(CILIUM_JSON_MOCK_IMAGE) $(CILIUM_JSON_MOCK_IMAGE_TAG)
	@docker pull $(HUBBLE_RELAY_IMAGE)
	@docker tag $(HUBBLE_RELAY_IMAGE) $(HUBBLE_RELAY_IMAGE_TAG)
	@docker pull $(HUBBLE_UI_IMAGE)
	@docker tag $(HUBBLE_UI_IMAGE) $(HUBBLE_UI_IMAGE_TAG)
	@docker pull $(HUBBLE_UI_BACKEND_IMAGE)
	@docker tag $(HUBBLE_UI_BACKEND_IMAGE) $(HUBBLE_UI_BACKEND_IMAGE_TAG)

.PHONY: other-images
other-images:  ## Pull the other images and tag them for use in kind
	@docker pull $(POSTGRES_IMAGE)
	@docker tag $(POSTGRES_IMAGE) $(POSTGRES_IMAGE_TAG)
	@docker pull $(ALPINE_CURL_IMAGE)
	@docker tag $(ALPINE_CURL_IMAGE) $(ALPINE_CURL_IMAGE_TAG)
	@docker pull $(NETWORK_MULTITOOL_IMAGE)
	@docker tag $(NETWORK_MULTITOOL_IMAGE) $(NETWORK_MULTITOOL_IMAGE_TAG)
	@docker pull $(BUSYBOX_IMAGE)
	@docker tag $(BUSYBOX_IMAGE) $(BUSYBOX_IMAGE_TAG)
