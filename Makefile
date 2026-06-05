# note: call scripts from /scripts

.PHONY: default build builder-image binary-image container test stop clean-images clean push apply deploy

BUILDER ?= gitwebhookproxy-builder
BINARY ?= GitWebhookProxy
DOCKER_IMAGE ?= stakater/gitwebhookproxy
REGISTRY ?= docker.io

# Default platforms to build for
PLATFORMS ?= linux/amd64 linux/arm64
# Default value "dev"
DOCKER_TAG ?= dev
REPOSITORY = ${DOCKER_IMAGE}:${DOCKER_TAG}

VERSION=$(shell cat .VERSION)
BUILD=

GOCMD = go
GOFLAGS ?= $(GOFLAGS:)
LDFLAGS =

# Detect container CLI (docker or podman)
CONTAINER_CLI := $(shell command -v podman 2> /dev/null || echo docker)
BUILD_DATE := $(shell date -u +'%Y%m%d%H%M%S')

goproxy := https://proxy.golang.org
golang_fips140_version=v1.0.0

default: build test

install:
	"$(GOCMD)" mod download

build:
	"$(GOCMD)" build ${GOFLAGS} ${LDFLAGS} -o "${BINARY}"

# Legacy builder image method (single architecture)
builder-image:
	@$(CONTAINER_CLI) build --network host -t "${BUILDER}" -f build/package/Dockerfile.build .

# Legacy binary image method (single architecture)
binary-image: builder-image
	@$(CONTAINER_CLI) run --network host --rm "${BUILDER}" | $(CONTAINER_CLI) build --network host -t "${REPOSITORY}" -f build/package/Dockerfile.run -

# Multi-architecture container build
container:
	@if ! echo "$(CONTAINER_CLI)" | grep -Eq '/podman$$|^podman$$'; then \
		echo "Error: the container target for a multi-arch containers build requires 'podman' to be installed"; \
		exit 1; \
	fi
	@echo "Building multi-architecture container using $(CONTAINER_CLI)"
	@$(eval manifest_digest_f := $(shell mktemp /tmp/manifest-digest.XXXXXX))

	# Create manifest
	@if ! $(CONTAINER_CLI) manifest create "${REPOSITORY}-b$(BUILD_DATE)"; then \
		echo "FATAL: error creating manifest" 1>&2 ; \
		exit 1 ; \
	fi

	# Build for each platform
	@for platform in $(PLATFORMS); do \
		echo "Building $(REPOSITORY)-b$(BUILD_DATE) for $${platform}..." ; \
		arch=$$(basename $${platform}) ; \
		container_digest_f=$$(mktemp /tmp/container-digest-$${arch}.XXXXXX) ; \
		if ! $(CONTAINER_CLI) build \
			--manifest="${REPOSITORY}-b$(BUILD_DATE)" \
			--no-cache \
			--rm=true \
			--platform $$platform \
			--iidfile $${container_digest_f} \
			--file build/package/Dockerfile.multi \
			--env=GOFIPS140="$(golang_fips140_version)" \
			--env=GOPROXY="$(goproxy)" \
			--tag "${REPOSITORY}-b$(BUILD_DATE)-$$arch" . ;  then \
				echo "FATAL: failed building ${REPOSITORY}-b$(BUILD_DATE)-$$arch." ; \
				exit 1 ; \
		fi ; \
		DIGEST=$$(cat $${container_digest_f}) ;\
		IMAGE_ID=$$(echo "$${DIGEST}" | sed 's#sha256:##') ;\
		if [ -n "$(REGISTRY)" ]; then \
			echo "INFO: pushing ${REPOSITORY}-b$(BUILD_DATE)-$$arch ($${IMAGE_ID}) to ${REGISTRY}" ; \
			if ! $(CONTAINER_CLI) push $${IMAGE_ID} "docker://$(REGISTRY)/${REPOSITORY}-b$(BUILD_DATE)-$$arch"; then \
				echo "FATAL: failed pushing container $${IMAGE_ID} to docker://$(REGISTRY)/${REPOSITORY}-b$(BUILD_DATE)-$$arch" ; \
				exit 1 ; \
			fi ; \
		else \
			echo "WARNING: Skip pushing ${REPOSITORY}-b$(BUILD_DATE)-$$arch to registry (as REGISTRY not defined)" ; \
		fi; \
		rm -f $${container_digest_f} ; \
	done

	# Push container list manifests to registry if REGISTRY is defined
	@if [ -n "$(REGISTRY)" ]; then \
		for m in "$(REGISTRY)/${REPOSITORY}-b$(BUILD_DATE)" "$(REGISTRY)/${REPOSITORY}"; do \
			echo "Pushing manifest to $$m" ; \
			if ! $(CONTAINER_CLI) manifest push --digestfile="$(manifest_digest_f)" "${REPOSITORY}-b$(BUILD_DATE)" "docker://$$m"; then \
				echo "FATAL: error pushing list manifest $$m" 1>&2 ; \
				exit 1 ; \
			fi; \
		done ; \
	else \
		echo "Skip pushing the container list manifest to registry (as REGISTRY not defined)" ; \
	fi

	# Clean up
	@rm -f $(manifest_digest_f)
	@echo "Cleaning up dangling images"
	@for container in $$($(CONTAINER_CLI) images -q -f "dangling=true"); do \
		$(CONTAINER_CLI) rmi $${container} || true; \
	done

test:
	"$(GOCMD)" test -v ./...

stop:
	@$(CONTAINER_CLI) stop "${BINARY}" || true

clean-images: stop
	@$(CONTAINER_CLI) rmi "${BUILDER}" "${BINARY}" || true

clean:
	"$(GOCMD)" clean -i

push: ## push the latest Docker image to DockerHub
	$(CONTAINER_CLI) push $(REPOSITORY)

apply:
	kubectl apply -f deployments/manifests/

# Legacy deploy method
deploy-legacy: binary-image push apply

# New multi-arch deploy method
deploy: container apply