package tcp

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/HDN-1D10T/divinity/src/util"
	"golang.org/x/crypto/ssh"
)

// SSHPreflight - checks if we want to use the SSH protocol and on which port
func SSHPreflight(chSuccess chan int, ipInfo chan IPinfo) {
	defer close(chSuccess)
	var successCount = 0
	results := make(chan bool)
	var wg sync.WaitGroup

	go func() {
		for info := range ipInfo {
			wg.Add(1)
			go func(info IPinfo) {
				defer wg.Done()
				doSSH, sshport := shouldSSH(info)
				if !doSSH {
					results <- false
					return
				}
				if trySSH(info.ip, sshport, info.user, info.pass) {
					msg := fmt.Sprintf("%s:%s %s:%s %s", info.ip, sshport, info.user, info.pass, info.alert)
					util.LogWrite(msg)
					results <- true
					return
				}
				results <- false
			}(info)
		}
		wg.Wait()
		close(results)
	}()

	for ok := range results {
		if ok {
			successCount += 1
		}
		chSuccess <- successCount
	}
}

func shouldSSH(info IPinfo) (bool, string) {
	if info.port == "22" || *Conf.Port == "22" {
		return true, info.port
	}
	if strings.Contains(info.hostString, ":") && len(info.port) > 0 {
		return *Conf.SSH, info.port
	}
	return *Conf.SSH && len(info.port) > 0, info.port
}

func trySSH(ip, port, user, pass string) bool {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", ip+":"+port, sshConfig)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
