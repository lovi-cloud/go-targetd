package targetd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Client is client of go-targetd
type Client struct {
	User     string
	Password string

	URL        *url.URL
	HTTPClient *http.Client
	Logger     *log.Logger
}

var (
	userAgent = fmt.Sprintf("go-targetd")
)

// New create go-targetd client
func New(baseURL, username, password string, httpClient *http.Client, logger *log.Logger) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse baseURL: %w", err)
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if logger == nil {
		l := log.New(os.Stdout, "", log.LstdFlags)
		logger = l
	}

	return &Client{
		User:     username,
		Password: password,

		URL:        u,
		HTTPClient: httpClient,
		Logger:     logger,
	}, nil
}
