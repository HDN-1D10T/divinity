package tcp

import (
	"fmt"
	"net"
	"time"

	"github.com/HDN-1D10T/divinity/src/util"
)

const clear = "\033[2K"

func psWorker(host string, port string, results chan string) {
	address := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("%sspotchecking: %s\r", clear, host)
	conn, _ := net.DialTimeout("tcp", address, time.Duration(*Conf.Timeout)*time.Millisecond)
	if conn != nil {
		results <- host
		conn.SetReadDeadline(time.Now().Add(time.Second))
		conn.Close()
		return
	}
	results <- "None"
}

// Scan with multi-threaded workers
func ScanMulti(ips []string, port string) {
	if len(ips) >= 1048575 {
		fmt.Println("ERROR: Too many IPs. Try breaking up into separate lists.")
		return
	}
	results := make(chan string, len(ips))
	//openports := make(chan string)

	go func() {
		for _, ip := range ips {
			go psWorker(ip, port, results)
		}
	}()

	for i := 1; i <= cap(results); i++ {
		host := <-results
		if host == "None" {
			continue
		}
		if *Conf.IPOnly {
			util.LogWrite(host)
		} else {
			util.LogWrite(fmt.Sprintf("%s:%s open", host, port))
		}
	}
}
