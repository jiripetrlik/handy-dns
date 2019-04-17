package dnszone

type DNSZoneData struct {
	Origin            string
	PrimaryNameServer string
	HostmasterEmail   string
	SerialNumber      int64
	Refresh           int64
	Retry             int64
	Expire            int64
	TTL               int64
	ItemIndex         int64
	ZoneItems         []ZoneItem
}

type ZoneItem struct {
	ID       int64
	Name     string
	Class    string
	ItemType string
	Data     string
}

func NewDNSZoneData(origin string, primaryNameServer string, hostmasterEmail string) DNSZoneData {
	zoneData := DNSZoneData{
		Origin:            origin,
		PrimaryNameServer: primaryNameServer,
		HostmasterEmail:   hostmasterEmail,
		SerialNumber:      1,
		Refresh:           3600,
		Retry:             600,
		Expire:            604800,
		TTL:               1800,
		ItemIndex:         1,
		ZoneItems:         make([]ZoneItem, 0),
	}

	return zoneData
}
