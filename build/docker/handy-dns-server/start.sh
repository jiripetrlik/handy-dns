#!/bin/sh

set -e

if [ ! -f "/dns-conf/named.conf" ]; then
   cp /etc/bind/named.conf.recursive /dns-conf/named.conf
   sed -i 's/127.0.0.1\/32/0.0.0.0\/0/' /dns-conf/named.conf
   sed -i 's/127.0.0.1;/any;/' /dns-conf/named.conf
fi

exec /usr/sbin/named -c /dns-conf/named.conf -f
