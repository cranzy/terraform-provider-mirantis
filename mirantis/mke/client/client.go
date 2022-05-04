package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
)

const (
	EndpointDefaultScheme = "https"
)

// Client MSR client
type Client struct {
	apiURL     *url.URL
	auth       *Auth
	HTTPClient *http.Client
}

// NewClient from a string URL and u/p
func NewClientSimple(endpoint, username, password string) (Client, error) {
	HTTPClient := &http.Client{}
	auth := NewAuthUP(username, password)

	apiURL, err := url.Parse(endpoint)
	if err != nil {
		return Client{}, err
	}

	return NewClient(apiURL, &auth, HTTPClient)
}

// NewUnsafeSSLClient that allows self-signed SSL from a string URL and u/p
func NewUnsafeSSLClient(endpoint, username, password string) (Client, error) {
	HTTPClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	auth := NewAuthUP(username, password)

	apiURL, err := url.Parse(endpoint)
	if err != nil {
		return Client{}, err
	}

	return NewClient(apiURL, &auth, HTTPClient)
}

// NewUnsafeSSLClient creates a new MKE API Client that ignores unsafe SSL
func NewClient(apiURL *url.URL, auth *Auth, HTTPClient *http.Client) (Client, error) {
	return Client{
		apiURL:     apiURL,
		HTTPClient: HTTPClient,
		auth:       auth,
	}, nil
}

// Build a request URL string from the client endpoint and an API target path
func (c *Client) reqURLFromTarget(target string) string {
	// target should be a relative path, and will be treated as a relative reference
	// to the client URL
	// @see https://pkg.go.dev/net/url#URL.ResolveReference
	targetURL, _ := url.Parse(target)
	relativeURL := c.apiURL.ResolveReference(targetURL)

	return relativeURL.String()
}
