#!/bin/sh

set -e

if [ ! -f "/dns-conf/named.conf" ]; then
   cp /etc/bind/named.conf.recursive /dns-conf/named.conf
fi

exec /usr/sbin/named -c /etc/bind/named.conf.recursive -f
