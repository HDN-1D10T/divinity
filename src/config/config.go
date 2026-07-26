package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
		BasicAuth:   flag.String("basic-auth", "", "plain-text HTTP Basic Auth value (username:password)"),
		Cidr:        flag.String("cidr", "", "specify CIDR range, '-' or 'stdin' instead of list of individual IPs"),
		ContentType: flag.String("content", "", "payload content type"),
		Credentials: flag.String("creds", "", "'username:password' formatted string for tcp connections"),
		Data:        flag.String("data", "", "POST form data"),
		Pages:       flag.Int("pages", 1, "[SHODAN] # of page results to return"),
		HeaderName:  flag.String("headername", "", "set a single header name"),
		HeaderValue: flag.String("headervalue", "", "set a single header value"),
		IPOnly:      flag.Bool("ips", false, "[SHODAN] setting ips will ONLY return a list of IPs that match the query, requires -passive"),
		List:        flag.String("list", "", "/path/to/ip_list, '-' or 'stdin'"),
		ListIPs:     flag.Bool("list-ips", false, "return list of IPs from -cidr or from -cidr list -list [/path/to/cidr_list]"),
		Masscan:     flag.Bool("masscan", false, "use masscan with -scan option. masscan must be installed. requires -cidr [range]"),
		Method:      flag.String("method", "GET", "HTTP Method"),
		OutputFile:  flag.String("out", "", "/path/to/outputfile"),
		Passive:     flag.Bool("passive", false, "[SHODAN] return IP passive info or actively check default creds"),
		Password:    flag.String("pass", "", "password for tcp connections"),
		Path:        flag.String("path", "/", "/path/to/login_page"),
		Port:        flag.String("port", "", "target port number"),
		Protocol:    flag.String("protocol", "", "protocol (http, https, or tcp)"),
		Routes:      flag.Bool("routes", false, "get CIDR ranges for ASNumber specified by -asn or from -list"),
		Scan:        flag.Bool("scan", false, "scan for open ports on a host, can use -masscan -cidr [range], or defaults to native portscanner"),
		ScanFast:    flag.Bool("scanfast", false, "used with -cidr or -list for multiple hosts, launches a multi-threaded scan -- may not be as accurate as single-threaded -scan option"),
		SSH:         flag.Bool("ssh", false, "force SSH connection on non-standard port"),
		SearchTerm:  flag.String("query", "", "[SHODAN] Shodan search query"),
		Success:     flag.String("success", "", "string match for successful login"),
		TopPorts:    flag.Bool("top", false, "used with -scan to scan top ports"),
		Telnet:      flag.Bool("telnet", false, "force telnet connection on non-standard port"),
		Timeout:     flag.Int("timeout", 500, "timeout in milliseconds for Telnet and scanfast TCP checks"),
		Username:    flag.String("user", "", "username for tcp connections"),
	}
	LocalConfig = flag.String("config", "", "Needs /path/to/config.json as argument")
	WebConfig   = flag.String("webconfig", "", "Needs URL to config.json as argument")
	parsed      bool
)

// ParseConfiguration from Options
func ParseConfiguration() Options {
	if parsed {
		return C
	}

	flag.Parse()
	cliValues := visitedFlagValues()

	// Parse JSON for config
	if len(*LocalConfig) > 0 {
		C.parseLocal(*LocalConfig)
	} else if len(*WebConfig) > 0 {
		C.parseRemote(*WebConfig)
	}
	if err := C.applyCLIOverrides(cliValues); err != nil {
		log.Fatalf("Error applying command-line options: %s . Terminating.", err)
	}
	parsed = true
	return C
}

func visitedFlagValues() map[string]string {
	values := make(map[string]string)
	flag.Visit(func(f *flag.Flag) {
		values[f.Name] = f.Value.String()
	})
	return values
}

func (c *Options) parseLocal(file string) {
	ct, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error opening JSON configuration (%s): %s . Terminating.", file, err)
	}
	defer ct.Close()

	ctb, err := io.ReadAll(ct)
	if err != nil {
		log.Fatalf("Error reading local JSON configuration (%s): %s . Terminating.", file, err)
	}
	err = c.applyJSON(ctb, file)
	if err != nil {
		log.Fatalf("Error applying local JSON configuration (%s): %s . Terminating.", file, err)
	}
}

func (c *Options) parseRemote(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting remote JSON configuration (%s): %s . Terminating.", url, err)
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		log.Fatalf("Error getting remote JSON configuration (%s): HTTP %s . Terminating.", url, res.Status)
	}

	ctb, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading remote JSON configuration (%s): %s . Terminating.", url, err)
	}
	err = c.applyJSON(ctb, url)
	if err != nil {
		log.Fatalf("Error applying remote JSON configuration (%s): %s . Terminating.", url, err)
	}
}

func (c *Options) applyJSON(data []byte, source string) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for name, value := range raw {
		if err := c.setFromJSON(name, value, source); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}

func (c *Options) applyCLIOverrides(values map[string]string) error {
	for name, value := range values {
		if err := c.setFromString(name, value); err != nil {
			return fmt.Errorf("-%s=%q: %w", name, value, err)
		}
	}
	return nil
}

func (c *Options) setFromJSON(name string, raw json.RawMessage, source string) error {
	switch name {
	case "alert", "asn", "basic-auth", "cidr", "content", "creds", "data", "list",
		"headername", "headervalue", "method", "out", "pass", "path", "port",
		"protocol", "query", "success", "user":
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			return err
		}
		return c.setFromString(name, value)
	case "all", "list-ips", "ips", "masscan", "passive", "routes", "scan",
		"scanfast", "ssh", "telnet", "top":
		var value bool
		if err := json.Unmarshal(raw, &value); err != nil {
			return err
		}
		return c.setFromString(name, strconv.FormatBool(value))
	case "pages", "timeout":
		var value int
		if err := json.Unmarshal(raw, &value); err != nil {
			return err
		}
		return c.setFromString(name, strconv.Itoa(value))
	default:
		log.Printf("Warning: ignoring unknown JSON configuration key %q in %s", name, source)
		return nil
	}
}

func (c *Options) setFromString(name, value string) error {
	switch name {
	case "alert":
		*c.Alert = value
	case "asn":
		*c.ASN = value
	case "basic-auth":
		*c.BasicAuth = value
	case "cidr":
		*c.Cidr = value
	case "content":
		*c.ContentType = value
	case "creds":
		*c.Credentials = value
	case "data":
		*c.Data = value
	case "list":
		*c.List = value
	case "headername":
		*c.HeaderName = value
	case "headervalue":
		*c.HeaderValue = value
	case "method":
		*c.Method = value
	case "out":
		*c.OutputFile = value
	case "pass":
		*c.Password = value
	case "path":
		*c.Path = value
	case "port":
		*c.Port = value
	case "protocol":
		*c.Protocol = value
	case "query":
		*c.SearchTerm = value
	case "success":
		*c.Success = value
	case "user":
		*c.Username = value
	case "all":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*c.All = parsed
	case "list-ips":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*c.ListIPs = parsed
	case "ips":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*c.IPOnly = parsed
	case "masscan":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*c.Masscan = parsed
	case "passive":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*c.Passive = parsed
	case "routes":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*c.Routes = parsed
	case "scan":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*c.Scan = parsed
	case "scanfast":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*c.ScanFast = parsed
	case "ssh":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*c.SSH = parsed
	case "telnet":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*c.Telnet = parsed
	case "top":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*c.TopPorts = parsed
	case "pages":
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		*c.Pages = parsed
	case "timeout":
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		*c.Timeout = parsed
	case "config", "webconfig":
		return nil
	default:
		return nil
	}
	return nil
}
