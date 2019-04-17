package dnszone

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type DNSZone struct {
	ZoneFile     string
	ZoneDataFile string
}

func (z *DNSZone) Initialize(origin string, primaryNameServer string, hostmasterEmail string) {
	_, err := os.Stat(z.ZoneDataFile)
	if os.IsNotExist(err) {
		dnsZoneDataDNSZoneData := NewDNSZoneData(origin, primaryNameServer, hostmasterEmail)
		z.writeZoneData(dnsZoneDataDNSZoneData)
	}
	z.writeZoneFile()
}

func (z *DNSZone) GetZoneData() DNSZoneData {
	dnsZoneData := z.readZoneData()

	return dnsZoneData
}

func (z *DNSZone) AddZoneItem(item ZoneItem) int64 {
	dnsZoneData := z.readZoneData()
	item.ID = dnsZoneData.ItemIndex
	dnsZoneData.ItemIndex++
	dnsZoneData.ZoneItems = append(dnsZoneData.ZoneItems, item)
	z.writeZoneData(dnsZoneData)
	z.writeZoneFile()

	return item.ID
}

func (z *DNSZone) UpdateZoneItem(item ZoneItem) {
	dnsZoneData := z.readZoneData()

	_, oldItem := findItem(item.ID, dnsZoneData.ZoneItems)
	oldItem.Name = item.Name
	oldItem.Class = item.Class
	oldItem.ItemType = item.ItemType
	oldItem.Data = item.Data

	z.writeZoneData(dnsZoneData)
	z.writeZoneFile()
}

func (z *DNSZone) DeleteZoneItem(id int64) {
	dnsZoneData := z.readZoneData()

	index, _ := findItem(id, dnsZoneData.ZoneItems)
	dnsZoneData.ZoneItems = append(dnsZoneData.ZoneItems[:index], dnsZoneData.ZoneItems[index+1:]...)

	z.writeZoneData(dnsZoneData)
	z.writeZoneFile()
}

func (z *DNSZone) readZoneData() DNSZoneData {
	var dnsZoneData DNSZoneData
	data, _ := ioutil.ReadFile(z.ZoneDataFile)
	json.Unmarshal(data, &dnsZoneData)

	return dnsZoneData
}

func (z *DNSZone) writeZoneData(dnsZoneData DNSZoneData) {
	data, _ := json.Marshal(dnsZoneData)
	ioutil.WriteFile(z.ZoneDataFile, data, 0664)
}

func (z *DNSZone) writeZoneFile() {
	zoneData := z.readZoneData()

	file, _ := os.OpenFile(
		z.ZoneFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0664,
	)
	defer file.Close()

	for _, item := range zoneData.ZoneItems {
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

func findItem(id int64, zoneItems []ZoneItem) (int, *ZoneItem) {
	var index int
	var foundItem *ZoneItem

	for i, item := range zoneItems {
		if item.ID == id {
			index = i
			foundItem = &zoneItems[i]
		}
	}

	return index, foundItem
}
