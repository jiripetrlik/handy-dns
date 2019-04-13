package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jiripetrlik/handy-dns/internal/app/rest"
)

func main() {
	zoneNamePtr := flag.String("n", "example-domain", "Domain name")
	zoneFilePtr := flag.String("f", "example-domain.hosts", "Zone file")
	zoneDataPtr := flag.String("d", "example-domain.json", "Data about zone")
	flag.Parse()

	log.Printf(
		"Starint handy-dns-manager for domain %v. Zonefile=%v and zone data=%v",
		*zoneNamePtr,
		*zoneFilePtr,
		*zoneDataPtr,
	)

	rest.HandleRestAPI()
	http.ListenAndServe(":8080", nil)
}
