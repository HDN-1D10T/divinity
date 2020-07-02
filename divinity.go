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
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"sync"

	"github.com/HDN-1D10T/divinity/src/util"

	"github.com/HDN-1D10T/divinity/src/config"
	"github.com/HDN-1D10T/divinity/src/masscan"
	"github.com/HDN-1D10T/divinity/src/shodan"
	"github.com/HDN-1D10T/divinity/src/tcp"
)

// Configuration imported from src/config
type Configuration struct{ config.Options }

func makeRange(min, max int) []int {
	r := make([]int, max-min+1)
	for i := range r {
		r[i] = min + i
	}
	return r
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
	var ips []string
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Fatal(err)
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	single := regexp.MustCompile(`/32$`)
	if single.MatchString(cidr) {
		return ips, nil
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
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
	cidr := *conf.Cidr
	list := *conf.List
	ipsOnly := *conf.IPOnly
	masscan := *conf.Masscan
	passive := *conf.Passive
	protocol := *conf.Protocol
	scan := *conf.Scan
	shodanSearch := *conf.SearchTerm
	// Process list from CIDR range
	if len(cidr) > 0 {
		ips, _ := getIPsFromCIDR(cidr)
		if scan {
			// Scan with Masscan
			if masscan {
				mScan(cidr)
				return
			}
			for _, host := range ips {
				// Scan with native scanner
				tcp.Scan(host)
			}
			return
		}
		if protocol == "tcp" {
			tcp.Handler(ips)
		}
		if protocol == "http" || protocol == "https" {
			for _, host := range ips {
				wg.Add(1)
				go tcp.DoHTTPLogin(host, &wg)
			}
			wg.Wait()
		}
		return
	}
	// Process list from stdin
	if len(list) == 1 || list == "stdin" {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanLines)
		var ips []string
		for scanner.Scan() {
			if scanner.Text() != "" {
				ips = append(ips, scanner.Text())
			}
		}
		if scan {
			// Scan with Masscan
			if masscan {
				mScan(cidr)
				return
			}
			for _, host := range ips {
				// Scan with native scanner
				tcp.Scan(host)
			}
			return
		}
		if *conf.Protocol == "tcp" {
			tcp.Handler(ips)
			return
		}
		wg.Add(len(ips))
		for _, host := range ips {
			go tcp.DoHTTPLogin(host, &wg)
		}
		wg.Wait()
		return
	}
	// Process list from file
	if len(list) > 1 {
		file, err := os.Open(list)
		if err != nil {
			util.PanicErr(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var ips []string
		for scanner.Scan() {
			if scanner.Text() != "" {
				ips = append(ips, scanner.Text())
			}
		}
		if scan {
			// Scan with Masscan
			if masscan {
				mScan(cidr)
				return
			}
			for _, host := range ips {
				// Scan with native scanner
				tcp.Scan(host)
			}
			return
		}
		if *conf.Protocol == "tcp" {
			tcp.Handler(ips)
			file.Close()
			return
		}
		file.Close()
		// wg.Add(len(ips))
		for _, host := range ips {
			wg.Add(1)
			go tcp.DoHTTPLogin(host, &wg)
		}
		wg.Wait()
		return
	}
	// Process list from Shodan in passive mode
	if passive {
		apiKey := os.Getenv("SHODAN_API_KEY")
		s := shodan.New(apiKey)
		info, err := s.APIInfo()
		util.PanicErr(err)
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
			util.PanicErr(err)
			// Run config from command line arguments:
			if ipsOnly {
				for _, host := range hostSearch.Matches {
					msg := host.IPString
					util.LogWrite(msg)
				}
			} else {
				for _, host := range hostSearch.Matches {
					msg := host.IPString + "\t\t" + host.Location.CountryName
					util.LogWrite(msg)
				}
			}
		}
		return
	}
	// Shodan active mode
	if len(shodanSearch) > 0 {
		apiKey := os.Getenv("SHODAN_API_KEY")
		s := shodan.New(apiKey)
		info, err := s.APIInfo()
		util.PanicErr(err)
		fmt.Printf(
			"Query Credits:\t%d\nScan Credits:\t%d\n\n",
			info.QueryCredits,
			info.ScanCredits)
		pageRange := makeRange(1, *conf.Pages)
		for _, num := range pageRange {
			pageStr := strconv.Itoa(num)
			query := shodanSearch + "&page=" + pageStr
			hostSearch, err := s.HostSearch(query)
			util.PanicErr(err)
			// wg.Add(len(hostSearch.Matches))
			for _, host := range hostSearch.Matches {
				wg.Add(1)
				go tcp.DoHTTPLogin(host.IPString, &wg)
			}
			wg.Wait()
		}
	}
	return
}
