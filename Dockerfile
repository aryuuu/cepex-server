FROM golang:1.17-alpine AS builder

WORKDIR /go/src/app
# copy src
COPY . .
# install deps
RUN go get -d -v ./...
# compile binary
RUN go build -o cepex-server

FROM alpine:latest AS base

COPY --from=builder /go/src/app/cepex-server /cepex-server
CMD ["/cepex-server"]
