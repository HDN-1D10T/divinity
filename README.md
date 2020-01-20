![enter image description here](https://img.shields.io/badge/platform-ALL-green)
# HAKDEFNET / 1D10T's Project Divinity

**Divinity** is an ever-expanding HDN-Offensive Security Framework that can be used for multiple security research purposes.  It *can integrate with online search tools, but does not rely* on them, an example being **Shodan** for some of its features, the main function of which is to test over HTTP/HTTPS and report IPs that are using **default credentials**. We hope that by making this public, we can help people to test thier own systems using this opensource Framework that we decided to release to the world in an effort to make it better and widely used in increasing security. It is my hope that anyone who uses this also adds functionality and gives us credit for the work we are investing in this project. If you would like to know more about Hakdefnet then check out the site at: https://hakdefnet.org 
Enjoy and contribute
Your 1D10T / PHX / HDN Team


## Ways to run Divinity

- **local JSON config file** (specified by `-config [FILE PATH]`)
- **remote JSON config file** (specified by `-webconfig [URL]`)
- **command line parameters** (which will *override* any duplicate parameters if config is also specified) 

---

|Config Parameters |Default |Description                  |Type  |                     
|------------------|--------|-----------------------------|------|
|`-config`         |none    |path to a JSON config file   |string|
|`-webconfig`      |none    |URL to a JSON config file    |string|
|`-list`           |none    |supply list of IPs (`-list -` or `-list stdin` allows processing from `stdin` instead of a file)|string|
|`-range`          |none    |specify a CIDR range of IP addresses|string
