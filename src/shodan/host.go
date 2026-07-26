package shodan

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type HostLocation struct {
	City         string  `json:"city"`
	RegionCode   string  `json:"region_code"`
	AreaCode     int     `json:"area_code"`
	Longitude    float32 `json:"longitude"`
	CountryCode3 string  `json:"country_code3"`
	CountryName  string  `json:"country_name"`
	PostalCode   string  `json:"postal_code"`
	DMACode      int     `json:"dma_code"`
	CountryCode  string  `json:"country_code"`
	Latitude     float32 `json:"latitude"`
}

type Host struct {
	OS        string       `json:"os"`
	Timestamp string       `json:"timestamp"`
	ISP       string       `json:"isp"`
	ASN       string       `json:"asn"`
	Hostnames []string     `json:"hostnames"`
	Location  HostLocation `json:"location"`
	IP        int64        `json:"ip"`
	Domains   []string     `json:"domains"`
	Org       string       `json:"org"`
	Data      string       `json:"data"`
	Port      int          `json:"port"`
	IPString  string       `json:"ip_str"`
}

type HostSearch struct {
	Matches []Host `json:"matches"`
}

func (s *Client) HostSearch(q string) (*HostSearch, error) {
	return s.HostSearchPage(q, 1)
}

func (s *Client) HostSearchPage(q string, page int) (*HostSearch, error) {
	res, err := http.Get(s.hostSearchURL(q, page))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var ret HostSearch
	if err := decodeResponse("Shodan host search", res, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (s *Client) hostSearchURL(q string, page int) string {
	params := url.Values{}
	params.Set("key", s.apiKey)
	params.Set("query", q)
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	return fmt.Sprintf("%s/shodan/host/search?%s", s.baseURL, params.Encode())
}
