FROM golang:1.12-alpine as builder

WORKDIR /go/src/github.com/jiripetrlik/handy-dns
COPY . .

RUN go install -v github.com/jiripetrlik/handy-dns/cmd/handy-dns-manager

FROM alpine:latest

COPY --from=builder /go/bin/handy-dns-manager /handy-dns-manager

ENV DNS_IP 127.0.0.1
ENV ORIGIN example-domain.
ENV PRIMARY_NAME_SERVER ns1
ENV EMAIL email.example-domain.
ENV DNS_ZONE /dns-conf/domain.hosts
ENV ZONE_DATA /dns-conf/zone-data.json

CMD exec /handy-dns-manager -i "$DNS_IP" -o "$ORIGIN" -p "$PRIMARY_NAME_SERVER" -e "$EMAIL" -f "$DNS_ZONE" -d "$ZONE_DATA"
