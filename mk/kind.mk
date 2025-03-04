EXAMPLE_NAMESPACE ?= kuma-example
KIND_KUBECONFIG_DIR ?= $(HOME)/.kube
KIND_KUBECONFIG ?= $(KIND_KUBECONFIG_DIR)/kind-kuma-config
KIND_CLUSTER_NAME ?= kuma

ifeq ($(KUMACTL_INSTALL_USE_LOCAL_IMAGES),true)
	KUMACTL_INSTALL_CONTROL_PLANE_IMAGES := --control-plane-image=$(KUMA_CP_DOCKER_IMAGE_NAME) --dataplane-image=$(KUMA_DP_DOCKER_IMAGE_NAME) --dataplane-init-image=$(KUMA_INIT_DOCKER_IMAGE_NAME)
else
	KUMACTL_INSTALL_CONTROL_PLANE_IMAGES :=
endif
ifeq ($(KUMACTL_INSTALL_USE_LOCAL_IMAGES),true)
	KUMACTL_INSTALL_METRICS_IMAGES := --kuma-prometheus-sd-image=$(KUMA_PROMETHEUS_SD_DOCKER_IMAGE_NAME)
else
	KUMACTL_INSTALL_METRICS_IMAGES :=
endif

define KIND_EXAMPLE_DATAPLANE_MESH
$(shell KUBECONFIG=$(KIND_KUBECONFIG) kubectl -n $(EXAMPLE_NAMESPACE) exec $$(kubectl -n $(EXAMPLE_NAMESPACE) get pods -l app=example-app -o=jsonpath='{.items[0].metadata.name}') -c kuma-sidecar printenv KUMA_DATAPLANE_MESH)
endef
define KIND_EXAMPLE_DATAPLANE_NAME
$(shell KUBECONFIG=$(KIND_KUBECONFIG) kubectl -n $(EXAMPLE_NAMESPACE) exec $$(kubectl -n $(EXAMPLE_NAMESPACE) get pods -l app=example-app -o=jsonpath='{.items[0].metadata.name}') -c kuma-sidecar printenv KUMA_DATAPLANE_NAME)
endef

CI_KIND_VERSION ?= v0.8.1
CI_KUBERNETES_VERSION ?= v1.18.2@sha256:7b27a6d0f2517ff88ba444025beae41491b016bc6af573ba467b70c5e8e0d85f

KIND_PATH := $(CI_TOOLS_DIR)/kind

KUMA_MODE ?= standalone
KUMA_NAMESPACE ?= kuma-system

.PHONY: ${KIND_KUBECONFIG_DIR}
${KIND_KUBECONFIG_DIR}:
	@mkdir -p ${KIND_KUBECONFIG_DIR}

.PHONY: kind/start
kind/start: ${KIND_KUBECONFIG_DIR}
	@kind get clusters | grep $(KIND_CLUSTER_NAME) >/dev/null 2>&1 && echo "Kind cluster already running." && exit 0 || \
		(kind create cluster \
			--name "$(KIND_CLUSTER_NAME)" \
			--image=kindest/node:$(CI_KUBERNETES_VERSION) \
			--kubeconfig $(KIND_KUBECONFIG) \
			--quiet --wait 120s && \
		until \
			KUBECONFIG=$(KIND_KUBECONFIG) kubectl wait -n kube-system --timeout=5s --for condition=Ready --all pods ; \
		do echo "Waiting for the cluster to come up" && sleep 1; done )
	@echo
	@echo '>>> You need to manually run the following command in your shell: >>>'
	@echo
	@echo export KUBECONFIG="${KIND_KUBECONFIG}"
	@echo
	@echo '<<< ------------------------------------------------------------- <<<'
	@echo

.PHONY: kind/stop
kind/stop:
	@kind delete cluster --name $(KIND_CLUSTER_NAME)

.PHONY: kind/stop/all
kind/stop/all:
	@kind delete clusters --all

.PHONY: kind/load/control-plane
kind/load/control-plane:
	@kind load docker-image $(KUMA_CP_DOCKER_IMAGE) --name=$(KIND_CLUSTER_NAME)

.PHONY: kind/load/kuma-dp
kind/load/kuma-dp:
	@kind load docker-image $(KUMA_DP_DOCKER_IMAGE) --name=$(KIND_CLUSTER_NAME)

.PHONY: kind/load/kuma-init
kind/load/kuma-init:
	@kind load docker-image $(KUMA_INIT_DOCKER_IMAGE) --name=$(KIND_CLUSTER_NAME)

.PHONY: kind/load/kuma-prometheus-sd
kind/load/kuma-prometheus-sd:
	@kind load docker-image $(KUMA_PROMETHEUS_SD_DOCKER_IMAGE) --name=$(KIND_CLUSTER_NAME)

.PHONY: kind/load/images
kind/load/images: kind/load/control-plane kind/load/kuma-dp kind/load/kuma-init kind/load/kuma-prometheus-sd

.PHONY: kind/load
kind/load: image/kuma-cp image/kuma-dp image/kuma-init image/kuma-prometheus-sd kind/load/images

.PHONY: kind/deploy/kuma
kind/deploy/kuma: build/kumactl kind/load
	@${BUILD_ARTIFACTS_DIR}/kumactl/kumactl install --mode $(KUMA_MODE) control-plane $(KUMACTL_INSTALL_CONTROL_PLANE_IMAGES) | KUBECONFIG=$(KIND_KUBECONFIG)  kubectl apply -f -
	@KUBECONFIG=$(KIND_KUBECONFIG) kubectl wait --timeout=60s --for=condition=Available -n $(KUMA_NAMESPACE) deployment/kuma-control-plane
	@KUBECONFIG=$(KIND_KUBECONFIG) kubectl wait --timeout=60s --for=condition=Ready -n $(KUMA_NAMESPACE) pods -l app=kuma-control-plane
	@KUBECONFIG=$(KIND_KUBECONFIG) kubectl delete -n $(EXAMPLE_NAMESPACE) pod -l app=example-app
	@until \
    	KUBECONFIG=$(KIND_KUBECONFIG) kubectl wait -n kube-system --timeout=5s --for condition=Ready --all pods ; \
    do \
    	echo "Waiting for the cluster to come up" && sleep 1; \
    done

.PHONY: kind/deploy/kuma/global
kind/deploy/kuma/global: KUMA_MODE=global
kind/deploy/kuma/global: kind/deploy/kuma

.PHONY: kind/deploy/kuma/local
kind/deploy/kuma/local: KUMA_MODE=local
kind/deploy/kuma/local: kind/deploy/kuma

.PHONY: kind/deploy/metrics
kind/deploy/metrics: build/kumactl
	@${BUILD_ARTIFACTS_DIR}/kumactl/kumactl install metrics $(KUMACTL_INSTALL_METRICS_IMAGES) | kubectl apply -f -
	@KUBECONFIG=$(KIND_KUBECONFIG) kubectl wait --timeout=60s --for=condition=Ready -n kuma-metrics pods -l app=prometheus

.PHONY: kind/deploy/example-app
kind/deploy/example-app:
	@KUBECONFIG=$(KIND_KUBECONFIG) kubectl create namespace $(EXAMPLE_NAMESPACE) || true
	@KUBECONFIG=$(KIND_KUBECONFIG) kubectl annotate namespace $(EXAMPLE_NAMESPACE) kuma.io/sidecar-injection=enabled --overwrite
	@KUBECONFIG=$(KIND_KUBECONFIG) kubectl apply -n $(EXAMPLE_NAMESPACE) -f dev/examples/k8s/example-app/example-app.yaml
	@KUBECONFIG=$(KIND_KUBECONFIG) kubectl wait --timeout=120s --for=condition=Available -n $(EXAMPLE_NAMESPACE) deployment/example-app
	@KUBECONFIG=$(KIND_KUBECONFIG) kubectl wait --timeout=60s --for=condition=Ready -n $(EXAMPLE_NAMESPACE) pods -l app=example-app

.PHONY: run/k8s
run/k8s: fmt vet ## Dev: Run Control Plane locally in Kubernetes mode
	@KUBECONFIG=$(KIND_KUBECONFIG) make crd/upgrade -C pkg/plugins/resources/k8s/native
	KUBECONFIG=$(KIND_KUBECONFIG) \
	KUMA_SDS_SERVER_GRPC_PORT=$(SDS_GRPC_PORT) \
	KUMA_GRPC_PORT=$(CP_GRPC_PORT) \
	KUMA_ENVIRONMENT=kubernetes \
	KUMA_STORE_TYPE=kubernetes \
	KUMA_SDS_SERVER_TLS_CERT_FILE=app/kuma-cp/cmd/testdata/tls.crt \
	KUMA_SDS_SERVER_TLS_KEY_FILE=app/kuma-cp/cmd/testdata/tls.key \
	KUMA_RUNTIME_KUBERNETES_ADMISSION_SERVER_PORT=$(CP_K8S_ADMISSION_PORT) \
	KUMA_RUNTIME_KUBERNETES_ADMISSION_SERVER_CERT_DIR=app/kuma-cp/cmd/testdata \
	$(GO_RUN) ./app/kuma-cp/main.go run --log-level=debug

run/example/envoy/k8s: EXAMPLE_DATAPLANE_MESH=$(KIND_EXAMPLE_DATAPLANE_MESH)
run/example/envoy/k8s: EXAMPLE_DATAPLANE_NAME=$(KIND_EXAMPLE_DATAPLANE_NAME)
run/example/envoy/k8s: run/example/envoy
