package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const DefaultTimeout = 20 * time.Second

// SSOClient struct to hold SSO credentials and token
type SSOClient struct {
	ClientID     string
	ClientSecret string
	Token        string
	TokenURL     string
	TokenExpiry  time.Time
	sync.Mutex
	Transport http.RoundTripper
}

// TokenResponse struct to unmarshal the token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// NewSSOClient creates a new SSOClient
func NewSSOClient(clientID, clientSecret, tokenURL string) *SSOClient {
	return &SSOClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Transport:    http.DefaultTransport,
	}
}

// RoundTrip implements the RoundTripper interface
func (c *SSOClient) RoundTrip(req *http.Request) (*http.Response, error) {
	c.Lock()
	defer c.Unlock()

	if time.Now().After(c.TokenExpiry) {
		if err := c.authenticate(); err != nil {
			return nil, err
		}
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	return c.Transport.RoundTrip(req)
}

// authenticate authenticates to the SSO server and retrieves a token
func (c *SSOClient) authenticate() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", c.TokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("authenticate - error creating request: %v", err)
	}

	req.SetBasicAuth(c.ClientID, c.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("authenticate - error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("authenticate - error reading response: %v", err)
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return fmt.Errorf("authenticate - error parsing JSON: %v", err)
	}

	c.Token = tokenResponse.AccessToken
	c.TokenExpiry = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	return nil
}

// newAuthenticatedClient returns an HTTP client with SSO authentication
func newAuthenticatedClient(clientID, clientSecret, tokenURL string, timeout time.Duration) *http.Client {
	ssoClient := NewSSOClient(clientID, clientSecret, tokenURL)
	return &http.Client{
		Transport: ssoClient,
		Timeout:   timeout,
	}
}

// GetConsoleDotClient returns an HTTP client with SSO authentication for the console.redhat.com APIs
func GetConsoleDotClient(_ context.Context, d *plugin.QueryData, timeout time.Duration) (*http.Client, error) {
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "crc"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*http.Client), nil
	}

	// Default to the env var settings
	baseUrl := os.Getenv("CRC_URL")
	tokenURL := os.Getenv("CRC_TOKEN_URL")
	clientID := os.Getenv("CRC_CLIENT_ID")
	clientSecret := os.Getenv("CRC_CLIENT_SECRET")

	// Prefer config options given in Steampipe
	crcConfig := GetConfig(d.Connection)

	if crcConfig.BaseUrl != nil {
		baseUrl = *crcConfig.BaseUrl
	}
	if crcConfig.TokenURL != nil {
		tokenURL = *crcConfig.TokenURL
	}
	if crcConfig.ClientID != nil {
		clientID = *crcConfig.ClientID
	}
	if crcConfig.ClientSecret != nil {
		clientSecret = *crcConfig.ClientSecret
	}

	if baseUrl == "" {
		return nil, errors.New("'base_url' must be set in the connection configuration")
	}
	if tokenURL == "" {
		return nil, errors.New("'token_url' must be set in the connection configuration")
	}
	if clientID == "" {
		return nil, errors.New("'client_id' must be set in the connection configuration")
	}
	if clientSecret == "" {
		return nil, errors.New("'client_secret' must be set in the connection configuration")
	}

	client := newAuthenticatedClient(clientID, clientSecret, tokenURL, timeout)

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, client)

	// Done
	return client, nil
}
