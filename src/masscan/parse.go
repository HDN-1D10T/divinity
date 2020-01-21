package masscan

import "encoding/xml"

// Address struct
type Address struct {
	Addr     string `xml:"addr,attr"`
	AddrType string `xml:"addrtype,attr"`
}

// State struct
type State struct {
	State     string `xml:"state,attr"`
	Reason    string `xml:"reason,attr"`
	ReasonTTL string `xml:"reason_ttl,attr"`
}

// Host struct
type Host struct {
	XMLName xml.Name `xml:"host"`
	Endtime string   `xml:"endtime,attr"`
	Address Address  `xml:"address"`
	Ports   Ports    `xml:"ports>port"`
}

// Ports struct
type Ports []struct {
	Protocol string  `xml:"protocol,attr"`
	Portid   string  `xml:"portid,attr"`
	State    State   `xml:"state"`
	Service  Service `xml:"service"`
}

// Service struct
type Service struct {
	Name   string `xml:"name,attr"`
	Banner string `xml:"banner,attr"`
}
