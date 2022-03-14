package tcp

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/HDN-1D10T/divinity/src/tcp/ssh" // slightly-altered from golang.org/x/crypto/ssh
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
		Timeout:         time.Duration(*Conf.SSHTimeout) * time.Millisecond,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	log.Printf("Trying %s:%s %s:%s...\n", ip, port, user, pass)
	// Pass a context with a timeout to tell a blocking function that it
	// should abandon its work after the timeout elapses.
	conn, err := ssh.Dial("tcp", ip+":"+port, sshConfig)
	if err != nil || conn == nil {
		return
	}
	defer conn.Close()
	/*
		session, err := conn.NewSession()
		defer session.Close()
		if err != nil {
			return
		}
	*/
	/*
		sessionErr := session.Run("help")
		if sessionErr != nil {
			log.Println(sessionErr)
			return
		}
	*/
	msg := fmt.Sprintf("%s:%s %s:%s %s", ip, port, user, pass, alert)
	util.LogWrite(msg)
	return
}
