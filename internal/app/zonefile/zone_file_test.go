package zonefile

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func tmpFile(prefix string) string {
	file, _ := ioutil.TempFile("", prefix)
	file.Close()

	return file.Name()
}

func createTestDNSZone() DNSZone {
	var dnsZone DNSZone
	dnsZone.ZoneName = "test-domain"
	dnsZone.ZoneFile = tmpFile("test-zone-file-")
	dnsZone.ZoneData = tmpFile("test-zone-data-")

	return dnsZone
}

func deleteTestDNSZone(zone *DNSZone) {
	os.Remove(zone.ZoneFile)
	os.Remove(zone.ZoneData)
}

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

func TestWriteReadZoneFile(t *testing.T) {
	dnsZone := createTestDNSZone()
	defer deleteTestDNSZone(&dnsZone)

	itemsList := createTestItemsList()
	dnsZone.WriteZoneFile(itemsList)
	itemsList2 := dnsZone.ReadZoneFile()
	if len(itemsList2) != 2 {
		t.Error("Loaded items list should have 2 items")
	}
}

func TestExportZoneFile(t *testing.T) {
	dnsZone := createTestDNSZone()
	defer deleteTestDNSZone(&dnsZone)

	itemsList := createTestItemsList()
	dnsZone.WriteZoneFile(itemsList)
	dnsZone.ExportZoneFile()

	data, _ := ioutil.ReadFile(dnsZone.ZoneFile)
	content := string(data)
	if strings.Contains(content, "text1") == false || strings.Contains(content, "text2") == false {
		t.Error("Zone file does not contain all required strings(text1, text2)")
	}
}
