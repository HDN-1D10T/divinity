package shodan

const BaseURL = "https://api.shodan.io"

type Client struct {
	apiKey  string
	baseURL string
}

func New(apiKey string) *Client {
	return newClient(apiKey, BaseURL)
}

func newClient(apiKey, baseURL string) *Client {
	return &Client{apiKey: apiKey, baseURL: baseURL}
}
