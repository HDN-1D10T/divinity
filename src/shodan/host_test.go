package shodan

import (
	"io"
	"net/http"
	"net/url"
	"strings"
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

func TestHostSearchReturnsHelpfulErrorForEmptyBody(t *testing.T) {
	res := &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       io.NopCloser(strings.NewReader("")),
	}
	var ret HostSearch
	err := decodeResponse("Shodan host search", res, &ret)
	if err == nil {
		t.Fatal("expected empty response body to return an error")
	}
	if !strings.Contains(err.Error(), "empty response body") {
		t.Fatalf("expected helpful empty-body error, got %v", err)
	}
}

func TestHostSearchReturnsShodanHTTPErrorMessage(t *testing.T) {
	res := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Status:     "401 Unauthorized",
		Body:       io.NopCloser(strings.NewReader(`{"error":"Invalid API key"}`)),
	}
	var ret HostSearch
	err := decodeResponse("Shodan host search", res, &ret)
	if err == nil {
		t.Fatal("expected HTTP error response to return an error")
	}
	if !strings.Contains(err.Error(), "HTTP 401 Unauthorized") || !strings.Contains(err.Error(), "Invalid API key") {
		t.Fatalf("expected status and Shodan error in message, got %v", err)
	}
}
