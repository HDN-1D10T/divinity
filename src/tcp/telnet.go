package tcp

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

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/HDN-1D10T/divinity/src/util"
	expect "github.com/google/goexpect"
)

const timeout = 500 * time.Millisecond

func doTelnet(ip, user, pass, alert, outputFile string) {
	fmt.Printf("Trying %s:23 %s:%s...\n", ip, user, pass)
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
	if err != nil {
		return
	}
	e.Send("exit\n")
	if promptRE.MatchString(res) {
		msg := fmt.Sprintf("%s:23 %s:%s\t%s\n", ip, user, pass, alert)
		fmt.Println(msg)
		if len(outputFile) > 0 {
			util.FileWrite(msg, outputFile)
		}
	}
}

// Telnet - check for valid credentials
func Telnet(ips []string, conf *Configuration) {
	alert := *conf.Alert
	outputFile := *conf.OutputFile
	for _, ip := range ips {
		go func(ip string) {
			dumpmatch := regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}:[0-9]{1,5} .*[A-Za-z0-9].*:.*`)
			if !dumpmatch.MatchString(ip) {
				fmt.Println("Error: string formatted incorrectly" + ip)
			}
			hostString := strings.Split(ip, " ")[0]
			credString := strings.Split(ip, " ")[1]
			ip = strings.Split(hostString, ":")[0]
			creds := strings.Split(credString, ":")
			if len(creds) == 2 {
				user := creds[0]
				pass := creds[1]
				doTelnet(ip, user, pass, alert, outputFile)
			} else {
				user := creds[0]
				pass := ""
				doTelnet(ip, user, pass, alert, outputFile)
			}
		}(ip)
		time.Sleep(1 * time.Millisecond)
	}
}
