/*
The included "github.com/google/goexpect" package is used under
the BSD-3-Clause listed below.  All other code falls under the
default Creative Commons License for this project.  BSD-3-Clause
for "github.com/google/goexpect" is as follows:

Copyright (c) 2015 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package tcp

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/HDN-1D10T/divinity/src/util"
	expect "github.com/google/goexpect"
)

const timeout = 500 * time.Millisecond

var m = sync.RWMutex{}

func doTelnet(ip, user, pass, outputFile string, wg *sync.WaitGroup) {
	m.Lock()
	defer m.Unlock()
	fmt.Println("Trying " + ip + ":23 " + user + ":" + pass + "...")
	userRE := regexp.MustCompile(`.*([Ll]ogin)|([Uu]sername).*`)
	passRE := regexp.MustCompile(".*[Pp]assword.*")
	promptRE := regexp.MustCompile(`.*[#\$%>].*`)
	stuffRE := regexp.MustCompile(`.*[A-Za-z0-9].*`)
	e, _, err := expect.Spawn(fmt.Sprintf("telnet %s", ip), timeout)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer e.Close()
	e.Expect(userRE, timeout)
	e.Send(user + "\n")
	e.Expect(passRE, timeout)
	if stuffRE.MatchString(pass) {
		e.Send(pass + "\n")
	} else {
		e.Send("\n")
	}
	res, _, err := e.Expect(promptRE, timeout)
	e.Send("exit\n")
	if promptRE.MatchString(res) {
		fmt.Printf("%s:23 %s:%s\t*** GOOD ***\n", ip, user, pass)
		if len(outputFile) > 0 {
			util.FileWrite(ip+":23 "+user+":"+pass+"\t*** GOOD ***", outputFile)
		}
	}
}

// Telnet - check for valid credentials
func Telnet(ips []string, conf *Configuration) {
	var wg = sync.WaitGroup{}
	outputFile := *conf.OutputFile
	wg.Add(len(ips))
	defer wg.Wait()
	for _, ip := range ips {
		dumpmatch := regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}:[0-9]{1,5} .*[A-Za-z0-9].*:.*`)
		if !dumpmatch.MatchString(ip) {
			//wg.Done()
			return
		}
		hostString := strings.Split(ip, " ")[0]
		credString := strings.Split(ip, " ")[1]
		ip = strings.Split(hostString, ":")[0]
		creds := strings.Split(credString, ":")
		if len(creds) == 2 {
			user := creds[0]
			pass := creds[1]
			go doTelnet(ip, user, pass, outputFile, &wg)
		} else {
			user := creds[0]
			pass := ""
			go doTelnet(ip, user, pass, outputFile, &wg)
		}
	}
}
