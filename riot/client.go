package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    map[string]string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: map[string]string{
			"summoner": "https://%s.api.riotgames.com",
			"match":    "https://%s.api.riotgames.com",
			"league":   "https://%s.api.riotgames.com",
		},
	}
}

func (c *Client) makeRequest(ctx context.Context, method, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("X-Riot-Token", c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "TFT-Arena/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp.StatusCode, body)
	}

	return body, nil
}

func (c *Client) handleErrorResponse(statusCode int, body []byte) error {
	var riotErr RiotAPIError
	if err := json.Unmarshal(body, &riotErr); err != nil {
		return NewRiotError(statusCode, string(body))
	}
	return NewRiotError(statusCode, riotErr.Status.Message)
}

func (c *Client) getRegionURL(service, region string) string {
	return fmt.Sprintf(c.baseURL[service], region)
}