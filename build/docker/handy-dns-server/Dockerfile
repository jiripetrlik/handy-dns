FROM alpine:latest

ENV DNS_ZONE /dns-conf/domain.hosts

RUN apk add bind
COPY build/docker/handy-dns-server/*.sh /
VOLUME /dns-conf
EXPOSE 53/udp

CMD ["/start.sh"]
