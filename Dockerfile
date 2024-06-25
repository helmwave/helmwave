ARG GOLANG_VERSION=1.24
ARG ALPINE_VERSION=3.20

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV PROJECT=helmwave

WORKDIR ${PROJECT}

COPY go.mod go.sum ./
RUN go mod download

# Copy src code from the host and compile it
COPY cmd cmd
COPY pkg pkg
RUN go build -a -o /${PROJECT} ./cmd/${PROJECT}

### Base image with shell
FROM alpine:${ALPINE_VERSION} as base-release
RUN apk --update --no-cache add ca-certificates && update-ca-certificates
ENTRYPOINT ["/bin/helmwave"]

### Base image with shell and debugging tools
FROM base-release as base-debug-release
RUN apk --update --no-cache add jq bash
COPY --chown=root:root --chmod=0775 --from=bitnami/kubectl:latest /opt/bitnami/kubectl/bin/kubectl /bin/kubectl

### Build with goreleaser
FROM base-release as goreleaser
COPY helmwave /bin/

### Debug tag
FROM base-debug-release as debug-goreleaser
COPY helmwave /bin/

### Build in docker
FROM base-release as release
COPY --from=builder /helmwave /bin/

### Scratch with build in docker
FROM scratch as scratch-release
COPY --from=base-release /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /helmwave /bin/
ENTRYPOINT ["/bin/helmwave"]
USER 65534

### Scratch with goreleaser
FROM scratch as scratch-goreleaser
COPY --from=base-release /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY helmwave /bin/
ENTRYPOINT ["/bin/helmwave"]
USER 65534
