![enter image description here](https://img.shields.io/badge/platform-ALL-green)
# HAKDEFNET / 1D10T's Project Divinity

**Divinity** is an ever-expanding HDN-Offensive Security Framework that can be used for multiple security research purposes.
It *can integrate with online search tools, but does not rely* on them.  An example of one of those services is **Shodan** for some of its features,
with the main function of which is to test over HTTP/HTTPS and report IPs that are using **default credentials**. 

Many people install basic and advanced services like NetScaler, SAP, Firewalls and Routers without changing default passwords,
this fact makes all these implementations faulty and vulnerable to very simple attacks. Any scriptkiddie with a hacking GUI based tool
or pentesting image can search known lists and break into your devices easily without any real knowledge.
We wanted to level the playing field here as many critical infrastructure systems also have the same problems and this needs to stop
(config errors, using standard credentials) because this is irresponsible and unsafe in this period of cyber warfare and espionage.

We hope that by making this public, we can help people to test thier own systems using this opensource framework, which we decided
to release to the world in an effort to make it better and widely-used in order to increase security awareness, and hopefully security itself.

It is our hope that anyone who uses this tool also adds functionality and gives credit for the work we are investing into this project.
If you would like to know more about Hakdefnet, then check out the site at [https://hakdefnet.org](https://hakdefnet.org). 

Enjoy and contribute!

*-- Your 1D10T / PHX / HDN Team*

## Installation

`go get github.com/HDN-1D10T/divinity`

## Ways to run Divinity
- **local JSON config file** (specified by `-config [FILE PATH]`)
- **remote JSON config file** (specified by `-webconfig [URL]`)
- **command line parameters** (which will *override* any duplicate parameters if config is also specified) 
---
### Configuration Parameters:
|Normal Parameters |Value Description|                   
|----------------|-----------------|
|`-config`|path to a JSON config file   |
|`-webconfig`|URL to a JSON config file    |
|`-list`|path to list of IPs (value `-` or `stdin` allows processing from `stdin` instead of file)|
|`-dumplist`|path to file with format `[ip]:[port] [user]:[pass]`|
|`-cidr`|specify a CIDR range of IP addresses to run login tests or scan against|
|`-out`|specify file name or file path to save results 
|`-protocol`|specify if login target uses `HTTP`, `HTTPS`, or generic `TCP`|
|`-port`|specify port used by login target|
|`-path`|specify URL path to login page (default: `"/"`)|
|`-method`|specify HTTP method (usually `GET` or `POST`)|
|`-basic-auth`|if basic auth is needed, value should be plain-text `username:password` format|
|`-creds`|same as `-basic-auth`, except for TCP|
|`-user`|TCP username|
|`-pass`|TCP password|
|`-content`|value of `Content-Type` header when used with `-method POST`|
|`-data`|payload body when used with `-method POST`|
|`-headername`|specify an additional HTTP request header name|
|`-headervalue`|specify an additional HTTP request header value when used with `-headername [NAME]`|
|`-success`|string to match on that *ONLY* appears in successful login response|
|`-alert`|string to display when `-success` string is matched (default: `"SUCCESS"`)|
|`-scan`|actively scan IP range (if used with -masscan, requires `sudo`, `masscan`, and `-cidr`)|
---
### Shodan Configuration Parameters (optional):
If Shodan is used, you will need to set the environment variable `SHODAN_API_KEY=[your shodan API key]`.  
This can be exported on the command line or sourced in your `~/.bashrc`, etc.

|Parameters |Value Description|                   
|-----------|-----------------|
|`-query`   |Shodan search string|
|`-pages`   |Number `[type: int]` of page results to display. Best practice is to use this flag manually as to not use unnecessary query credits (default: `1`)|
|`-passive` |If this flag is set, only IPs along with associated countries will be displayed, without testing them for default credentials|
|`-ips`     |If this flag is set, *ONLY* a list of IPs will be returned in the output that matches the `-query` value (requires `-passive`). This option is good for searching a large number of `-pages` along with the `-out` parameter set, so that you can later run the tool multiple times using the `-list` parameter without an additional increment to your Shodan query credits.
---
## Example Configurations
The following configurations can be referenced locally with the `-config` parameter or hosted remotely and referenced with the `-webconfig` parameter.  Command line parameters will override any existing parameters included in the JSON configurations.

#### Basic-Auth GET Request (Device Manufacturer A):
- The following configuration would search for "Device Manufacturer A" listening on port 80 at `http://[IP ADDRESS]/login.html` and would return 1 page of Shodan results containing 100 IPs.
- An HTTP GET request would be sent with basic authentication with the credentials `admin:password` to all 100 IPs from the results.
- If the string `authenticated.html` is found in the HTTP response, the alert `*** DEFAULT CREDENTIALS ***` will be displayed next to the IP address in the output.

```
{
    "query": "Device Manufacturer A port:80",
    "protocol": "http",
    "port": "80",
    "path": "/login.html",
    "method": "GET",
    "basic-auth": "admin:password",
    "success": "authenticated.html",
    "alert": "*** DEFAULT CREDENTIALS ***"
}
```

#### POST Request (Device Manufacturer B):
- The following configuration would search for "device manufacturer 2" listening on port 8443 at `https://[IP ADDRESS]:8443/data/login` and would return 1 page of results, containing 100 IPs.
- An HTTP POST request would be sent with the credentials `admin:password` as form data to all 100 IPs.
- If the string `<authResult>0</authResult>` is found in the HTTP response, the alert `SUCCESS` will be displayed next to the IP address in the output (since `-alert` has a default value of `SUCCESS`, it doesn't have to be explicitly specified in the configuration.

```
{
    "query": "Device Manufacturer B port:8443",
    "protocol": "https",
    "port": "8443",
    "path": "/data/login",
    "method": "POST",
    "content": "application/x-www-form-url-encoded",
    "success": "<authResult>0</authResult>"
}
```
#### Example - Check 500 individual  results from Device Manufacturer B for default credentials:
Let's say you wanted to get 500 individual IP results from Device Manufacturer B, and you have the config stored in a directory on the web along with additional configs you have created for finding different devices.  You want to store these IPs for later use without increasing your Shodan API query credits.

You could use the following command to save a list of just the IP addresses to subsequently feed back into the application using the `-list` parameter, which would override the call to the Shodan API on the second run:

`divinity -webconfig http://example.com/divinity_configs/device-manufacturer-b.json -pages 5 -passive -ips -out manufacturer_b_ips.txt`

If you wanted to do everything in one go, just make sure to save your results with the `-out` parameter.  Let's also say that you want to include the text `*** DEFAULT ***` next to the successful attempts.  The output file will include only the IPs with successful default logins.

`divinity -webconfig http://example.com/divinity_configs/device-manufacturer-b.json -pages 5 -alert "*** DEFAULT ***" -out manufacturer_b_default_creds.txt`

#### Shodan-less Example - Check internal app tier for default credentials:
Let's say you have a numerous applications running a specific framework for which you have created a configuration file.  These applications are running in your DMZ on the 10.2.2.0/24 network.  As long as you have access to these applications, you can run the following command to test for default credentials:

`divinity -config /path/to/app.json -cidr 10.2.2.0/24 -out dmz_default_creds.txt`

## Portscanning
**Note:** `masscan` integration is not complete and is a work in progress.  There are also plans to implement `nmap` integration.  Both of these require spawning OS processes that require these utilities to be installed,
and additionally require `sudo` or root-level permission.

That being said, there is a native golang portscanner implemented directly in this project, and it works quite well.  Additionally, it doesn't require root-level permissions to do its thing, *and it does its thing very efficiently*.

When scanning for a single port, `divinity` can knock out a /24 in around 2 minutes.  That's over the Internet.  But of course, you would only use this tool on local networks...

### Scan Example

Let's say you wanted to find a list of IPs on a local network that were running Telnet servers.  You want to create a list of these IPs to feed back into `divinity` and test for default credentials `admin:admin`.

#### Create your list:

`divinity -scan -cidr 192.168.1.0/24 -port 23 -out telnet.txt`

#### Feed your list back in to check for default creds:

`divinity -list telnet.txt -protocol tcp -port 23 -out default_creds.txt`

#### If you want, just feed the `stdout` to `stdin`

`divinity -scan -cidr 192.168.1.0/24 -port 23 | divinity -protocol tcp -port 23 -creds admin:admin -list - -out default_creds.txt`

