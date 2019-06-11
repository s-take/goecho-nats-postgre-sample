FROM golang:1.10.2-alpine3.7 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/s-take/goecho-nats-postgre-sample

COPY Gopkg.lock Gopkg.toml ./
COPY vendor vendor
COPY retry retry
COPY db db
COPY schema schema
COPY event event
COPY api-query-service api-query-service
COPY command-service command-service

RUN go install ./...

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .
