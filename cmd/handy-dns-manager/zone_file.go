package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ZoneItem struct {
	ID       int64
	Name     string
	Class    string
	ItemType string
	Data     string
}

func readZoneFile(fileName string) []ZoneItem {
	var zoneItems []ZoneItem
	data, _ := ioutil.ReadFile(fileName)
	json.Unmarshal(data, &zoneItems)

	return zoneItems
}

func writeZoneFile(zoneItems []ZoneItem, fileName string) {
	data, _ := json.Marshal(zoneItems)
	ioutil.WriteFile(fileName, data, 0664)
}

func exportZoneFile(zoneItems []ZoneItem, fileName string) {
	file, _ := os.OpenFile(
		fileName,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0664,
	)
	defer file.Close()

	for _, item := range zoneItems {
		fmt.Fprintf(
			file,
			"%v\t%v\t%v\t%v\n",
			item.Name,
			item.Class,
			item.ItemType,
			item.Data,
		)
	}
}
