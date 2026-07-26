package tcp

import (
	"crypto/tls"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/HDN-1D10T/divinity/src/util"
)

var m = sync.RWMutex{}

// DoHTTPLogin checks for default credentials against an HTTP/HTTS endpoint
func DoHTTPLogin(ip string, wg *sync.WaitGroup) {
	m.RLock()
	defer m.RUnlock()
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
	protocol := strings.ToLower(*Conf.Protocol)
	port := httpPort(protocol, *Conf.Port)
	path := httpPath(*Conf.Path)
	method := strings.ToUpper(*Conf.Method)
	basicAuth := *Conf.BasicAuth
	basicAuth = base64.StdEncoding.EncodeToString([]byte(basicAuth))
	contentType := *Conf.ContentType
	headerName := *Conf.HeaderName
	headerValue := *Conf.HeaderValue
	data := *Conf.Data
	success := *Conf.Success
	alert := *Conf.Alert
	urlString := protocol + "://" + ip + ":" + port + path
	log.Println("Trying " + ip + " ...")
	// HTTP Request
	req, err := http.NewRequest(method, urlString, strings.NewReader(data))
	util.PanicErr(err)
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
		if strings.Contains(bodyString, success) {
			msg := ip + "\t" + alert
			util.LogWrite(msg)
			return
		}
		for _, v := range res.Header {
			if strings.Contains(strings.Join(v, ""), success) {
				msg := ip + "\t" + alert
				util.LogWrite(msg)
				return
			}
		}
	} else if len(basicAuth) > 0 {
		msg := ip + "\t" + alert
		util.LogWrite(msg)
	}
}

func httpPort(protocol, port string) string {
	if len(port) > 0 {
		return port
	}
	switch protocol {
	case "http":
		return "80"
	case "https":
		return "443"
	default:
		return port
	}
}

func httpPath(path string) string {
	if len(path) == 0 {
		return "/"
	}
	if strings.HasPrefix(path, "/") {
		return path
	}
	return "/" + path
}
