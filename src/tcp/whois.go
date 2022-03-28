package tcp

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/HDN-1D10T/divinity/src/util"
)

// GetRoutes returns IPv4 CIDR routes for the queried ASNumber
func GetRoutes(asn string) (string, error) {
	conn, err := net.Dial("tcp", "whois.radb.net:43")
	if err != nil {
		return "", err
	}
	conn.Write([]byte("-i origin " + asn + "\r\n"))
	buf := make([]byte, 1024)
	res := []byte{}
	for {
		numbytes, err := conn.Read(buf)
		sbuf := buf[0:numbytes]
		res = append(res, sbuf...)
		if err != nil {
			break
		}
	}
	conn.Close()
	return string(res), nil
}

// GetAllRoutes calls GetRoute for multiple ASNumbers
func GetAllRoutes(asns []string) []string {
	var ipv4_routes []string
	ipv4_routes_regex := regexp.MustCompile(`^route:\s+[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}/[0-9]{1,2}`)
	ipv4_cidr_regex := regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}/[0-9]{1,2}`)
	for _, asn := range asns {
		res, _ := GetRoutes(asn)
		lines := strings.Split(res, "\n")
		for _, line := range lines {
			route_matches := ipv4_routes_regex.FindAllString(line, -1)
			for _, route := range route_matches {
				if len(route) > 0 {
					cidr_matches := ipv4_cidr_regex.FindAllString(line, -1)
					for _, match := range cidr_matches {
						if len(match) > 1 {
							ipv4_routes = append(ipv4_routes, match)
							if len(*Conf.OutputFile) > 0 {
								util.FileWrite(match + "\n")
							}
							fmt.Println(match)
						}
					}
				}
			}
		}
	}
	return ipv4_routes
}
