/*
##############################################################################
#    Copyright (C) 2020  Hakdefnet International <https://hakdefnet.org>
#
#    Authors:
#    phx <https://github.com/phx>
#    1D10T <https://github.com/HDN-1D10T>
#
#    This program free software licensed under Creative Commons BY-NC-ND 4.0.
#    You can redistribute it and/or modify it under the terms of the
#    Attribution-NonCommerical-NoDerivatives 4.0 International License,
#    as published by Creative Commons.
#
#    This program is distributed in the hope that it will be useful,
#    but WITHOUT ANY WARRANTY; without even the implied warranty of
#    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
#    See CC Attribution-NonCommerical-NoDerivatives 4.0 International License
#    for more details.
#
#    You should have received a copy of the CC BY-NC-ND 4.0 License along with
#    this program.  If not, see <https://creativecommons.org/licenses/>.
##############################################################################
*/

package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/HDN-1D10T/divinity/src/config"
	"github.com/HDN-1D10T/divinity/src/masscan"
	"github.com/HDN-1D10T/divinity/src/shodan"
)

// Configuration imported from src/config
type Configuration struct{ config.Options }

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func makeRange(min, max int) []int {
	r := make([]int, max-min+1)
	for i := range r {
		r[i] = min + i
	}
	return r
}

func filewrite(chunk, outputFile string) {
	f, err := os.OpenFile(outputFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)
	defer f.Close()
	if _, err := f.WriteString(chunk + "\n"); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func getIPsFromCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	check(err)
	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

var m = sync.RWMutex{}

func doLogin(ip string, conf Configuration, wg *sync.WaitGroup) {
	m.RLock()
	defer wg.Done()
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}
	protocol := *conf.Protocol
	port := *conf.Port
	path := *conf.Path
	method := *conf.Method
	basicAuth := *conf.BasicAuth
	basicAuth = base64.StdEncoding.EncodeToString([]byte(basicAuth))
	contentType := *conf.ContentType
	headerName := *conf.HeaderName
	headerValue := *conf.HeaderValue
	data := *conf.Data
	success := *conf.Success
	alert := *conf.Alert
	outputFile := *conf.OutputFile
	urlString := protocol + "://" + ip + ":" + port + path
	fmt.Println("Trying " + ip + " ...")
	// HTTP Request
	req, err := http.NewRequest(method, urlString, strings.NewReader(data))
	check(err)
	if len(headerName) > 0 {
		req.Header.Set(headerName, headerValue)
	}
	if len(basicAuth) > 0 {
		req.Header.Set("Authorization", "Basic "+basicAuth)
	}
	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return
	}
	if err != nil {
		return
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	bodyString := string(bodyBytes)
	if len(success) > 0 {
		check(err)
		if strings.Contains(bodyString, success) {
			msg := ip + "\t" + alert
			println(msg)
			if len(outputFile) > 0 {
				filewrite(msg, outputFile)
			}
		}
	} else if len(basicAuth) > 0 {
		msg := ip + "\t" + alert
		println(msg)
		if len(outputFile) > 0 {
			filewrite(msg, outputFile)
		}
	}
	m.RUnlock()
}

func mScan(cidr string) {
	m := masscan.New()
	m.SetPorts("0-65535")
	m.SetRanges(cidr)
	m.SetRate("2000")
	m.SetExclude("127.0.0.1")
	err := m.Run()
	if err != nil {
		fmt.Println("scanner failed", err)
		return
	}
	results, err := m.Parse()
	if err != nil {
		fmt.Println("Scan Results:", err)
		return
	}
	for _, result := range results {
		fmt.Println(result)
	}
}

func main() {
	runtime.GOMAXPROCS(100)
	var wg = sync.WaitGroup{}
	conf := Configuration{
		config.ParseConfiguration(),
	}
	list := *conf.List
	shodanSearch := *conf.SearchTerm
	passive := *conf.Passive
	scan := *conf.Scan
	cidr := *conf.Cidr
	ipsOnly := *conf.IPOnly
	outputFile := *conf.OutputFile
	if len(cidr) > 0 {
		if scan {
			mScan(cidr)
		} else {
			//Process IPs from CIDR range
			ips, _ := getIPsFromCIDR(cidr)
			wg.Add(len(ips))
			for _, host := range ips {
				go doLogin(host, Configuration{config.ParseConfiguration()}, &wg)
			}
			wg.Wait()
		}
	} else if len(list) == 1 || list == "stdin" {
		// Process list from stdin
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanLines)
		var ips []string
		for scanner.Scan() {
			ips = append(ips, scanner.Text())
		}
		wg.Add(len(ips))
		for _, host := range ips {
			go doLogin(host, Configuration{config.ParseConfiguration()}, &wg)
		}
		wg.Wait()
	} else if len(list) > 1 {
		// Process list from file
		file, err := os.Open(list)
		check(err)
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var ips []string
		for scanner.Scan() {
			ips = append(ips, scanner.Text())
		}
		file.Close()
		wg.Add(len(ips))
		for _, host := range ips {
			go doLogin(host, Configuration{config.ParseConfiguration()}, &wg)
		}
		wg.Wait()
	} else if passive {
		// Process list from Shodan
		apiKey := os.Getenv("SHODAN_API_KEY")
		s := shodan.New(apiKey)
		info, err := s.APIInfo()
		check(err)
		// Get Shodan IP Results
		if !ipsOnly {
			fmt.Printf(
				"Query Credits:\t%d\nScan Credits:\t%d\n\n",
				info.QueryCredits,
				info.ScanCredits)
		}
		pageRange := makeRange(1, *conf.Pages)
		for _, num := range pageRange {
			pageStr := strconv.Itoa(num)
			query := shodanSearch + "&page=" + pageStr
			hostSearch, err := s.HostSearch(query)
			check(err)
			// Run config from command line arguments:
			if ipsOnly {
				for _, host := range hostSearch.Matches {
					msg := host.IPString
					fmt.Println(msg)
					if len(outputFile) > 0 {
						filewrite(msg, outputFile)
					}
				}
			} else {
				for _, host := range hostSearch.Matches {
					msg := host.IPString + "\t\t" + host.Location.CountryName
					fmt.Println(msg)
					if len(outputFile) > 0 {
						filewrite(msg, outputFile)
					}
				}
			}
		}
	} else if len(shodanSearch) > 0 {
		// Get Shodan IP Results
		apiKey := os.Getenv("SHODAN_API_KEY")
		s := shodan.New(apiKey)
		info, err := s.APIInfo()
		check(err)
		fmt.Printf(
			"Query Credits:\t%d\nScan Credits:\t%d\n\n",
			info.QueryCredits,
			info.ScanCredits)
		pageRange := makeRange(1, *conf.Pages)
		for _, num := range pageRange {
			pageStr := strconv.Itoa(num)
			query := shodanSearch + "&page=" + pageStr
			hostSearch, err := s.HostSearch(query)
			check(err)
			wg.Add(len(hostSearch.Matches))
			// Run config from command line arguments:
			for _, host := range hostSearch.Matches {
				go doLogin(host.IPString, Configuration{config.ParseConfiguration()}, &wg)
			}
			wg.Wait()
		}
	}
}
