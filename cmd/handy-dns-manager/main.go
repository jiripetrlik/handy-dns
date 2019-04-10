package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

func createTestItemsList() []ZoneItem {
	item1 := ZoneItem{
		1,
		"text1",
		"IN",
		"NS",
		"127.0.0.1"}
	item2 := ZoneItem{
		2,
		"text2",
		"IN",
		"A",
		"127.0.0.2"}
	items := []ZoneItem{item1, item2}

	return items
}

func endpointListItems(writer http.ResponseWriter, request *http.Request) {
	itemsList := createTestItemsList()
	itemsListJSON, _ := json.Marshal(itemsList)

	writer.Header().Set("Content-Type", "text/json")
	io.WriteString(writer, string(itemsListJSON))
	log.Printf("List with %v zone items was send", len(itemsList))
}

func endpointCreateItem(writer http.ResponseWriter, request *http.Request) {
	var item ZoneItem

	item.Name = request.URL.Query().Get("name")
	item.Class = request.URL.Query().Get("class")
	item.ItemType = request.URL.Query().Get("itemType")
	item.Data = request.URL.Query().Get("data")

	id := 5
	jsonID, _ := json.Marshal(id)
	writer.Header().Set("Content-Type", "text/json")
	io.WriteString(writer, string(jsonID))

	log.Printf("Item %v was created with id %v", item, id)
}

func endpointUpdateItem(writer http.ResponseWriter, request *http.Request) {
	var item ZoneItem

	item.ID, _ = strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
	item.Name = request.URL.Query().Get("name")
	item.Class = request.URL.Query().Get("class")
	item.ItemType = request.URL.Query().Get("itemType")
	item.Data = request.URL.Query().Get("data")

	jsonItem, _ := json.Marshal(item)
	writer.Header().Set("Content-Type", "text/json")
	io.WriteString(writer, string(jsonItem))

	log.Printf("Item %v was updated", item)
}

func endpointDeleteItem(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)

	jsonID, _ := json.Marshal(id)
	writer.Header().Set("Content-Type", "text/json")
	io.WriteString(writer, string(jsonID))

	log.Printf("Item %v was deleted", id)
}

func main() {
	http.HandleFunc("/api/list", endpointListItems)
	http.HandleFunc("/api/create", endpointCreateItem)
	http.HandleFunc("/api/update", endpointUpdateItem)
	http.HandleFunc("/api/delete", endpointDeleteItem)
	http.ListenAndServe(":8080", nil)
}
