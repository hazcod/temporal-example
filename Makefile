NAME = "temporal"
DOCKER_REGISTRY = "eu.gcr.io/ironsecurity"
STAGE = "dev"
NAMESPACE = $(STAGE)-$(NAME)
DEFAULT_RELEASE = "0.0"
RELEASE = $(DEFAULT_RELEASE)
CONTEXT = "docker-desktop"

K8S_DIR="kubernetes"
DOCKER_DIR=$(K8S_DIR)/docker
CHART_DIR=$(K8S_DIR)/helm

all: build up

clean:
	helm --kube-context $(CONTEXT) uninstall --namespace $(NAMESPACE) $(NAME) || true

build:
	DOCKER_BUILDKIT=1 docker build -t $(DOCKER_REGISTRY)/$(NAME):$(RELEASE) -f $(DOCKER_DIR)/$(STAGE).Dockerfile .

up:
	# helm dependency build $(CHART_DIR)
	kubectl --context $(CONTEXT) create namespace $(NAMESPACE) || true
	helm upgrade -v 3 --debug --wait --install \
		--kube-context $(CONTEXT) \
 		--namespace $(NAMESPACE) \
		--values $(K8S_DIR)/$(STAGE).yaml \
		--set devWorkingDirectory=$(shell pwd) \
		--set image.registry=$(DOCKER_REGISTRY) \
		--set namespace=$(NAMESPACE) \
		--version $(RELEASE) \
		$(NAME) $(CHART_DIR)
	helm --kube-context $(CONTEXT) test --logs --namespace $(NAMESPACE) $(NAME)
	kubectl --context $(CONTEXT) get pods --namespace $(NAMESPACE)

logs:
	kubectl --context $(CONTEXT) logs --max-log-requests=100 --namespace $(NAMESPACE) -f -l app=$(NAME) --all-containers

push:
	@if [ $(RELEASE) = $(DEFAULT_RELEASE) ]; then \
		echo "cannot push for dev";\
		exit 1;\
	fi
	docker push $(DOCKER_REGISTRY)/$(NAME):$(RELEASE)

test:
	go test -v ./...

scan:
	trivy image $(DOCKER_REGISTRY)/$(NAME):$(RELEASE)
