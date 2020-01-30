package tcp

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/HDN-1D10T/divinity/src/util"
)

const timeout = 120 * time.Millisecond

var (
	userRE   = regexp.MustCompile(`.*([Ll]ogin)|([Uu]sername).*`)
	passRE   = regexp.MustCompile(".*[Pp]assword.*")
	promptRE = regexp.MustCompile(`.*[#\$>].*`)
	badRE    = regexp.MustCompile(`.*[Ii]ncorrect.*`)
)

// TelnetPreflight - checks if we want to use the telnet protocol and on which port
func TelnetPreflight(hostString, ip, port, user, pass, Alert, OutputFile string) {
	if Port == "23" {
		Telnet(ip, Port, user, pass, Alert, OutputFile)
	}
	if len(strings.Split(hostString, ":")) > 1 {
		port = strings.Split(hostString, ":")[1]
		port = strings.Replace(hostString, " ", "", -1)
		Telnet(ip, port, user, pass, Alert, OutputFile)
	}
	if *Conf.Telnet {
		Telnet(ip, port, user, pass, Alert, OutputFile)
	}
}

// Telnet - Check for valid credentials
func Telnet(ip, port, user, pass, alert, outputFile string) {
	log.Printf("Trying %s:%s %s:%s...\n", ip, port, user, pass)
	conn, err := DialTimeout("tcp", ip+":"+port, timeout)
	if err != nil {
		//log.Println(err)
		return
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(time.Second))
	authlogin, err := conn.ReadUntil("login:")
	if err != nil {
		//log.Println(err)
		return
	}
	loginString := string(authlogin)
	if userRE.MatchString(loginString) {
		conn.Write([]byte(user + "\r\n"))
		conn.SetReadDeadline(time.Now().Add(time.Second))
		authpass, err := conn.ReadUntil("Password:")
		if err != nil {
			//log.Println(err)
			return
		}
		passString := string(authpass)
		if passRE.MatchString(passString) {
			conn.Write([]byte(pass + "\r\n"))
			conn.SetReadDeadline(time.Now().Add(time.Second))
		}
		prompt, err := conn.ReadUntil("$", ">", "#")
		if err != nil {
			//log.Println(err)
			return
		}
		promptString := string(prompt)
		if promptRE.MatchString(promptString) {
			if !badRE.MatchString(promptString) {
				msg := fmt.Sprintf("%s:%s %s:%s %s", ip, port, user, pass, alert)
				util.LogWrite(msg)
			}
		}
	}
}
