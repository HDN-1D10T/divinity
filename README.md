![THE WORD](./logo.png?raw=true)

![platform](https://img.shields.io/badge/platform-ALL-green)

# Project Divinity by HAKDEFNET / 1D10T / PHX

**Divinity** is an HDN offensive security framework for authorized security research and internal defensive validation. It can work from local target lists, CIDR ranges, piped input, or Shodan search results. Its main tested workflows are:

- Building target lists from scans, CIDR ranges, Shodan, or ASN route lookups.
- Checking HTTP/HTTPS endpoints for default or known credentials.
- Checking SSH and Telnet services for default or known credentials.
- Saving results so scans and credential checks can be chained without repeating discovery work.

Use Divinity only on systems you own or have explicit permission to test.

The original project goal is still the same: expose weak deployments that use default or standard credentials so operators can find and fix them before someone else does.

## Installation

```sh
go install github.com/HDN-1D10T/divinity@latest
```

From a local checkout:

```sh
go run .
```

or build a binary:

```sh
go build -o divinity .
```

## Configuration Order

Divinity can be configured three ways:

1. Defaults built into the tool.
2. A JSON config from `-config` or `-webconfig`.
3. Explicit command-line flags.

Command-line flags override duplicate values from JSON config files. If both `-config` and `-webconfig` are supplied, the local `-config` file is used.

There are no positional arguments. Put every option behind a flag such as `-cidr`, `-list`, `-protocol`, or `-port`.

## Quick Patterns

Check an HTTP endpoint list with Basic Auth:

```sh
divinity -list targets.txt -protocol http -port 80 -path /login.html -basic-auth admin:password -success authenticated.html
```

Check a CIDR for Telnet default credentials and save successful logins:

```sh
divinity -cidr 192.168.1.0/24 -protocol tcp -port 23 -creds admin:admin -out telnet-defaults.txt
```

Scan for hosts with port 23 open, then feed the host list into a credential check:

```sh
divinity -scan -cidr 192.168.1.0/24 -port 23 -out telnet-hosts.txt
divinity -list telnet-hosts.txt -protocol tcp -port 23 -creds admin:admin -out telnet-defaults.txt
```

Pipe scan output directly into a credential check:

```sh
divinity -scan -cidr 192.168.1.0/24 -port 23 | divinity -protocol tcp -port 23 -creds admin:admin -list - -out telnet-defaults.txt
```

Export Shodan IPs without credential checks:

```sh
divinity -query "Device Manufacturer B port:8443" -pages 5 -passive -ips -out manufacturer-b-ips.txt
```

## Inputs

### `-list`

`-list` reads one target per line from a file:

```text
192.0.2.10
192.0.2.11:23
192.0.2.12 admin:admin
192.0.2.13:2222 root:toor
```

Use `-list -` or `-list stdin` to read targets from standard input.

Inline ports are useful when targets use mixed ports. A global `-port` overrides inline ports.

### `-cidr`

`-cidr` expands an IPv4 CIDR range:

```sh
divinity -cidr 10.2.2.0/24 -list-ips
```

`/32` produces the single host. `/0` is rejected intentionally. For normal network ranges, Divinity skips the network and broadcast addresses.

Use `-cidr -` or `-cidr stdin` to read CIDR ranges from standard input.

To expand CIDRs stored in a file, use `-cidr list` with `-list`:

```sh
divinity -cidr list -list cidrs.txt -list-ips -out expanded-ips.txt
```

## Option Reference

| Flag | JSON key | Type | Description |
| --- | --- | --- | --- |
| `-config` | none | string | Path to a local JSON config file. |
| `-webconfig` | none | string | URL to a remote JSON config file. Ignored when `-config` is also set. |
| `-out` | `out` | string | File to append results to. Results also print to stdout in most modes. |
| `-list` | `list` | string | Target list path, `-`, or `stdin`. |
| `-cidr` | `cidr` | string | IPv4 CIDR range, `-`, `stdin`, or `list` when expanding CIDRs from `-list`. |
| `-list-ips` | `list-ips` | bool | Print expanded IPs from `-cidr` or `-cidr list -list`. |
| `-protocol` | `protocol` | string | `http`, `https`, or `tcp`. Input is case-insensitive. |
| `-port` | `port` | string | Target port. Overrides inline ports in `-list`. HTTP defaults to `80`; HTTPS defaults to `443`. |
| `-path` | `path` | string | HTTP path. Default is `/`; missing leading slash is added automatically. |
| `-method` | `method` | string | HTTP method. Default is `GET`; input is uppercased before use. |
| `-basic-auth` | `basic-auth` | string | Plain-text HTTP Basic Auth value, formatted as `username:password`. Divinity base64-encodes it for the request. |
| `-content` | `content` | string | HTTP `Content-Type` header, usually with `-method POST`. |
| `-data` | `data` | string | HTTP request body, usually with `-method POST`. |
| `-headername` | `headername` | string | Additional HTTP request header name. |
| `-headervalue` | `headervalue` | string | Additional HTTP request header value. Used with `-headername`. |
| `-success` | `success` | string | String to match in the HTTP response body or response headers. |
| `-alert` | `alert` | string | Text shown next to successful results. Default: `SUCCESS`. |
| `-creds` | `creds` | string | TCP credential pair, formatted as `username:password`. Passwords may contain additional `:` characters. |
| `-user` | `user` | string | TCP username. Overrides `-creds` and per-line credentials when set. |
| `-pass` | `pass` | string | TCP password. Overrides `-creds` and per-line credentials when set. |
| `-ssh` | `ssh` | bool | Force SSH checks on a non-standard port. Use with `-port` or inline `IP:PORT` targets. |
| `-telnet` | `telnet` | bool | Force Telnet checks on a non-standard port. Use with `-port` or inline `IP:PORT` targets. |
| `-timeout` | `timeout` | int | Timeout in milliseconds for Telnet and `-scanfast` TCP checks. Default: `500`. |
| `-scan` | `scan` | bool | Run the native scanner, or masscan when paired with `-masscan`. |
| `-scanfast` | `scanfast` | bool | Fast multi-host single-port scanner. Requires `-port`. |
| `-all` | `all` | bool | With native `-scan`, scan ports `1-65535`. |
| `-top` | `top` | bool | With native `-scan`, scan the built-in top port list. |
| `-masscan` | `masscan` | bool | With `-scan`, use the external `masscan` binary. Requires `masscan` and usually root/sudo privileges. |
| `-query` | `query` | string | Shodan search query. |
| `-pages` | `pages` | int | Number of Shodan result pages to request. Default: `1`. |
| `-passive` | `passive` | bool | Query Shodan and print results without credential checks. |
| `-ips` | `ips` | bool | With `-passive`, print only IPs. |
| `-routes` | `routes` | bool | Query RADB whois for IPv4 CIDR routes from `-asn` or ASNs in `-list`. |
| `-asn` | `asn` | string | ASN for `-routes`. Accepts `12345` or `AS12345`. |

Unknown JSON keys are ignored with a warning. JSON values must use the types shown above.

## JSON Configuration

JSON config keys match the option names without the leading dash. For example:

```json
{
  "protocol": "https",
  "port": "8443",
  "path": "/data/login",
  "method": "POST",
  "basic-auth": "admin:password",
  "content": "application/x-www-form-urlencoded",
  "data": "username=admin&password=password",
  "success": "<authResult>0</authResult>",
  "alert": "*** DEFAULT CREDENTIALS ***"
}
```

Run it locally:

```sh
divinity -config /path/to/device.json -cidr 10.2.2.0/24 -out defaults.txt
```

Override a JSON value on the command line:

```sh
divinity -config /path/to/device.json -port 9443 -cidr 10.2.2.0/24
```

In this example, `-port 9443` wins over the `port` value inside the config file.

## HTTP/HTTPS Credential Checks

HTTP checks are selected with `-protocol http` or `-protocol https`. Divinity builds one request per target:

```sh
divinity -list targets.txt -protocol https -port 8443 -path /login -method POST -content application/x-www-form-urlencoded -data "user=admin&pass=password" -success "Welcome"
```

Behavior notes:

- `-method` defaults to `GET`.
- `-path` defaults to `/`.
- `http` defaults to port `80` and `https` defaults to port `443` when `-port` is omitted.
- `-basic-auth` is for HTTP Basic Auth only.
- `-creds`, `-user`, and `-pass` are TCP credential options, not HTTP form fillers.
- If `-success` is set, Divinity logs the target only when that string appears in the response body or headers.
- If `-success` is omitted but `-basic-auth` is set, a `200 OK` response is treated as success.

### Basic Auth GET Example

```json
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

```sh
divinity -config manufacturer-a.json -cidr 192.0.2.0/24 -out manufacturer-a-defaults.txt
```

### POST Example

```json
{
  "query": "Device Manufacturer B port:8443",
  "protocol": "https",
  "port": "8443",
  "path": "/data/login",
  "method": "POST",
  "content": "application/x-www-form-urlencoded",
  "data": "username=admin&password=password",
  "success": "<authResult>0</authResult>"
}
```

```sh
divinity -config manufacturer-b.json -cidr 10.2.2.0/24 -alert "*** DEFAULT ***" -out manufacturer-b-defaults.txt
```

## TCP Credential Checks

`-protocol tcp` currently routes credential checks through the SSH and Telnet handlers. There is not a separate generic raw TCP authentication workflow in the code.

Credential precedence for TCP is:

1. `-user` and/or `-pass`
2. `-creds`
3. Per-line credentials from `-list`

Standard ports are selected by `-port`:

```sh
divinity -cidr 192.168.1.0/24 -protocol tcp -port 23 -creds admin:admin
divinity -cidr 192.168.1.0/24 -protocol tcp -port 22 -creds root:root
```

Use `-telnet` or `-ssh` for non-standard ports:

```sh
divinity -list telnet-custom.txt -protocol tcp -telnet -creds admin:admin
divinity -list ssh-custom.txt -protocol tcp -ssh -creds root:toor
```

`telnet-custom.txt` can contain:

```text
192.0.2.10:2323
192.0.2.11:2023 admin:admin
```

`ssh-custom.txt` can contain:

```text
192.0.2.20:2222
192.0.2.21:2200 root:toor
```

## Scanning

### Native Scanner

Scan a single port across a CIDR:

```sh
divinity -scan -cidr 192.168.1.0/24 -port 23 -out telnet-hosts.txt
```

Single-port native scans output plain IPs, which are convenient to feed back into credential checks.

Scan default ports `1-1024`:

```sh
divinity -scan -cidr 192.168.1.0/24 -out open-ports.txt
```

Scan all ports:

```sh
divinity -scan -cidr 192.168.1.0/24 -all -out all-open-ports.txt
```

Scan the built-in top port list:

```sh
divinity -scan -cidr 192.168.1.0/24 -top -out top-open-ports.txt
```

Multi-port scans output `IP:PORT`.

### Fast Single-Port Scanner

`-scanfast` is a faster single-port check for many hosts:

```sh
divinity -scanfast -cidr 192.168.1.0/24 -port 443 -out https-open.txt
```

Use `-ips` to print only IPs:

```sh
divinity -scanfast -cidr 192.168.1.0/24 -port 443 -ips -out https-ips.txt
```

`-timeout` controls the TCP dial timeout in milliseconds.

### Masscan

Masscan integration is intentionally thin and still depends on the external `masscan` binary:

```sh
sudo divinity -scan -masscan -cidr 192.168.1.0/24 -port 80 -out masscan-80.txt
```

When `-port` is omitted, Divinity asks masscan to scan `0-65535`. Native `-top` selection does not currently apply to masscan.

## Shodan

Set your Shodan API key first:

```sh
export SHODAN_API_KEY=your-shodan-api-key
```

Passive mode prints Shodan results without credential checks:

```sh
divinity -query "Device Manufacturer B port:8443" -pages 5 -passive
```

Print only IPs and save them for later:

```sh
divinity -query "Device Manufacturer B port:8443" -pages 5 -passive -ips -out manufacturer-b-ips.txt
```

Active Shodan mode uses the Shodan result IPs as HTTP/HTTPS credential-check targets:

```sh
divinity -config manufacturer-b.json -query "Device Manufacturer B port:8443" -pages 5 -out manufacturer-b-defaults.txt
```

Best practice: use `-passive -ips -out` for large Shodan exports, then run credential checks from `-list` so you do not spend query credits repeatedly.

## ASN Route Lookup

Look up IPv4 routes for one ASN:

```sh
divinity -routes -asn AS12345 -out routes.txt
```

The `AS` prefix is optional:

```sh
divinity -routes -asn 12345
```

Look up routes for multiple ASNs from a file:

```sh
divinity -routes -list asns.txt -out routes.txt
```

`asns.txt` can contain:

```text
AS12345
AS64496
64500
```

## Output Behavior

`-out` appends to the target file; it does not truncate existing files.

Most modes print to stdout and append to `-out` when set. Some scanner output is intentionally formatted for piping:

- `-scan -port PORT` prints plain IPs.
- `-scan` without `-port`, `-scan -all`, and `-scan -top` print `IP:PORT`.
- `-scanfast -ips` prints plain IPs.
- `-scanfast` without `-ips` prints `IP:PORT open`.

## Development

Run tests with a writable Go build cache if your environment restricts writes to the default cache:

```sh
GOCACHE=/private/tmp/divinity-gocache go test ./...
```

Run vet:

```sh
GOCACHE=/private/tmp/divinity-gocache go vet ./...
```
