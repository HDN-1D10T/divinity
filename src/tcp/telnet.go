package tcp

import (
	"fmt"
	"log"
	"regexp"
	"time"

	expect "github.com/google/goexpect"
)

const timeout = 500 * time.Millisecond

// Telnet - check for valid credentials
func Telnet(ip, user, pass string) {
	userRE := regexp.MustCompile(`.*([Ll]ogin)|([Uu]sername).*`)
	passRE := regexp.MustCompile(".*[Pp]assword.*")
	promptRE := regexp.MustCompile(`.*[#\$%>].*`)
	e, _, err := expect.Spawn(fmt.Sprintf("telnet %s", ip), timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	e.Expect(userRE, timeout)
	e.Send(user + "\n")
	e.Expect(passRE, timeout)
	e.Send(pass + "\n")
	res, _, err := e.Expect(promptRE, timeout)
	e.Send("exit\n")
	if promptRE.MatchString(res) {
		fmt.Println("good login")
	} else {
		fmt.Println("bad login")
	}
}
