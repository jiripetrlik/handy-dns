#!/bin/sh -x

set -e

if [ ! -f "/dns-conf/named.conf" ]; then
   if [ -n "$DOMAIN_NAME" ]; then
      cp /etc/bind/named.conf.authoritative /dns-conf/named.conf
      sed -i 's/127.0.0.1;/any;/' /dns-conf/named.conf
      echo "zone \"$DOMAIN_NAME\" IN {" >> /dns-conf/named.conf
      echo "   type master;" >> /dns-conf/named.conf
      echo "   file \"$DNS_ZONE\";" >> /dns-conf/named.conf
      echo "};" >> /dns-conf/named.conf
   else
      cp /etc/bind/named.conf.recursive /dns-conf/named.conf
      sed -i 's/127.0.0.1\/32/0.0.0.0\/0/' /dns-conf/named.conf
      sed -i 's/127.0.0.1;/any;/' /dns-conf/named.conf
   fi
fi

/reload.sh &
exec /usr/sbin/named -c /dns-conf/named.conf -d 1 -f
