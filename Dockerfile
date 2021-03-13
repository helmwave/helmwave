ARG GOLANG_VERSION=1.15
ARG ALPINE_VERSION=3.13

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder

LABEL maintainer="helmwave+zhilyaev.dmitriy@gmail.com"
MAINTAINER "helmwave+zhilyaev.dmitriy@gmail.com"
LABEL name="helmwave"

# enable Go modules support
ENV GO111MODULE=on

WORKDIR helmwave

COPY go.mod go.sum ./
RUN go mod download

# Copy src code from the host and compile it
COPY cmd cmd
COPY pkg pkg
RUN CGO_ENABLED=0 go build -a -o /helmwave ./cmd/helmwave

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
