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
	List        *string `json:"list"`
	Cidr        *string `json:"cidr"`
	SearchTerm  *string `json:"query"`
	Pages       *int    `json:"pages"`
	Passive     *bool   `json:"passive"`
	IPOnly      *bool   `json:"ips"`
	Protocol    *string `json:"protocol"`
	Port        *string `json:"port"`
	Path        *string `json:"path"`
	Method      *string `json:"method"`
	BasicAuth   *string `json:"basic-auth"`
	ContentType *string `json:"content"`
	HeaderName  *string `json:"headername"`
	HeaderValue *string `json:"headervalue"`
	Data        *string `json:"data"`
	Success     *string `json:"success"`
	Alert       *string `json:"alert"`
	OutputFile  *string `json:"out"`
	Scan        *bool   `json:"scan"`
	Masscan     *bool   `json:"masscan"`
	Username    *string `json:"user"`
	Password    *string `json:"pass"`
	Credentials *string `json:"creds"`
	DumpList    *string `json:"dumplist"`
}

// Options for Configuration
var (
	C = Options{
		List:        flag.String("list", "", "/path/to/ip_list"),
		Cidr:        flag.String("cidr", "", "specify CIDR range instead of list of individual IPs"),
		SearchTerm:  flag.String("query", "", "[SHODAN] Shodan search query"),
		Pages:       flag.Int("pages", 1, "[SHODAN] # of page results to return"),
		Passive:     flag.Bool("passive", false, "[SHODAN] return IP passive info or actively check default creds"),
		IPOnly:      flag.Bool("ips", false, "[SHODAN] setting ips will ONLY return a list of IPs that match the query, requires -passive"),
		Protocol:    flag.String("protocol", "", "protocol (http or https)"),
		Port:        flag.String("port", "", "port number"),
		Path:        flag.String("path", "/", "/path/to/login_page"),
		Method:      flag.String("method", "", "HTTP Method"),
		BasicAuth:   flag.String("basic-auth", "", "base64-decoded (plain-text) BasicAuth header value (username:password)"),
		ContentType: flag.String("content", "", "payload content type"),
		HeaderName:  flag.String("headername", "", "set a single header name"),
		HeaderValue: flag.String("headervalue", "", "set a single header value"),
		Data:        flag.String("data", "", "POST form data"),
		Success:     flag.String("success", "", "string match for successful login"),
		Alert:       flag.String("alert", "SUCCESS", "alert message upon success"),
		OutputFile:  flag.String("out", "", "/path/to/outputfile"),
		Scan:        flag.Bool("scan", false, "scan for open ports on a host, can use -masscan -cidr [range], or defaults to native portscanner"),
		Masscan:     flag.Bool("masscan", false, "use masscan with -scan option. masscan must be installed. requires -cidr [range]"),
		Username:    flag.String("user", "", "username for tcp connections"),
		Password:    flag.String("pass", "", "password for tcp connections"),
		Credentials: flag.String("creds", "", "'username:password' formatted string for tcp connections"),
		DumpList:    flag.String("dumplist", "", "path to file with format '[ip]:[port] [user]:[pass]'"),
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
