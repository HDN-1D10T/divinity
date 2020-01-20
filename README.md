![enter image description here](https://img.shields.io/badge/platform-ALL-green)
# Project Divinity

**Divinity** is an ever-expanding framework that can be used for multiple purposes.  It *integrates, but does not rely* on **Shodan** for some of its features, the main function of which is to test over HTTP/HTTPS and report IPs that are using **default credentials**.


## Ways to run Divinity

- **local JSON config file** (specified by `-config [FILE PATH]`)
- **remote JSON config file** (specified by `-webconfig [URL]`)
- **command line parameters** (which will *override* any duplicate parameters if config is also specified) 

---

|Config Parameters |Default |Description                  |Type  |                     
|------------------|--------|-----------------------------|------|
|`-config`         |none    |path to a JSON config file   |string|
|`-webconfig`      |none    |URL to a JSON config file    |string|
|`-list`           |none    |supply list of IPs (`-list -` or `-list stdin` allows processing from `stdin` instead of a file|string|
|`-range`          |none    |specify a CIDR range of IP addresses|string
