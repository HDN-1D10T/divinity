package shodan

import (
	"net/url"
	"testing"
)

func TestHostSearchURLEncodesQueryAndPage(t *testing.T) {
	client := New("test-key")
	rawURL := client.hostSearchURL(`Device Manufacturer A port:80`, 5)

	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("hostSearchURL returned invalid URL: %v", err)
	}

	values := parsed.Query()
	if values.Get("key") != "test-key" {
		t.Fatalf("expected API key query parameter, got %q", values.Get("key"))
	}
	if values.Get("query") != "Device Manufacturer A port:80" {
		t.Fatalf("expected unescaped Shodan query after parsing, got %q", values.Get("query"))
	}
	if values.Get("page") != "5" {
		t.Fatalf("expected page query parameter, got %q", values.Get("page"))
	}
}
