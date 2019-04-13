package zonefile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type DNSZone struct {
	ZoneName string
	ZoneFile string
	ZoneData string
}

type ZoneItem struct {
	ID       int64
	Name     string
	Class    string
	ItemType string
	Data     string
}

func (z *DNSZone) Initialize() {
	_, err := os.Stat(z.ZoneData)
	if os.IsNotExist(err) {
		zoneItems := []ZoneItem{}
		z.WriteZoneFile(zoneItems)
	}
	z.ExportZoneFile()
}

func (z *DNSZone) ReadZoneFile() []ZoneItem {
	var zoneItems []ZoneItem
	data, _ := ioutil.ReadFile(z.ZoneData)
	json.Unmarshal(data, &zoneItems)

	return zoneItems
}

func (z *DNSZone) WriteZoneFile(zoneItems []ZoneItem) {
	data, _ := json.Marshal(zoneItems)
	ioutil.WriteFile(z.ZoneData, data, 0664)
}

func (z *DNSZone) ExportZoneFile() {
	zoneItems := z.ReadZoneFile()

	file, _ := os.OpenFile(
		z.ZoneFile,
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
