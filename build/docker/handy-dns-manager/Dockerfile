FROM golang:1.14-alpine as builder

RUN apk update && apk add git openssh && rm -rf /var/cache/apk/*

WORKDIR /handy-dns
COPY . .

RUN go build

FROM alpine:latest

COPY --from=builder /handy-dns/handy-dns-manager /handy-dns-manager

ENV DNS_IP 127.0.0.1
ENV ORIGIN example-domain.
ENV PRIMARY_NAME_SERVER ns1
ENV EMAIL email.example-domain.
ENV DNS_ZONE /dns-conf/domain.hosts
ENV ZONE_DATA /dns-conf/zone-data.json

RUN mkdir /dns-conf

CMD exec /handy-dns-manager -i "$DNS_IP" -o "$ORIGIN" -p "$PRIMARY_NAME_SERVER" \
     -e "$EMAIL" -f "$DNS_ZONE" -d "$ZONE_DATA" -s "$HTPASSWD" \
     -certfile "$CERT_FILE" -keyfile "$KEY_FILE"
