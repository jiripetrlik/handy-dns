package dnszone

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type DNSZone struct {
	ZoneFile     string
	ZoneDataFile string
	Mutex        *sync.Mutex
}

func (z *DNSZone) Initialize(ip string, origin string, primaryNameServer string, hostmasterEmail string) {
	z.Mutex.Lock()
	defer z.Mutex.Unlock()

	_, err := os.Stat(z.ZoneDataFile)
	if os.IsNotExist(err) {
		dnsZoneDataDNSZoneData := NewDNSZoneData(ip, origin, primaryNameServer, hostmasterEmail)
		z.writeZoneData(dnsZoneDataDNSZoneData)
	}
	z.writeZoneFile()
}

func (z *DNSZone) GetZoneData() DNSZoneData {
	z.Mutex.Lock()
	defer z.Mutex.Unlock()

	dnsZoneData := z.readZoneData()

	return dnsZoneData
}

func (z *DNSZone) AddZoneItem(item ZoneItem) int64 {
	z.Mutex.Lock()
	defer z.Mutex.Unlock()

	dnsZoneData := z.readZoneData()
	item.ID = dnsZoneData.ItemIndex
	dnsZoneData.ItemIndex++
	dnsZoneData.ZoneItems = append(dnsZoneData.ZoneItems, item)
	dnsZoneData.SerialNumber++
	z.writeZoneData(dnsZoneData)
	z.writeZoneFile()

	return item.ID
}

func (z *DNSZone) UpdateZoneItem(item ZoneItem) error {
	z.Mutex.Lock()
	defer z.Mutex.Unlock()

	dnsZoneData := z.readZoneData()

	_, oldItem, err := findItem(item.ID, dnsZoneData.ZoneItems)
	if err != nil {
		return fmt.Errorf("Item %v does not exist", item.ID)
	}
	oldItem.Name = item.Name
	oldItem.Class = item.Class
	oldItem.ItemType = item.ItemType
	oldItem.Data = item.Data

	dnsZoneData.SerialNumber++

	z.writeZoneData(dnsZoneData)
	z.writeZoneFile()

	return nil
}

func (z *DNSZone) DeleteZoneItem(id int64) error {
	z.Mutex.Lock()
	defer z.Mutex.Unlock()

	dnsZoneData := z.readZoneData()

	index, _, err := findItem(id, dnsZoneData.ZoneItems)
	if err != nil {
		return fmt.Errorf("Item %v does not exist", id)
	}
	dnsZoneData.ZoneItems = append(dnsZoneData.ZoneItems[:index], dnsZoneData.ZoneItems[index+1:]...)

	dnsZoneData.SerialNumber++

	z.writeZoneData(dnsZoneData)
	z.writeZoneFile()

	return nil
}

func (z *DNSZone) readZoneData() DNSZoneData {
	var dnsZoneData DNSZoneData
	data, err := ioutil.ReadFile(z.ZoneDataFile)
	if err != nil {
		log.Fatal("Unable to read zone data. Caused by " + err.Error())
	}
	err = json.Unmarshal(data, &dnsZoneData)
	if err != nil {
		log.Fatal("Unable to unmarshall zone data. Caused by " + err.Error())
	}

	return dnsZoneData
}

func (z *DNSZone) writeZoneData(dnsZoneData DNSZoneData) {
	data, err := json.MarshalIndent(dnsZoneData, "", "\t")
	if err != nil {
		log.Fatal("Unable to marshall zone data. Caused by " + err.Error())
	}
	err = ioutil.WriteFile(z.ZoneDataFile, data, 0664)
	if err != nil {
		log.Fatal("Unable to write zone data. Caused by " + err.Error())
	}
}

func (z *DNSZone) writeZoneFile() {
	zoneData := z.readZoneData()

	file, err := os.OpenFile(
		z.ZoneFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0664,
	)
	if err != nil {
		log.Fatal("Unable to write zone file. Caused by " + err.Error())
	}
	defer file.Close()

	file.WriteString(zoneData.zoneFileHeader())
	for _, item := range zoneData.ZoneItems {
		_, err = fmt.Fprintf(
			file,
			"%v\t%v\t%v\t%v\n",
			item.Name,
			item.Class,
			item.ItemType,
			item.Data,
		)
		if err != nil {
			log.Fatal("Unable to write zone file. Caused by " + err.Error())
		}
	}
}

func findItem(id int64, zoneItems []ZoneItem) (int, *ZoneItem, error) {
	for i, item := range zoneItems {
		if item.ID == id {
			return i, &zoneItems[i], nil
		}
	}

	return 0, nil, fmt.Errorf("Can not find item with id: %v", id)
}
