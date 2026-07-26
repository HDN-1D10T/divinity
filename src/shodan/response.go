package shodan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type apiError struct {
	Error string `json:"error"`
}

func decodeResponse(operation string, res *http.Response, target interface{}) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("%s failed: error reading response: %w", operation, err)
	}

	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		return fmt.Errorf("%s failed: empty response body from Shodan", operation)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("%s failed: HTTP %s%s", operation, res.Status, responseDetail(body))
	}

	var shodanErr apiError
	if err := json.Unmarshal(body, &shodanErr); err == nil && shodanErr.Error != "" {
		return fmt.Errorf("%s failed: %s", operation, shodanErr.Error)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("%s failed: invalid JSON response: %w", operation, err)
	}
	return nil
}

func responseDetail(body []byte) string {
	var shodanErr apiError
	if err := json.Unmarshal(body, &shodanErr); err == nil && shodanErr.Error != "" {
		return ": " + shodanErr.Error
	}

	detail := strings.TrimSpace(string(body))
	if detail == "" {
		return ""
	}
	if len(detail) > 240 {
		detail = detail[:240] + "..."
	}
	return ": " + detail
}
