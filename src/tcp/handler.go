package tcp

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/HDN-1D10T/divinity/src/config"
	"github.com/HDN-1D10T/divinity/src/util"
)

// Configuration imported from src/config
type Configuration struct{ config.Options }

type IPinfo struct {
	hostString string
	ip         string
	port       string
	user       string
	pass       string
	alert      string
}

const timeout = 120 * time.Millisecond

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
	userRE   = regexp.MustCompile(`.*([Ll]ogin)|([Uu]sername).*`)
	passRE   = regexp.MustCompile(".*[Pp]assword.*")
	promptRE = regexp.MustCompile(`.*[#\$>].*`)
	badRE    = regexp.MustCompile(`.*(Using username)|([Pp]assword:)|([Dd]enied)|([Ii]ncorrect).*`)
)

var wg sync.WaitGroup

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
	if len(creds) > 0 {
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

func doList(ipinfo chan IPinfo, lines []string) {
	runtime.GOMAXPROCS(100)
	listMatch := regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}(:[0-9]{1,5})?.*:?.*`)
	for _, line := range lines {
		if !listMatch.MatchString(line) {
			log.Println("string formatted incorrectly: " + line)
			return
		}
		connectionString := strings.Split(line, " ")
		hostString, credString := func(connectionString []string) (string, string) {
			if len(connectionString) > 1 {
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
		info := IPinfo{
			hostString: hostString,
			ip:         ip,
			port:       port,
			user:       user,
			pass:       pass,
			alert:      Alert,
		}
		ipinfo <- info
		time.Sleep(time.Millisecond)
	}
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// Handler for TCP
// Parses config options and handles as necessary
func Handler(lines []string) {
	if *Conf.OutputFile == "" {
		fmt.Println("You MUST specify an output file via '-out [/path/to/file]'")
		return
	}
	totalTimeouts := time.Duration(*Conf.Timeout) * time.Millisecond
	timeDuration := time.Duration(len(lines)) * totalTimeouts
	timeStart := time.Now()
	timeComplete := timeStart.Add(timeDuration)
	alertMatch := regexp.MustCompile(".*" + Alert + ".*")
	statusMatch := regexp.MustCompile("^STATUS:.*")
	successMatch := regexp.MustCompile("^@SUCCESS@.*")
	if len(*Conf.List) > 0 || len(*Conf.Cidr) > 0 {
		ipInfo := make(chan IPinfo, 0)
		// get all of the ip info:
		go func() {
			doList(ipInfo, lines)
			close(ipInfo)
		}()
		status := make(chan string, 0)
		messages := make(chan string, 0)
		exit := make(chan struct{})
		if *Conf.SSH || *Conf.Port == "22" {
			//SSHPreflight(hostString, ip, port, user, pass, Alert, OutputFile)
			go func() {
				SSHPreflight(timeStart, timeDuration, timeComplete, status, messages, ipInfo)
				close(status)
				close(messages)
			}()
			for {
				select {
				case msg, open := <-status:
					if !open {
						close(exit)
						return
					}
					if msg == "CLEARSCREEN" {
						clearScreen()
					}
					if statusMatch.MatchString(msg) {
						msg = strings.Split(msg, "STATUS:")[1]
						fmt.Print(msg)
					} else {
						fmt.Print(msg)
					}
				case msg, open := <-messages:
					if !open {
						close(exit)
						return
					}
					if successMatch.MatchString(msg) {
						msg = strings.Split(msg, "@SUCCESS@")[1]
						if alertMatch.MatchString(msg) {
							util.FileWrite(msg)
						}
						fmt.Print(msg)
					} else if msg != "nil" {
						fmt.Print(msg)
					}
				}
			}
		}
		if *Conf.Telnet || *Conf.Port == "23" {
			ipInfo = make(chan IPinfo, 0)
			go func() {
				doList(ipInfo, lines)
				close(ipInfo)
			}()
			// TODO: Implement concurrency in Telnet (it is already pretty quick and accurate with -timeout set to 120ms)
			for info := range ipInfo {
				hostString := info.hostString
				ip := info.ip
				port := info.port
				user := info.user
				pass := info.pass
				alert := info.alert
				TelnetPreflight(hostString, ip, port, user, pass, alert, OutputFile)
			}
		}
		return
	}
}
