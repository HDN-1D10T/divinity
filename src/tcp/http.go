package tcp

import (
	"crypto/tls"
	"encoding/base64"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/HDN-1D10T/divinity/src/util"
)

var m = sync.RWMutex{}

const maxHTTPBodyBytes = 10 << 20

// DoHTTPLogin checks for default credentials against an HTTP/HTTS endpoint
func DoHTTPLogin(ip string, wg *sync.WaitGroup) {
	m.RLock()
	defer m.RUnlock()
	defer wg.Done()
	timeout := httpTimeout()
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: timeout,
			}).Dial,
			TLSHandshakeTimeout:   timeout,
			ResponseHeaderTimeout: timeout,
			ExpectContinueTimeout: timeout,
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
	bodyBytes, err := io.ReadAll(io.LimitReader(res.Body, maxHTTPBodyBytes))
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

func httpTimeout() time.Duration {
	timeout := *Conf.HTTPTimeout
	if timeout <= 0 {
		timeout = 10000
	}
	return time.Duration(timeout) * time.Millisecond
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
