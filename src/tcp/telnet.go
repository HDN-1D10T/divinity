package tcp

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/HDN-1D10T/divinity/src/util"
)

const timeout = 500 * time.Millisecond

// Telnet - Check for valid credentials
func Telnet(ips []string, conf *Configuration) {
	alert := *conf.Alert
	outputFile := *conf.OutputFile
	nouserRE := regexp.MustCompile(`^:.+`)
	nopassRE := regexp.MustCompile(`.+:$`)
	for _, ip := range ips {
		go func(ip string) {
			dumpmatch := regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}:[0-9]{1,5} .*:.*`)
			if !dumpmatch.MatchString(ip) {
				log.Println("string formatted incorrectly" + ip)
				return
			}
			hostString := strings.Split(ip, " ")[0]
			credString := strings.Split(ip, " ")[1]
			ip = strings.Split(hostString, ":")[0]
			creds := strings.Split(credString, ":")
			if len(creds) == 2 {
				user := creds[0]
				pass := creds[1]
				doTelnet(ip, user, pass, alert, outputFile)
			} else if nouserRE.MatchString(credString) {
				user := ""
				pass := creds[0]
				doTelnet(ip, user, pass, alert, outputFile)
			} else if nopassRE.MatchString(credString) {
				user := creds[0]
				pass := ""
				doTelnet(ip, user, pass, alert, outputFile)
			}
		}(ip)
		time.Sleep(20 * time.Millisecond)
	}
}

func doTelnet(ip, user, pass, alert, outputFile string) {
	//log.Printf("Trying %s:23 %s:%s...\n", ip, user, pass)
	userRE := regexp.MustCompile(`.*([Ll]ogin)|([Uu]sername).*`)
	passRE := regexp.MustCompile(".*[Pp]assword.*")
	promptRE := regexp.MustCompile(`.*[#\$%>].*`)
	badRE := regexp.MustCompile(`.*[Ii]ncorrect.*`)
	conn, err := DialTimeout("tcp", ip+":23", timeout)
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
				msg := fmt.Sprintf("%s:23 %s:%s %s\n", ip, user, pass, alert)
				fmt.Printf(msg)
				if len(outputFile) > 0 {
					util.FileWrite(msg, outputFile)
				}
			}
		}
	}
}
