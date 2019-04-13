package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jiripetrlik/handy-dns/internal/app/zonefile"
)

type HandyDnsRestServer struct {
	DNSZone *zonefile.DNSZone
}

func (s *HandyDnsRestServer) endpointListItems(writer http.ResponseWriter, request *http.Request) {
	itemsList := s.DNSZone.ReadZoneFile()
	itemsListJSON, _ := json.Marshal(itemsList)

	writer.Header().Set("Content-Type", "text/json")
	io.WriteString(writer, string(itemsListJSON))
	log.Printf("List with %v zone items was send", len(itemsList))
}

func (s *HandyDnsRestServer) endpointCreateItem(writer http.ResponseWriter, request *http.Request) {
	itemsList := s.DNSZone.ReadZoneFile()

	var item zonefile.ZoneItem
	item.ID = nextID(itemsList)
	item.Name = request.URL.Query().Get("name")
	item.Class = request.URL.Query().Get("class")
	item.ItemType = request.URL.Query().Get("itemType")
	item.Data = request.URL.Query().Get("data")

	itemsList = append(itemsList, item)
	s.DNSZone.WriteZoneFile(itemsList)
	s.DNSZone.ExportZoneFile()

	jsonID, _ := json.Marshal(item.ID)
	writer.Header().Set("Content-Type", "text/json")
	io.WriteString(writer, string(jsonID))

	log.Printf("Item %v was created with id %v", item, item.ID)
}

func (s *HandyDnsRestServer) endpointUpdateItem(writer http.ResponseWriter, request *http.Request) {
	itemsList := s.DNSZone.ReadZoneFile()
	id, _ := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
	_, item := findItem(id, itemsList)

	item.ID, _ = strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
	item.Name = request.URL.Query().Get("name")
	item.Class = request.URL.Query().Get("class")
	item.ItemType = request.URL.Query().Get("itemType")
	item.Data = request.URL.Query().Get("data")

	s.DNSZone.WriteZoneFile(itemsList)
	s.DNSZone.ExportZoneFile()

	jsonItem, _ := json.Marshal(item)
	writer.Header().Set("Content-Type", "text/json")
	io.WriteString(writer, string(jsonItem))

	log.Printf("Item %v was updated", item)
}

func (s *HandyDnsRestServer) endpointDeleteItem(writer http.ResponseWriter, request *http.Request) {
	itemsList := s.DNSZone.ReadZoneFile()
	id, _ := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
	index, _ := findItem(id, itemsList)

	itemsList = append(itemsList[:index], itemsList[index+1:]...)
	s.DNSZone.WriteZoneFile(itemsList)
	s.DNSZone.ExportZoneFile()

	jsonID, _ := json.Marshal(id)
	writer.Header().Set("Content-Type", "text/json")
	io.WriteString(writer, string(jsonID))

	log.Printf("Item %v was deleted", id)
}

func (s *HandyDnsRestServer) HandleRestAPI() {
	http.HandleFunc("/api/list", s.endpointListItems)
	http.HandleFunc("/api/create", s.endpointCreateItem)
	http.HandleFunc("/api/update", s.endpointUpdateItem)
	http.HandleFunc("/api/delete", s.endpointDeleteItem)
}

func nextID(zoneItems []zonefile.ZoneItem) int64 {
	var max int64 = 0

	for _, item := range zoneItems {
		if item.ID > max {
			max = item.ID
		}
	}

	return max + 1
}

func findItem(id int64, zoneItems []zonefile.ZoneItem) (int, *zonefile.ZoneItem) {
	var index int
	var foundItem *zonefile.ZoneItem

	for i, item := range zoneItems {
		if item.ID == id {
			index = i
			foundItem = &zoneItems[i]
		}
	}

	return index, foundItem
}
