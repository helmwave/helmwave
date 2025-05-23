ARG GOLANG_VERSION=1.23
ARG ALPINE_VERSION=3.21
ARG PROJECT=helmwave

# Bulder
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder
ENV GO111MODULE=on
ENV CGO_ENABLED=0
WORKDIR ${PROJECT}
COPY go.mod go.sum ./
RUN go mod download

# Copy src code from the host and compile it
COPY cmd cmd
COPY pkg pkg

RUN go build -a -o "/${PROJECT}" "./cmd/${PROJECT}"

## Base image with shell
FROM alpine:${ALPINE_VERSION} as base-release
RUN apk --update --no-cache add ca-certificates && update-ca-certificates
ENTRYPOINT ["/bin/${PROJECT}"]


### Base image with shell and debugging tools
FROM base-release as base-debug-release
RUN apk --update --no-cache add jq bash
COPY --chown=root:root --chmod=0775 --from=bitnami/kubectl:latest /opt/bitnami/kubectl/bin/kubectl /bin/kubectl


### Build with goreleaser
FROM base-release as goreleaser
COPY ${PROJECT} /bin/

### Debug tag, just manually
FROM base-debug-release as debug-goreleaser
COPY ${PROJECT} /bin/

### Build in docker
FROM base-release as release
COPY --from=builder /${PROJECT} /bin/

### Scratch with build in docker
FROM scratch as scratch-release
COPY --from=base-release /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /${PROJECT} /bin/
ENTRYPOINT ["/bin/${PROJECT}"]
USER 65534

### Scratch with goreleaser
FROM scratch as scratch-goreleaser
COPY --from=base-release /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ${PROJECT} /bin/
ENTRYPOINT ["/bin/${PROJECT}"]
USER 65534
