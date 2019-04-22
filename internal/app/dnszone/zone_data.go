package dnszone

import (
	"bytes"
	"text/template"
)

type DNSZoneData struct {
	IP                string
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

func NewDNSZoneData(ip string, origin string, primaryNameServer string, hostmasterEmail string) DNSZoneData {
	zoneData := DNSZoneData{
		IP:                ip,
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

func (data *DNSZoneData) zoneFileHeader() string {
	t, _ := template.New("header").Parse(zoneFileHeaderTemplate)
	var tmpBuffer bytes.Buffer
	t.Execute(&tmpBuffer, *data)

	return tmpBuffer.String()
}

const zoneFileHeaderTemplate = `
$ORIGIN {{.Origin}}
@                      3600 SOA   {{.PrimaryNameServer}}.{{.Origin}} {{.HostmasterEmail}} (
                              {{.SerialNumber}}
                              {{.Refresh}}
                              {{.Retry}}
                              {{.Expire}}
                              {{.TTL}})

	IN	NS	{{.PrimaryNameServer}}.{{.Origin}}
{{.PrimaryNameServer}}	IN	A	{{.IP}}
`
