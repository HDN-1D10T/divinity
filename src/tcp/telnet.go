package tcp

import (
	"fmt"
	"log"
	"time"

	"github.com/HDN-1D10T/divinity/src/util"
)

// TelnetPreflight - checks if we want to use the telnet protocol and on which port
func TelnetPreflight(hostString, ip, port, user, pass, Alert, OutputFile string) {
	if shouldTelnet(port) {
		Telnet(ip, port, user, pass, Alert, OutputFile)
	}
}

func shouldTelnet(port string) bool {
	if port == "23" || *Conf.Port == "23" {
		return true
	}
	return *Conf.Telnet && len(port) > 0
}

// Telnet - Check for valid credentials
func Telnet(ip, port, user, pass, alert, outputFile string) {
	log.Printf("Trying %s:%s %s:%s...\n", ip, port, user, pass)
	conn, err := DialTimeout("tcp", ip+":"+port, time.Duration(*Conf.Timeout)*time.Millisecond)
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
