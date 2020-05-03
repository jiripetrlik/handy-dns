FROM debian:10

RUN apt-get update \
    && apt-get install -y curl inetutils-ping net-tools dnsutils traceroute ssh telnet less vim mc \
    && rm -rf /var/lib/apt/lists/*
