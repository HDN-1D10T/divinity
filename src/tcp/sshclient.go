package tcp

import (
	"strings"
	"time"

	"github.com/HDN-1D10T/divinity/src/util"
	"golang.org/x/crypto/ssh"
)

// SSHPreflight - checks if we want to use the SSH protocol and on which port
func SSHPreflight(chSuccess chan int, ipInfo chan IPinfo) {
	var successCount = 0
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
			go func() {
				sshConfig := &ssh.ClientConfig{
					User: user,
					Auth: []ssh.AuthMethod{
						ssh.Password(pass),
					},
					Timeout: time.Duration(*Conf.Timeout) * time.Millisecond,
					//Timeout:         time.Duration(10 * time.Second),
					HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				}
				//fmt.Print("Trying " + ip + ":" + sshport + " " + user + ":" + pass + "...\033[K\r")
				conn, _ := ssh.Dial("tcp", ip+":"+sshport, sshConfig)
				//time.Sleep(time.Duration(*Conf.Timeout) * time.Millisecond)
				if conn != nil {
					// start session
					sess, err := conn.NewSession()
					status := func() bool {
						if err != nil {
							return false
						}
						if sess != nil {
							// run single command
							err = sess.Run("uptime")
							if err == nil {
								return true
							}
						}
						return false
					}()
					if status {
						successCount += 1
						msg := ip + ":" + sshport + " " + user + ":" + pass + " " + alert + "\n"
						//fmt.Print("\033[K\r" + msg)
						util.FileWrite(msg)
						sess.Close()
					}
					conn.Close()
				}
			}()
		}
		chSuccess <- successCount
	}
}
