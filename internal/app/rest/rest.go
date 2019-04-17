package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jiripetrlik/handy-dns/internal/app/dnszone"
)

type HandyDnsRestServer struct {
	DNSZone *dnszone.DNSZone
}

func (s *HandyDnsRestServer) endpointListItems(writer http.ResponseWriter, request *http.Request) {
	itemsList := s.DNSZone.GetZoneData().ZoneItems
	itemsListJSON, _ := json.Marshal(itemsList)

	writer.Header().Set("Content-Type", "text/json")
	io.WriteString(writer, string(itemsListJSON))
	log.Printf("List with %v zone items was send", len(itemsList))
}

func (s *HandyDnsRestServer) endpointCreateItem(writer http.ResponseWriter, request *http.Request) {
	item := dnszone.ZoneItem{
		ID:       0,
		Name:     request.URL.Query().Get("name"),
		Class:    request.URL.Query().Get("class"),
		ItemType: request.URL.Query().Get("itemType"),
		Data:     request.URL.Query().Get("data"),
	}

	s.DNSZone.AddZoneItem(item)

	log.Printf("Item %v was created with id %v", item, item.ID)
}

func (s *HandyDnsRestServer) endpointUpdateItem(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
	item := dnszone.ZoneItem{
		ID:       id,
		Name:     request.URL.Query().Get("name"),
		Class:    request.URL.Query().Get("class"),
		ItemType: request.URL.Query().Get("itemType"),
		Data:     request.URL.Query().Get("data"),
	}

	s.DNSZone.UpdateZoneItem(item)

	jsonItem, _ := json.Marshal(item)
	writer.Header().Set("Content-Type", "text/json")
	io.WriteString(writer, string(jsonItem))

	log.Printf("Item %v was updated", item)
}

func (s *HandyDnsRestServer) endpointDeleteItem(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)

	s.DNSZone.DeleteZoneItem(id)

	log.Printf("Item %v was deleted", id)
}

func (s *HandyDnsRestServer) HandleRestAPI() {
	http.HandleFunc("/api/list", s.endpointListItems)
	http.HandleFunc("/api/create", s.endpointCreateItem)
	http.HandleFunc("/api/update", s.endpointUpdateItem)
	http.HandleFunc("/api/delete", s.endpointDeleteItem)
}
