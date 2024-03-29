package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Options struct for Configuration
type Options struct {
	Alert       *string `json:"alert"`
	All         *bool   `json:"all"`
	ASN         *string `json:"asn"`
	BasicAuth   *string `json:"basic-auth"`
	Cidr        *string `json:"cidr"`
	ContentType *string `json:"content"`
	Credentials *string `json:"creds"`
	Data        *string `json:"data"`
	List        *string `json:"list"`
	ListIPs     *bool   `json:"list-ips"`
	HeaderName  *string `json:"headername"`
	HeaderValue *string `json:"headervalue"`
	IPOnly      *bool   `json:"ips"`
	Masscan     *bool   `json:"masscan"`
	Method      *string `json:"method"`
	OutputFile  *string `json:"out"`
	Path        *string `json:"path"`
	Pages       *int    `json:"pages"`
	Passive     *bool   `json:"passive"`
	Password    *string `json:"pass"`
	Port        *string `json:"port"`
	Protocol    *string `json:"protocol"`
	Routes      *bool   `json:"routes"`
	SearchTerm  *string `json:"query"`
	Scan        *bool   `json:"scan"`
	ScanFast    *bool   `json:"scanfast"`
	SSH         *bool   `json:"ssh"`
	Success     *string `json:"success"`
	TopPorts    *bool   `json:"top"`
	Telnet      *bool   `json:"telnet"`
	Timeout     *int    `json:"timeout"`
	Username    *string `json:"user"`
}

// Options for Configuration
var (
	C = Options{
		Alert:       flag.String("alert", "SUCCESS", "alert message upon success"),
		All:         flag.Bool("all", false, "used with -scan to scan all ports"),
		ASN:         flag.String("asn", "", "used with -routes to show CIDR blocks for ASNumber"),
		BasicAuth:   flag.String("basic-auth", "", "base64-decoded (plain-text) BasicAuth header value (username:password)"),
		Cidr:        flag.String("cidr", "", "specify CIDR range instead of list of individual IPs"),
		ContentType: flag.String("content", "", "payload content type"),
		Credentials: flag.String("creds", "", "'username:password' formatted string for tcp connections"),
		Data:        flag.String("data", "", "POST form data"),
		Pages:       flag.Int("pages", 1, "[SHODAN] # of page results to return"),
		HeaderName:  flag.String("headername", "", "set a single header name"),
		HeaderValue: flag.String("headervalue", "", "set a single header value"),
		IPOnly:      flag.Bool("ips", false, "[SHODAN] setting ips will ONLY return a list of IPs that match the query, requires -passive"),
		List:        flag.String("list", "", "/path/to/ip_list"),
		ListIPs:     flag.Bool("list-ips", false, "return list of IPs from -cidr or -list containing one or more CIDR ranges"),
		Masscan:     flag.Bool("masscan", false, "use masscan with -scan option. masscan must be installed. requires -cidr [range]"),
		Method:      flag.String("method", "", "HTTP Method"),
		OutputFile:  flag.String("out", "", "/path/to/outputfile"),
		Passive:     flag.Bool("passive", false, "[SHODAN] return IP passive info or actively check default creds"),
		Password:    flag.String("pass", "", "password for tcp connections"),
		Path:        flag.String("path", "/", "/path/to/login_page"),
		Port:        flag.String("port", "", "port number"),
		Protocol:    flag.String("protocol", "", "protocol (http or https)"),
		Routes:      flag.Bool("routes", false, "get CIDR ranges for ASNumber specified by -asn or from -list"),
		Scan:        flag.Bool("scan", false, "scan for open ports on a host, can use -masscan -cidr [range], or defaults to native portscanner"),
		ScanFast:    flag.Bool("scanfast", false, "used with -cidr or -list for multiple hosts, launches a multi-threaded scan -- may not be as accurate as single-threaded -scan option"),
		SSH:         flag.Bool("ssh", false, "force SSH connection on non-standard port"),
		SearchTerm:  flag.String("query", "", "[SHODAN] Shodan search query"),
		Success:     flag.String("success", "", "string match for successful login"),
		TopPorts:    flag.Bool("top", false, "used with -scan to scan top ports"),
		Telnet:      flag.Bool("telnet", false, "force telnet connection on non-standard port"),
		Timeout:     flag.Int("timeout", 500, "timeout in milliseconds to receive Telnet response (longer is more accurate)"),
		Username:    flag.String("user", "", "username for tcp connections"),
	}
	LocalConfig = flag.String("config", "", "Needs /path/to/config.json as argument")
	WebConfig   = flag.String("webconfig", "", "Needs URL to config.json as argument")
)

// ParseConfiguration from Options
func ParseConfiguration() Options {

	flag.Parse()

	// Parse JSON for config
	if len(*LocalConfig) > 0 {
		C.parseLocal(*LocalConfig)
	} else if len(*WebConfig) > 0 {
		C.parseRemote(*WebConfig)
	}
	return C
}

func (c *Options) parseLocal(file string) {
	ct, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error opening JSON configuration (%s): %s . Terminating.", file, err)
	}
	defer ct.Close()

	ctb, _ := ioutil.ReadAll(ct)
	err = json.Unmarshal(ctb, &c)
	if err != nil {
		log.Fatalf("Error unmarshalling local JSON configuration (%s): %s . Terminating.", file, err)
	}
}

func (c *Options) parseRemote(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting remote JSON configuration (%s): %s . Terminating.", url, err)
	}
	defer res.Body.Close()

	//var ret Options
	err = json.NewDecoder(res.Body).Decode(&c)
	if err != nil {
		log.Fatalf("Error decoding remote JSON configuration (%s): %s . Terminating.", url, err)
	}
}
