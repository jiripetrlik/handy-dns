FROM golang:1.12-alpine as builder

WORKDIR /go/src/handy-dns
COPY . .

RUN go install -v ./cmd/handy-dns-manager

FROM alpine:latest

COPY --from=builder /go/bin/handy-dns-manager /handy-dns-manager

CMD ["/handy-dns-manager"]
