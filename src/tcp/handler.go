package tcp

import (
	"github.com/HDN-1D10T/divinity/src/config"
)

// Configuration imported from src/config
type Configuration struct{ config.Options }

// Handler for TCP
// Parses config options and handles as necessary
func Handler(ips []string) {
	conf := Configuration{
		config.ParseConfiguration(),
	}
	if *conf.Protocol == "tcp" {
		if *conf.Port == "23" {
			if len(*conf.DumpList) > 1 {
				Telnet(ips, &conf)
			}
		}
	}
}
