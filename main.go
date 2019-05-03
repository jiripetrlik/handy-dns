package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jiripetrlik/handy-dns-manager/internal/app/dnszone"
	"github.com/jiripetrlik/handy-dns-manager/internal/app/rest"
)

func main() {
	ipPtr := flag.String("i", "127.0.0.1", "IP of primary nameserver")
	originPtr := flag.String("o", "example-domain.", "Domain origin")
	primaryNameServerPtr := flag.String("p", "ns1", "Primary name server")
	emailPtr := flag.String("e", "email.example-domain.", "Hostmaster email")
	dnszonePtr := flag.String("f", "example-domain.hosts", "Zone file")
	zoneDataPtr := flag.String("d", "example-domain.json", "Data about zone")
	htpasswdPtr := flag.String("s", "", "Htpasswd file")
	certFilePtr := flag.String("certfile", "", "Cert file")
	keyFilePtr := flag.String("keyfile", "", "Key file")
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
		&sync.Mutex{},
	}
	dnsZone.Initialize(*ipPtr, *originPtr, *primaryNameServerPtr, *emailPtr)

	restServer := rest.NewHandyDNSRestServer(&dnsZone, *htpasswdPtr)
	restServer.HandleRestAPI()

	var srv http.Server
	idleConnsClosed := make(chan struct{})
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		log.Print("Registering handler for graceful termination")
		<-sigs
		log.Print("Closing server")

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatal("HTTP server Shutdown: " + err.Error())
		}
		close(idleConnsClosed)
	}()

	var err error
	if len(*certFilePtr) > 0 && len(*keyFilePtr) > 0 {
		srv.Addr = ":8443"
		err = srv.ListenAndServeTLS(*certFilePtr, *keyFilePtr)
	} else {
		srv.Addr = ":8080"
		err = srv.ListenAndServe()
	}
	if err != http.ErrServerClosed {
		log.Fatal("HTTP server ListenAndServe: " + err.Error())
	}

	<-idleConnsClosed
}
