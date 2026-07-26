package tcp

import "testing"

func withTCPConfig(t *testing.T, fn func()) {
	t.Helper()
	port := *Conf.Port
	creds := *Conf.Credentials
	user := *Conf.Username
	pass := *Conf.Password
	telnet := *Conf.Telnet
	ssh := *Conf.SSH
	defer func() {
		*Conf.Port = port
		*Conf.Credentials = creds
		*Conf.Username = user
		*Conf.Password = pass
		*Conf.Telnet = telnet
		*Conf.SSH = ssh
	}()
	fn()
}

func TestGetIPPortUsesInlinePortAndCLIOverride(t *testing.T) {
	withTCPConfig(t, func() {
		*Conf.Port = ""
		ip, port := GetIPPort("192.0.2.10:2222")
		if ip != "192.0.2.10" || port != "2222" {
			t.Fatalf("expected inline host:port, got ip=%q port=%q", ip, port)
		}

		*Conf.Port = "23"
		ip, port = GetIPPort("192.0.2.10:2222")
		if ip != "192.0.2.10" || port != "23" {
			t.Fatalf("expected CLI port override, got ip=%q port=%q", ip, port)
		}
	})
}

func TestGetCredsPrecedenceAndColonPasswords(t *testing.T) {
	withTCPConfig(t, func() {
		*Conf.Credentials = ""
		*Conf.Username = ""
		*Conf.Password = ""
		user, pass := GetCreds("admin:p:a:s:s")
		if user != "admin" || pass != "p:a:s:s" {
			t.Fatalf("expected SplitN credential parsing, got user=%q pass=%q", user, pass)
		}

		*Conf.Credentials = "root:toor"
		user, pass = GetCreds("admin:ignored")
		if user != "root" || pass != "toor" {
			t.Fatalf("expected -creds precedence, got user=%q pass=%q", user, pass)
		}

		*Conf.Username = "cli"
		*Conf.Password = "secret"
		user, pass = GetCreds("admin:ignored")
		if user != "cli" || pass != "secret" {
			t.Fatalf("expected -user/-pass precedence, got user=%q pass=%q", user, pass)
		}
	})
}

func TestHTTPDefaults(t *testing.T) {
	if got := httpPort("http", ""); got != "80" {
		t.Fatalf("expected default HTTP port 80, got %q", got)
	}
	if got := httpPort("https", ""); got != "443" {
		t.Fatalf("expected default HTTPS port 443, got %q", got)
	}
	if got := httpPort("https", "8443"); got != "8443" {
		t.Fatalf("expected explicit port to win, got %q", got)
	}
	if got := httpPath("login.html"); got != "/login.html" {
		t.Fatalf("expected path to be normalized, got %q", got)
	}
	if got := httpPath("/admin"); got != "/admin" {
		t.Fatalf("expected existing leading slash to be preserved, got %q", got)
	}
}

func TestProtocolPreflightUsesParsedPort(t *testing.T) {
	withTCPConfig(t, func() {
		*Conf.Port = ""
		*Conf.Telnet = true
		if !shouldTelnet("2222") {
			t.Fatal("expected -telnet with parsed inline port to run")
		}

		*Conf.Telnet = false
		if !shouldTelnet("23") {
			t.Fatal("expected parsed standard telnet port to run")
		}

		*Conf.SSH = true
		ok, port := shouldSSH(IPinfo{hostString: "192.0.2.10:2222", port: "2222"})
		if !ok || port != "2222" {
			t.Fatalf("expected SSH preflight to keep parsed inline port, got ok=%v port=%q", ok, port)
		}

		*Conf.SSH = false
		ok, _ = shouldSSH(IPinfo{hostString: "192.0.2.10:2222", port: "2222"})
		if ok {
			t.Fatal("expected non-standard SSH port to require -ssh")
		}
	})
}

func TestSSHPreflightClosesProgressChannel(t *testing.T) {
	ipInfo := make(chan IPinfo)
	close(ipInfo)
	chSuccess := make(chan int)

	go SSHPreflight(chSuccess, ipInfo)

	if _, ok := <-chSuccess; ok {
		t.Fatal("expected SSHPreflight to close progress channel when there is no input")
	}
}
