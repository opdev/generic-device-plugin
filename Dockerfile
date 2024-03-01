ARG ARCH=amd64
ARG OS=linux

FROM docker.io/golang:1.21 AS builder
ARG ARCH
ARG OS
COPY . /go/src/generic-device-plugin
WORKDIR /go/src/generic-device-plugin
RUN make build OS=${OS} ARCH=${ARCH}

FROM registry.access.redhat.com/ubi9/ubi-micro
LABEL name="generic-device-plugin" \
      maintainer="Edmund Ochieng" \
      version="alpha" \
      summary="Generic Device plugin for USB Devices" \
      description="A generic device plugin to expose usb and serial devices to Kubelet"
COPY --from=builder /go/src/generic-device-plugin/bin/generic-device-plugin /usr/local/bin/generic-device-plugin
ENTRYPOINT ["/usr/local/bin/generic-device-plugin"]
