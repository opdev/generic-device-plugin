ARG ARCH=amd64
ARG OS=linux

FROM docker.io/golang:1.21 AS builder
ARG ARCH
ARG OS
COPY . /go/src/edge-device-plugin
WORKDIR /go/src/edge-device-plugin
RUN make build OS=${OS} ARCH=${ARCH}

FROM registry.access.redhat.com/ubi9/ubi-micro
LABEL name="edge-device-plugin" \
      maintainer="Edmund Ochieng" \
      version="alpha" \
      summary="Device plugin for Litmus Edge" \
      description="A plugin to expose usb and serial devices to the Litmus application"
COPY --from=builder /go/src/edge-device-plugin/bin/edge-deviceplugin /usr/local/bin/edge-device-plugin
ENTRYPOINT ["/usr/local/bin/edge-device-plugin"]
