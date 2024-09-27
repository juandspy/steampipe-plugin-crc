package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// MakeAPIRequest makes an API request to the specified endpoint
func MakeAPIRequest(ctx context.Context, d *plugin.QueryData, method, endpoint string, body interface{}, timeout time.Duration) (*http.Response, error) {
	client, err := GetConsoleDotClient(ctx, d, timeout)
	if err != nil {
		return nil, err
	}

	baseURL := GetConfig(d.Connection).BaseUrl
	url := *baseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code %d and body: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}
