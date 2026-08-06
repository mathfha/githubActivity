package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://api.github.com"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) doRequest(ctx context.Context, url string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "GitHub-Activity-CLI")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to parse JSON response: %w", err)
		}
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("user not found (404)")
	case http.StatusForbidden:
		return fmt.Errorf("access forbidden (403), GitHub API rate limit may have been exceeded")
	default:
		return fmt.Errorf("unexpected API status code: %d", resp.StatusCode)
	}
}

func (c *Client) FetchUserActivity(ctx context.Context, username string) ([]Event, error) {
	url := fmt.Sprintf("%s/users/%s/events/public?per_page=100", baseURL, username)
	var events []Event
	err := c.doRequest(ctx, url, &events)
	if err != nil {
		if err.Error() == "user not found (404)" {
			return nil, fmt.Errorf("user '%s' not found (404)", username)
		}
		return nil, err
	}
	return events, nil
}

func (c *Client) FetchUserRepositories(ctx context.Context, username string) ([]Repository, error) {
	url := fmt.Sprintf("%s/users/%s/repos?per_page=100&sort=updated", baseURL, username)
	var repos []Repository
	err := c.doRequest(ctx, url, &repos)
	if err != nil {
		if err.Error() == "user not found (404)" {
			return nil, fmt.Errorf("user '%s' not found (404)", username)
		}
		return nil, err
	}
	return repos, nil
}
