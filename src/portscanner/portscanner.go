package portscanner

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"time"

	"github.com/HDN-1D10T/divinity/src/util"

	"github.com/HDN-1D10T/divinity/src/config"
)

// Configuration imported from src/config
type Configuration struct{ config.Options }

func worker(ports, results chan int, host string) {
	for p := range ports {
		//fmt.Printf("Trying %s:%d...\n", host, p)
		address := fmt.Sprintf("%s:%d", host, p)
		conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

// Scan with native portscanner
func Scan(host string, portString string) {
	conf := Configuration{
		config.ParseConfiguration(),
	}
	ports := make(chan int, 100)
	results := make(chan int)

	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results, host)
	}

	if portString == "all" {
		go func() {
			for i := 1; i <= 1024; i++ {
				ports <- i
			}
		}()

		for i := 0; i < 1024; i++ {
			port := <-results
			if port != 0 {
				openports = append(openports, port)
			}
		}
		close(ports)
		close(results)
		sort.Ints(openports)
		for _, port := range openports {
			host := host
			msg := fmt.Sprintf("%s:%d\n", host, port)
			fmt.Printf(msg)
			if len(*conf.OutputFile) > 0 {
				util.FileWrite(msg, *conf.OutputFile)
			}
		}
	} else {
		thePort, _ := strconv.Atoi(portString)
		ports <- thePort
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
		close(ports)
		close(results)
		sort.Ints(openports)
		for range openports {
			host := host
			msg := fmt.Sprintf("%s\n", host)
			fmt.Printf(msg)
			if len(*conf.OutputFile) > 0 {
				util.FileWrite(msg, *conf.OutputFile)
			}
		}
	}
}
