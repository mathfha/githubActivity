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

func (c *Client) FetchUserActivity(ctx context.Context, username string) ([]Event, error) {
	url := fmt.Sprintf("%s/users/%s/events", baseURL, username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "GitHub-Activity-CLI")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal mengirim HTTP request: %w", err)
	}
	
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var events []Event
		if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
			return nil, fmt.Errorf("gagal mem-parsing respons JSON: %w", err)
		}
		return events, nil

	case http.StatusNotFound:
		return nil, fmt.Errorf("pengguna '%s' tidak ditemukan (404)", username)
		
	case http.StatusForbidden:
		return nil, fmt.Errorf("akses ditolak (403), kemungkinan batas *rate limit* API GitHub telah tercapai")
		
	default:
		return nil, fmt.Errorf("API mengembalikan status yang tidak diharapkan: %d", resp.StatusCode)
	}
}
