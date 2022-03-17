package tcp

import (
	"strings"
	"time"

	"github.com/HDN-1D10T/divinity/src/util"
	"golang.org/x/crypto/ssh"
)

// SSHPreflight - checks if we want to use the SSH protocol and on which port
func SSHPreflight(messages chan string, ipInfo chan IPinfo) {
	for info := range ipInfo {
		hostString := info.hostString
		ip := info.ip
		port := info.port
		user := info.user
		pass := info.pass
		alert := info.alert
		doSSH, sshport := func() (bool, string) {
			if port == "22" {
				return true, port
			}
			if len(strings.Split(hostString, ":")) > 1 {
				port = strings.Split(hostString, ":")[1]
				port = strings.Replace(hostString, " ", "", -1)
				return true, port
			}
			if *Conf.SSH {
				return true, port
			}
			return false, ""
		}()
		if doSSH {
			sshConfig := &ssh.ClientConfig{
				User: user,
				Auth: []ssh.AuthMethod{
					ssh.Password(pass),
				},
				//Timeout: time.Duration(*Conf.SSHTimeout) * time.Millisecond,
				Timeout:         time.Duration(10 * time.Second),
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}
			go func() {
				messages <- "Trying " + ip + ":" + sshport + " " + user + ":" + pass + "..."
				conn, _ := ssh.Dial("tcp", ip+":"+sshport, sshConfig)
				if conn != nil {
					msg := ip + ":" + sshport + " " + user + ":" + pass + " " + alert
					util.FileWrite(msg)
					messages <- msg
					conn.Close()
				}
			}()
			time.Sleep(time.Duration(*Conf.Timeout) * time.Millisecond)
		}
	}
}
