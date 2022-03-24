package tcp

import (
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

var counter = 0

// SSHPreflight - checks if we want to use the SSH protocol and on which port
func SSHPreflight(timeStart time.Time, timeDuration time.Duration, timeComplete time.Time, status chan string, messages chan string, ipInfo chan IPinfo) {
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
					//Timeout: time.Duration(*Conf.SSHTimeout) * time.Millisecond,
					Timeout:         time.Duration(10 * time.Second),
					HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				}
				timeElapsed := time.Duration(time.Now().Sub(timeStart))
				timeUntilCompletion := time.Duration(timeComplete.Sub(time.Now()))
				status <- "CLEARSCREEN"
				status <- "STATUS:\033[K\r"
				status <- "STATUS:Start date:\t\t" + timeStart.Format(time.RFC1123) + "\n"
				status <- "STATUS:Duration:\t\t" + timeDuration.String() + "\n"
				status <- "STATUS:Elapsed:\t\t" + timeElapsed.String() + "\n"
				status <- "STATUS:Time until complete:\t" + timeUntilCompletion.String() + "\n"
				status <- "STATUS:Est. completion date:\t" + timeComplete.Format(time.RFC1123) + "\n"
				status <- "STATUS:Found:\t\t\t" + strconv.Itoa(counter) + "\n"
				status <- "STATUS:\n"
				status <- "Trying " + ip + ":" + sshport + " " + user + ":" + pass + "...\033[K\r"
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
							// run single command
							err = sess.Run("ps")
							if err == nil {
								return true
							}
							err = sess.Run("uname -a")
							if err == nil {
								return true
							}
							err = sess.Run("show help")
							if err == nil {
								return true
							}
							err = sess.Run("?")
							if err == nil {
								return true
							}
							err = sess.Run("help")
							if err == nil {
								return true
							}
							return false
						}
						return false
					}()
					if status {
						counter = counter + 1
						messages <- "@SUCCESS@" + ip + ":" + sshport + " " + user + ":" + pass + " " + alert + "\n"
						messages <- "\n"
						sess.Close()
					}
					conn.Close()
				}
			}()
			time.Sleep(time.Duration(*Conf.Timeout) * time.Millisecond)
		}
	}
}
