package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jiripetrlik/handy-dns/internal/app/dnszone"
	"github.com/jiripetrlik/handy-dns/internal/app/rest"
)

func main() {
	ipPtr := flag.String("i", "127.0.0.1", "IP of primary nameserver")
	originPtr := flag.String("o", "example-domain.", "Domain origin")
	primaryNameServerPtr := flag.String("p", "ns1", "Primary name server")
	emailPtr := flag.String("e", "email.example-domain.", "Hostmaster email")
	dnszonePtr := flag.String("f", "example-domain.hosts", "Zone file")
	zoneDataPtr := flag.String("d", "example-domain.json", "Data about zone")
	flag.Parse()

	log.Printf(
		"Starting handy-dns-manager for domain %v. dnszone=%v and zone data=%v",
		*originPtr,
		*dnszonePtr,
		*zoneDataPtr,
	)

	dnsZone := dnszone.DNSZone{
		*dnszonePtr,
		*zoneDataPtr,
	}
	dnsZone.Initialize(*ipPtr, *originPtr, *primaryNameServerPtr, *emailPtr)

	restServer := rest.HandyDnsRestServer{
		&dnsZone,
	}

	restServer.HandleRestAPI()
	http.ListenAndServe(":8080", nil)
}
