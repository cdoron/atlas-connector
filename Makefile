ETCD_CONNECTOR_DIR := ${PWD}
USER_ID := $(shell id -u)
GROUP_ID := $(shell id -g)

GIT_USER_ID := fybrik
GIT_REPO_ID := datacatalog-go

DOCKER_HOSTNAME ?= ghcr.io
DOCKER_NAMESPACE ?= cdoron
DOCKER_TAG ?= 0.0.0
DOCKER_NAME ?= atlas-connector

IMG := ${DOCKER_HOSTNAME}/${DOCKER_NAMESPACE}/${DOCKER_NAME}:${DOCKER_TAG}
export HELM_EXPERIMENTAL_OCI=1


all: build

.PHONY: compile
compile:
	go build .

.PHONY: build
build: compile
	docker build . -t ${IMG}; cd ..

.PHONY: docker-push
docker-push:
	docker push ${IMG}

.PHONY: push-to-kind
push-to-kind:
	kind load docker-image ${IMG}


generate-code:
	git clone https://github.com/fybrik/fybrik/
	cd fybrik && git checkout v0.7.0
	docker run --rm \
           -v ${ETCD_CONNECTOR_DIR}:/local \
           -u "${USER_ID}:${GROUP_ID}" \
           openapitools/openapi-generator-cli generate -g go-server \
           --git-user-id=${GIT_USER_ID} \
           --git-repo-id=${GIT_REPO_ID} \
           --config=/local/openapi-configs/go-server-config.yaml \
           -o /local/api \
           -i /local/fybrik/connectors/api/datacatalog.spec.yaml
	rm -Rf fybrik

patch:
	sed -i 's/\t"github.com\/gorilla\/mux"//' api/api/api_default.go
