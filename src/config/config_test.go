package config

import "testing"

func testString(value string) *string {
	return &value
}

func testBool(value bool) *bool {
	return &value
}

func testInt(value int) *int {
	return &value
}

func testOptions() Options {
	return Options{
		Alert:       testString("SUCCESS"),
		All:         testBool(false),
		ASN:         testString(""),
		BasicAuth:   testString(""),
		Cidr:        testString(""),
		ContentType: testString(""),
		Credentials: testString(""),
		Data:        testString(""),
		List:        testString(""),
		ListIPs:     testBool(false),
		HeaderName:  testString(""),
		HeaderValue: testString(""),
		HTTPTimeout: testInt(10000),
		IPOnly:      testBool(false),
		Masscan:     testBool(false),
		Method:      testString("GET"),
		OutputFile:  testString(""),
		Path:        testString("/"),
		Pages:       testInt(1),
		Passive:     testBool(false),
		Password:    testString(""),
		Port:        testString(""),
		Protocol:    testString(""),
		Routes:      testBool(false),
		SearchTerm:  testString(""),
		Scan:        testBool(false),
		ScanFast:    testBool(false),
		SSH:         testBool(false),
		Success:     testString(""),
		TopPorts:    testBool(false),
		Telnet:      testBool(false),
		Timeout:     testInt(500),
		Username:    testString(""),
	}
}

func TestApplyJSONPreservesOptionPointers(t *testing.T) {
	opts := testOptions()
	originalPort := opts.Port
	originalScanFast := opts.ScanFast
	originalPages := opts.Pages
	originalHTTPTimeout := opts.HTTPTimeout

	err := opts.applyJSON([]byte(`{"port":"8080","scanfast":true,"pages":4,"method":"POST","http-timeout":1500}`), "test")
	if err != nil {
		t.Fatalf("applyJSON returned error: %v", err)
	}

	if opts.Port != originalPort {
		t.Fatal("port pointer was replaced instead of updated")
	}
	if opts.ScanFast != originalScanFast {
		t.Fatal("scanfast pointer was replaced instead of updated")
	}
	if opts.Pages != originalPages {
		t.Fatal("pages pointer was replaced instead of updated")
	}
	if opts.HTTPTimeout != originalHTTPTimeout {
		t.Fatal("http-timeout pointer was replaced instead of updated")
	}
	if *opts.Port != "8080" || !*opts.ScanFast || *opts.Pages != 4 || *opts.Method != "POST" || *opts.HTTPTimeout != 1500 {
		t.Fatalf("unexpected parsed config: port=%q scanfast=%v pages=%d method=%q http-timeout=%d", *opts.Port, *opts.ScanFast, *opts.Pages, *opts.Method, *opts.HTTPTimeout)
	}
}

func TestCLIOverridesReplaceJSONValues(t *testing.T) {
	opts := testOptions()
	err := opts.applyJSON([]byte(`{"port":"80","scanfast":false,"pages":1}`), "test")
	if err != nil {
		t.Fatalf("applyJSON returned error: %v", err)
	}

	err = opts.applyCLIOverrides(map[string]string{
		"port":         "443",
		"scanfast":     "true",
		"pages":        "7",
		"http-timeout": "3000",
	})
	if err != nil {
		t.Fatalf("applyCLIOverrides returned error: %v", err)
	}

	if *opts.Port != "443" {
		t.Fatalf("expected CLI port override, got %q", *opts.Port)
	}
	if !*opts.ScanFast {
		t.Fatal("expected CLI scanfast override")
	}
	if *opts.Pages != 7 {
		t.Fatalf("expected CLI pages override, got %d", *opts.Pages)
	}
	if *opts.HTTPTimeout != 3000 {
		t.Fatalf("expected CLI http-timeout override, got %d", *opts.HTTPTimeout)
	}
}

func TestInvalidBooleanOverrideReturnsError(t *testing.T) {
	opts := testOptions()
	err := opts.applyCLIOverrides(map[string]string{"scan": "definitely"})
	if err == nil {
		t.Fatal("expected invalid boolean override to return an error")
	}
}
