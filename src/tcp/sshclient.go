package tcp

import (
	"strings"
	"time"

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
				conn, err := ssh.Dial("tcp", ip+":"+sshport, sshConfig)
				if err != nil {
					messages <- "nil"
				}
				if conn != nil {
					// start session
					sess, err := conn.NewSession()
					if err != nil {
						messages <- "nil"
					}
					status := func() bool {
						if sess != nil {
							//sess.Stdout = os.Stdout
							//sess.Stderr = os.Stderr
							// run single command
							err = sess.Run("ps")
							if err == nil {
								return true
							}
							err = sess.Run("help")
							if err == nil {
								return true
							}
							err = sess.Run("sh help")
							if err == nil {
								return true
							}
							err = sess.Run("?")
							if err == nil {
								return true
							}
							return false
						}
						return false
					}()
					if status {
						messages <- ip + ":" + sshport + " " + user + ":" + pass + " " + alert
						sess.Close()
					}
					conn.Close()
				}
			}()
			time.Sleep(time.Duration(*Conf.Timeout) * time.Millisecond)
		}
	}
}
