package tcp

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/HDN-1D10T/divinity/src/util"
)

// SSHPreflight - checks if we want to use the telnet protocol and on which port
func SSHPreflight(hostString, ip, port, user, pass, Alert, OutputFile string) {
	if Port == "22" {
		SSH(ip, Port, user, pass, Alert, OutputFile)
	}
	if len(strings.Split(hostString, ":")) > 1 {
		port = strings.Split(hostString, ":")[1]
		port = strings.Replace(hostString, " ", "", -1)
		SSH(ip, port, user, pass, Alert, OutputFile)
	}
	if *Conf.SSH {
		SSH(ip, port, user, pass, Alert, OutputFile)
	}
}

// SSH - Check for valid credentials
func SSH(ip, port, user, pass, alert, outputFile string) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		Timeout:         250 * time.Millisecond,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	log.Printf("Trying %s:%s %s:%s...\n", ip, port, user, pass)
	conn, err := ssh.Dial("tcp", ip+":"+port, sshConfig)
	if err != nil {
		return
	}
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		// log.Println(err)
		return
	}
	sessionErr := session.Run("help")
	if sessionErr != nil {
		session.Close()
		return
	}
	defer session.Close()
	msg := fmt.Sprintf("%s:%s %s:%s %s", ip, port, user, pass, alert)
	util.LogWrite(msg)
	return
	//conn.SetReadDeadline(time.Now().Add(time.Second))
	/*
			authlogin, err := conn.ReadUntil("login:")
			if err != nil {
				//log.Println(err)
				return
			}

		conn.SetReadDeadline(time.Now().Add(time.Second))
		authpass, err := conn.ReadUntil("password:")
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
	*/
}
