FROM golang:1.15-alpine AS builder

LABEL maintainer="helmwave+zhilyaev.dmitriy@gmail.com"
LABEL name="helmwave"

# enable Go modules support
ENV GO111MODULE=on

WORKDIR $GOPATH/src/github.com/zhilyaev/helmwave

COPY go.mod go.sum ./
RUN go mod download

# Copy src code from the host and compile it
COPY cmd cmd
COPY pkg pkg
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /helmwave github.com/zhilyaev/helmwave/cmd/helmwave

FROM alpine:3.13
RUN apk --no-cache add ca-certificates
COPY --from=builder /helmwave /bin
ENTRYPOINT ["/bin/helmwave"]
