package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestWriteReadZoneFile(t *testing.T) {
	file, _ := ioutil.TempFile("", "test-zone-file-")
	defer os.Remove(file.Name())
	file.Close()

	itemsList := createTestItemsList()
	writeZoneFile(itemsList, file.Name())
	itemsList2 := readZoneFile(file.Name())
	if len(itemsList2) != 2 {
		t.Error("Loaded items list should have 2 items")
	}
}

func TestExportZoneFile(t *testing.T) {
	file, _ := ioutil.TempFile("", "test-export-zone-file-")
	file.Close()
	defer os.Remove(file.Name())

	itemsList := createTestItemsList()
	exportZoneFile(itemsList, file.Name())

	data, _ := ioutil.ReadFile(file.Name())
	content := string(data)
	if strings.Contains(content, "text1") == false || strings.Contains(content, "text2") == false {
		t.Error("Zone file does not contain all required strings(text1, text2)")
	}
}
