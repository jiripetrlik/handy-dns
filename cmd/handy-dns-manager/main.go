package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jiripetrlik/handy-dns/internal/app/rest"
	"github.com/jiripetrlik/handy-dns/internal/app/zonefile"
)

func main() {
	zoneNamePtr := flag.String("n", "example-domain", "Domain name")
	zoneFilePtr := flag.String("f", "example-domain.hosts", "Zone file")
	zoneDataPtr := flag.String("d", "example-domain.json", "Data about zone")
	flag.Parse()

	log.Printf(
		"Starting handy-dns-manager for domain %v. Zonefile=%v and zone data=%v",
		*zoneNamePtr,
		*zoneFilePtr,
		*zoneDataPtr,
	)

	dnsZone := zonefile.DNSZone{
		*zoneNamePtr,
		*zoneFilePtr,
		*zoneDataPtr,
	}
	dnsZone.Initialize()

	restServer := rest.HandyDnsRestServer{
		&dnsZone,
	}

	restServer.HandleRestAPI()
	http.ListenAndServe(":8080", nil)
}
