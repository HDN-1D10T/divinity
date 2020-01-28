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

// Handler for TCP
// Parses config options and handles as necessary
func Handler(lines []string) {
	if len(*Conf.DumpList) > 0 || len(*Conf.List) > 0 {
		doList(lines)
		return
	}
}

func doList(lines []string) {
	dumplistMatch := regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}(:[0-9]{1,5})?.*:?.*`)
	for _, line := range lines {
		go func(line string) {
			if !dumplistMatch.MatchString(line) {
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
			ip, port := func(hostString string) (string, string) {
				host := strings.Split(hostString, ":")
				if len(hostString) == 2 {
					ip := host[0]
					port := host[1]
					return ip, port
				}
				ip := host[0]
				return ip, ""
			}(hostString)
			user, pass := func(credString string) (string, string) {
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
					return "", ""
				}
				creds = strings.Split(credString, ":")
				if len(credString) > 1 {
					if nouserRE.MatchString(creds[0]) {
						user := ""
						pass := creds[0]
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
				return "", ""
			}(credString)
			if len(Port) > 0 {
				TelnetPreflight(hostString, ip, Port, user, pass, Alert, OutputFile)
				return
			}
			TelnetPreflight(hostString, ip, port, user, pass, Alert, OutputFile)
			return
		}(line)
		time.Sleep(20 * time.Millisecond)
	}
}

/*
func doIPList(lines []string) {
	listMatch := regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}(:[0-9]{1,5})?(\s+|\t+)?$`)
	for _, line := range lines {
		go func(line string) {
			if !listMatch.MatchString(line) {
				log.Println("string formatted incorrectly " + line)
				return
			}
			connectionString := strings.Split(line, ":")
			hostString, portString := func(connectionString []string) (string, string) {
				if len(connectionString) > 1 {
					ip := connectionString[0]
					port := connectionString[1]
					return ip, port
				}
				ip := connectionString[0]
				return ip, ""
			}(connectionString)
			ip := strings.Replace(hostString, " ", "", -1)
			port := strings.Replace(portString, " ", "", -1)
			port = func(port string) string {
				if len(*Conf.Port) > 0 {
					return port
				}
				return *Conf.Port
			}(port)
			credString := *Conf.Credentials
			user, pass := func(credString string) (string, string) {
				if len(Username) > 0 || len(Password) > 0 {
					user := Username
					pass := Password
					return user, pass
				}
				if len(*Conf.Credentials) > 0 {
					creds := strings.Split(credString, ":")
					if len(creds) > 0 {
						user := creds[0]
						pass := creds[1]
						return user, pass
					}
					if nouserRE.MatchString(credString) {
						user := ""
						pass := credString
						return user, pass
					}
					if nopassRE.MatchString(credString) {
						user := credString
						pass := ""
						return user, pass
					}
				}
				user := ""
				pass := ""
				return user, pass
			}(credString)
			if len(Port) > 0 {
				TelnetPreflight(line, ip, Port, user, pass, Alert, OutputFile)
				return
			}
			TelnetPreflight(line, ip, port, user, pass, Alert, OutputFile)
			return
		}(line)
		time.Sleep(400 * time.Millisecond)
	}
}
*/
