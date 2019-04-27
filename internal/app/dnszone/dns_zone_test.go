package dnszone

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
)

func tmpFile(prefix string) (string, error) {
	file, err := ioutil.TempFile("", prefix)
	if err != nil {
		return "", fmt.Errorf("Unable to create temporary file %v. Caused by %v", file.Name(), err.Error())
	}
	file.Close()
	err = os.Remove(file.Name())
	if err != nil {
		return "", fmt.Errorf("Unable to delete temporary file %v. Caused by %v", file.Name(), err.Error())
	}

	return file.Name(), nil
}

func createTestDNSZone() (*DNSZone, error) {
	tmpZoneFile, err := tmpFile("test-zone-file-")
	if err != nil {
		return nil, fmt.Errorf("Unable to create test-zone-* tmp file. Caused by %v", err.Error())
	}
	tmpZoneDataFile, err := tmpFile("test-zone-data-")
	if err != nil {
		return nil, fmt.Errorf("Unable to create test-zone-data-* tmp file. Caused by %v", err.Error())
	}

	dnsZone := DNSZone{
		ZoneFile:     tmpZoneFile,
		ZoneDataFile: tmpZoneDataFile,
		Mutex:        &sync.Mutex{},
	}
	dnsZone.Initialize("127.0.0.1", "test-domain.", "ns1", "email.test-domain.")

	return &dnsZone, nil
}

func deleteTestDNSZone(zone *DNSZone) {
	os.Remove(zone.ZoneFile)
	os.Remove(zone.ZoneDataFile)
}

func createTestItemsList() []ZoneItem {
	item1 := ZoneItem{
		0,
		"text1",
		"IN",
		"NS",
		"127.0.0.1"}
	item2 := ZoneItem{
		0,
		"text2",
		"IN",
		"A",
		"127.0.0.2"}
	item3 := ZoneItem{
		0,
		"text3",
		"IN",
		"NS",
		"127.0.0.3"}
	items := []ZoneItem{item1, item2, item3}

	return items
}

func TestAddZoneItem(t *testing.T) {
	dnsZone, err := createTestDNSZone()
	if err != nil {
		t.Errorf("Unable to create test DNSZone. Caused by %v", err.Error())
	}
	defer deleteTestDNSZone(dnsZone)

	testItemsList := createTestItemsList()
	for _, item := range testItemsList {
		dnsZone.AddZoneItem(item)
	}

	loadedItemsList := dnsZone.GetZoneData().ZoneItems
	if len(loadedItemsList) != 3 {
		t.Errorf("Loaded items list should have exactly 3 items, but has %v items", len(loadedItemsList))
	}
	for index := range testItemsList {
		if itemsEqualIgnoreID(testItemsList[index], loadedItemsList[index]) == false {
			t.Errorf("Loaded item %v is different", index)
		}
	}

	data, _ := ioutil.ReadFile(dnsZone.ZoneFile)
	for index, item := range testItemsList {
		if strings.Contains(string(data), item.Name) == false {
			t.Errorf("Zone file does not contain name of item %v (%v)", index, item.Name)
		}
	}
}

func TestUpdateZoneItem(t *testing.T) {
	newText := "new-text"
	dnsZone, err := createTestDNSZone()
	if err != nil {
		t.Errorf("Unable to create test DNSZone. Caused by %v", err.Error())
	}
	defer deleteTestDNSZone(dnsZone)

	testItemsList := createTestItemsList()
	oldText := testItemsList[2].Name
	for _, item := range testItemsList {
		dnsZone.AddZoneItem(item)
	}

	loadedItemsList := dnsZone.GetZoneData().ZoneItems
	loadedItemsList[2].Name = newText
	dnsZone.UpdateZoneItem(loadedItemsList[2])

	loadedItemsList = dnsZone.GetZoneData().ZoneItems
	if loadedItemsList[2].Name != newText {
		t.Error("Updated item does not contain new text")
	}

	data, _ := ioutil.ReadFile(dnsZone.ZoneFile)
	if strings.Contains(string(data), oldText) {
		t.Error("Zone file contains old text")
	}
	if strings.Contains(string(data), newText) == false {
		t.Error("Zone file does not contain new text")
	}
}

func TestDeleteZoneItem(t *testing.T) {
	dnsZone, err := createTestDNSZone()
	if err != nil {
		t.Errorf("Unable to create test DNSZone. Caused by %v", err.Error())
	}
	defer deleteTestDNSZone(dnsZone)

	testItemsList := createTestItemsList()
	for _, item := range testItemsList {
		dnsZone.AddZoneItem(item)
	}

	dnsZone.DeleteZoneItem(2)

	loadedItemsList := dnsZone.GetZoneData().ZoneItems
	if len(loadedItemsList) != 2 {
		t.Errorf("Expected number of items is 2, but was %v", len(loadedItemsList))
	}
}

func itemsEqualIgnoreID(item1 ZoneItem, item2 ZoneItem) bool {
	item1.ID = 0
	item2.ID = 0

	if item1 == item2 {
		return true
	}
	return false
}
