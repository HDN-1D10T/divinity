package tcp

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/HDN-1D10T/divinity/src/config"
)

// Configuration imported from src/config
type Configuration struct{ config.Options }

var (
	// Conf - Gets configuration values
	Conf = Configuration{config.ParseConfiguration()}
	// Alert ...
	Alert = *Conf.Alert
	// OutputFile ...
	OutputFile = *Conf.OutputFile
	// Protocol ...
	Protocol = *Conf.Protocol
	// Port ...
	Port = *Conf.Port
	// Username ...
	Username = *Conf.Username
	// Password ...
	Password = *Conf.Password
)

var (
	nouserRE = regexp.MustCompile(`^:.+`)
	nopassRE = regexp.MustCompile(`.+:$`)
)

// GetCreds returns username string and password string
func GetCreds(credString string) (string, string) {
	if len(Username) > 0 || len(Password) > 0 {
		user := Username
		pass := Password
		return user, pass
	}
	creds := strings.Split(*Conf.Credentials, ":")
	if len(*Conf.Credentials) > 0 {
		if len(creds) > 1 {
			if nouserRE.MatchString(creds[0]) {
				user := ""
				pass := creds[1]
				return user, pass
			}
			if nopassRE.MatchString(creds[1]) {
				user := creds[0]
				pass := ""
				return user, pass
			}
			user := creds[0]
			pass := creds[1]
			return user, pass
		}
	}
	creds = strings.Split(credString, ":")
	if len(*Conf.Credentials) > 0 {
		if len(creds) > 1 {
			if nouserRE.MatchString(creds[0]) {
				user := ""
				pass := creds[1]
				return user, pass
			}
			if nopassRE.MatchString(creds[1]) {
				user := creds[0]
				pass := ""
				return user, pass
			}
			user := creds[0]
			pass := creds[1]
			return user, pass
		}
	}
	return "", ""
}

// GetIPPort takes a 'ip:port' string and returns the ip and port
func GetIPPort(connectionString string) (string, string) {
	hostString := strings.Split(connectionString, ":")
	if len(hostString) == 2 {
		ip := hostString[0]
		port := hostString[1]
		if len(Port) > 0 {
			return ip, Port
		}
		return ip, port
	}
	ip := hostString[0]
	if len(Port) > 0 {
		return ip, Port
	}
	return ip, ""
}

// Handler for TCP
// Parses config options and handles as necessary
func Handler(lines []string) {
	if len(*Conf.List) > 0 || len(*Conf.Cidr) > 0 {
		doList(lines)
		return
	}
}

func doList(lines []string) {
	listMatch := regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}(:[0-9]{1,5})?.*:?.*`)
	for _, line := range lines {
		go func(line string) {
			if !listMatch.MatchString(line) {
				log.Println("string formatted incorrectly: " + line)
				return
			}
			connectionString := strings.Split(line, " ")
			hostString, credString := func(connectionString []string) (string, string) {
				if len(connectionString) == 2 {
					fmt.Println(len(connectionString))
					hostString := connectionString[0]
					credString := connectionString[1]
					hostString = strings.Replace(hostString, " ", "", -1)
					credString = strings.Replace(credString, " ", "", -1)
					return hostString, credString
				}
				hostString := connectionString[0]
				return hostString, ""
			}(connectionString)
			ip, port := GetIPPort(hostString)
			user, pass := GetCreds(credString)
			if *Conf.Telnet || port == "23" {
				TelnetPreflight(hostString, ip, port, user, pass, Alert, OutputFile)
			}
			return
		}(line)
		time.Sleep(50 * time.Millisecond)
	}
}
