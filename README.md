# Handy DNS
Handy DNS project provides containers with DNS server and tools for managing DNS zone using
simple REST API.

## Docker images

### DNS server image

DNS server image contains pre-installed bind server. It supports two
possible modes - recursive and iterative. Recursive mode is enabled
by default if no parameters are specified.
```
# Run recursive DNS server
docker run -d --name handy-dns-server -p 53:53/udp jiripetrlik/handy-dns-server
```

In iterative mode the DNS server works as authoritative for specific zone.
To run in iterative mode it is necessary to specify DNS zone file and
domain name.

Parameters:
* DOMAIN_NAME - name of the managed domain
* DNS_ZONE - path to zone file (default is `/dns-conf/domain.hosts`)

### DNS Manager image

DNS manager is an application written in Go which provides simple REST API
to manage DNS zone file. Use following commands to run iterative DNS
server with zone file managed by DNS manager:

```
# Create volume for DNS configuration files
docker volume create dns-config

# Run DNS manager for "example-domain"
docker run -d --name handy-dns-manager -e "DNS_IP=dns-server-ip" -e "ORIGIN=example-domain." --mount 'type=volume,src=dns-config,dst=/dns-conf' -p 8080:8080 jiripetrlik/handy-dns-manager

# Run DNS server
docker run -d --name handy-dns-server -e "DOMAIN_NAME=example-domain" --mount 'type=volume,src=dns-config,dst=/dns-conf' -p 53:53/udp jiripetrlik/handy-dns-server
```

#### Additional environment variables
* ZONE_DATA ... JSON file with data about DNS zone (default value: /dns-conf/zone-data.json)
* EMAIL ... Hostmaster email (default value: email.example-domain.)
* DNS_ZONE ... Zone file (default value: /dns-conf/domain.hosts)
* PRIMARY_NAME_SERVER ... name of primary nameserver (default value: ns1)
* ORIGIN ... Domain origin (default value: example-domain.)
* HTPASSWD ... Htpasswd file
* CERT_FILE ... Cert file for HTTPS
* KEY_FILE ... Key file for HTTPS

### Net tools image

[Nettools image](https://hub.docker.com/r/jiripetrlik/nettools) is an Ubuntu based image with pre-installed network utilities:
* ping
* dig
* nslookup
* traceroute
* telnet
* ssh
* curl

Dockerfile can be found [here](https://github.com/jiripetrlik/handy-dns/blob/master/build/docker/nettools/Dockerfile).
It is possible to run Docker container with nettools using following command:
```
docker run -i -t --rm jiripetrlik/nettools /bin/bash
```

## CURL examples for DNS Manager REST API

```
# List items
curl http://localhost:8080/api/list

# Add item
curl "http://localhost:8080/api/create?name=machine1&class=IN&itemType=A&data=127.0.0.1"

# Update item
curl "http://localhost:8080/api/update?id=1&name=machine5&class=IN&itemType=A&data=127.0.0.1"

# Delete item
curl "http://localhost:8080/api/delete?id=1"
```

## Building Docker images

```
# Build handy-dns-server image
docker build -f build/docker/handy-dns-server/Dockerfile -t handy-dns-server .

# Build handy-dns-manager image
docker build -f build/docker/handy-dns-manager/Dockerfile -t handy-dns-manager .

# Build nettools image
docker build -f build/docker/nettools/Dockerfile -t nettools .
```
