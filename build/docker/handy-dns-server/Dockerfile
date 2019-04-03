FROM alpine:latest

RUN apk add bind
COPY start.sh /start.sh
VOLUME /dns-conf
EXPOSE 53/udp

CMD ["/start.sh"]
