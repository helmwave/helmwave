ARG GOLANG_VERSION=1.16
ARG ALPINE_VERSION=3.14

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder

LABEL maintainer="helmwave+zhilyaev.dmitriy@gmail.com"
LABEL name="helmwave"

# enable Go modules support
ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR helmwave

COPY go.mod go.sum ./
RUN go mod download

# Copy src code from the host and compile it
COPY cmd cmd
COPY pkg pkg
RUN go build -a -o /helmwave ./cmd/helmwave

###
FROM alpine:${ALPINE_VERSION} as base-release
RUN apk --no-cache add ca-certificates
ENTRYPOINT ["/bin/helmwave"]

###
FROM base-release as goreleaser
COPY helmwave /bin/

###
FROM base-release
COPY --from=builder /helmwave /bin/
