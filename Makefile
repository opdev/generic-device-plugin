OS ?= linux
ARCH ?= amd64
IMAGE_BUILDER ?= podman
IMAGE_REPO ?= quay.io/opdev
VERSION ?= latest

build:
	GOOS=$(OS) GOARCH=$(ARCH) go build -o bin/generic-device-plugin cmd/plugin/main.go

.PHONY: image-build
image-build:
	$(IMAGE_BUILDER) build -t $(IMAGE_REPO)/generic-device-plugin:$(VERSION) .
