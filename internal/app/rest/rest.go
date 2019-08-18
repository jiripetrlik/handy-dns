package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jiripetrlik/handy-dns-manager/internal/app/dnszone"
	auth "github.com/abbot/go-http-auth"
	"github.com/rakyll/statik/fs"
	_ "github.com/jiripetrlik/handy-dns-manager/statik"
)

type HandyDnsRestServer struct {
	DNSZone *dnszone.DNSZone
	authenticator *auth.BasicAuth
}

type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func NewHandyDNSRestServer(dnsZone *dnszone.DNSZone, htpasswd string) *HandyDnsRestServer {
	var authenticator *auth.BasicAuth
	if len(htpasswd) > 0 {
		htpasswd := auth.HtpasswdFileProvider(htpasswd)
		authenticator = auth.NewBasicAuthenticator("Handy DNS Realm", htpasswd)
	}

	server := HandyDnsRestServer{
		DNSZone: dnsZone,
		authenticator: authenticator,
	}

	return &server
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err != nil {
		log.Printf("Error processing http request. Caused by %v", err.Error)
		http.Error(w, err.Message, err.Code)
	}
}

func (s *HandyDnsRestServer) isAllowed(request *http.Request) bool {
	if s.authenticator == nil {
		return true
	}

	if len(s.authenticator.CheckAuth(request)) > 0 {
		return true
	}

	return false
}

func denyAccess(writer http.ResponseWriter) {
	writer.Header().Set("WWW-Authenticate", "Basic")
	writer.WriteHeader(http.StatusUnauthorized)
	writer.Write([]byte("Unauthorized"))
}

func (s *HandyDnsRestServer) endpointListItems(writer http.ResponseWriter, request *http.Request) *appError {
	if s.isAllowed(request) == false {
		denyAccess(writer)
		return nil
	}

	itemsList := s.DNSZone.GetZoneData().ZoneItems
	itemsListJSON, err := json.MarshalIndent(itemsList, "", "\t")
	if err != nil {
		return &appError{
			Error:   err,
			Message: "Can not marshal zone items. Caused by " + err.Error(),
			Code:    500,
		}
	}

	writer.Header().Set("Content-Type", "text/json")
	io.WriteString(writer, string(itemsListJSON))
	log.Printf("List with %v zone items was send", len(itemsList))
	return nil
}

func (s *HandyDnsRestServer) endpointCreateItem(writer http.ResponseWriter, request *http.Request) *appError {
	if s.isAllowed(request) == false {
		denyAccess(writer)
		return nil
	}

	item := dnszone.ZoneItem{
		ID:       0,
		Name:     request.URL.Query().Get("name"),
		Class:    request.URL.Query().Get("class"),
		ItemType: request.URL.Query().Get("itemType"),
		Data:     request.URL.Query().Get("data"),
	}

	id := s.DNSZone.AddZoneItem(item)
	idJSON, err := json.MarshalIndent(id, "", "\t")
	if err != nil {
		return &appError{
			Error:   err,
			Message: "Can not marshal item id. Caused by " + err.Error(),
			Code:    500,
		}
	}
	writer.Header().Set("Content-Type", "text/json")
	io.WriteString(writer, string(idJSON))
	log.Printf("Item %v was created with id %v", item, id)
	return nil
}

func (s *HandyDnsRestServer) endpointUpdateItem(writer http.ResponseWriter, request *http.Request) *appError {
	if s.isAllowed(request) == false {
		denyAccess(writer)
		return nil
	}

	id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
	if err != nil {
		return &appError{
			Error:   err,
			Message: "Unable to convert id to int64",
			Code:    404,
		}
	}
	item := dnszone.ZoneItem{
		ID:       id,
		Name:     request.URL.Query().Get("name"),
		Class:    request.URL.Query().Get("class"),
		ItemType: request.URL.Query().Get("itemType"),
		Data:     request.URL.Query().Get("data"),
	}

	err = s.DNSZone.UpdateZoneItem(item)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return &appError{
				Error:   err,
				Message: "Item does not exist",
				Code:    404,
			}
		}
		return &appError{
			Error:   err,
			Message: "Error updating item. Caused by " + err.Error(),
			Code:    500,
		}
	}

	jsonItem, err := json.MarshalIndent(item, "", "\t")
	if err != nil {
		return &appError{
			Error:   err,
			Message: "Can not marshal item. Caused by " + err.Error(),
			Code:    500,
		}
	}
	writer.Header().Set("Content-Type", "text/json")
	io.WriteString(writer, string(jsonItem))

	log.Printf("Item %v was updated", item)
	return nil
}

func (s *HandyDnsRestServer) endpointDeleteItem(writer http.ResponseWriter, request *http.Request) *appError {
	if s.isAllowed(request) == false {
		denyAccess(writer)
		return nil
	}

	id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
	if err != nil {
		return &appError{
			Error:   err,
			Message: "Unable to convert id to int64",
			Code:    404,
		}
	}

	err = s.DNSZone.DeleteZoneItem(id)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return &appError{
				Error:   err,
				Message: "Item does not exist",
				Code:    404,
			}
		}
		return &appError{
			Error:   err,
			Message: "Error deleting item. Caused by " + err.Error(),
			Code:    500,
		}
	}

	idJSON, err := json.MarshalIndent(id, "", "\t")
	if err != nil {
		return &appError{
			Error:   err,
			Message: "Can not marshal item id. Caused by " + err.Error(),
			Code:    500,
		}
	}
	writer.Header().Set("Content-Type", "text/json")
	io.WriteString(writer, string(idJSON))
	log.Printf("Item %v was deleted", id)
	return nil
}

func (s *HandyDnsRestServer) HandleRestAPI() {
	http.Handle("/api/list", appHandler(s.endpointListItems))
	http.Handle("/api/create", appHandler(s.endpointCreateItem))
	http.Handle("/api/update", appHandler(s.endpointUpdateItem))
	http.Handle("/api/delete", appHandler(s.endpointDeleteItem))

	statikFS, err := fs.New()
	if err != nil {
	    panic(err)
	}
	staticServer := http.FileServer(statikFS)
	sh := http.StripPrefix("/swaggerui/", staticServer)
	http.Handle("/swaggerui/", sh)
}
