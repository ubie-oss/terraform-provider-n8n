package n8n

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// URL joins Endpoint with path segments (no query string).
func (c *Client) URL(parts ...string) string {
	base := strings.TrimSuffix(c.Endpoint, "/")
	for _, p := range parts {
		base += "/" + strings.Trim(p, "/")
	}
	return base
}

// DoJSON performs an HTTP JSON request. 404 responses become NotFoundError.
func (c *Client) DoJSON(ctx context.Context, method, rawURL string, body any, out any, notFoundResource, notFoundID string) error {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, reader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		id := notFoundID
		if id == "" {
			id = rawURL
		}
		resource := notFoundResource
		if resource == "" {
			resource = "resource"
		}
		return &NotFoundError{Resource: resource, ID: id}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Body: strings.TrimSpace(string(respBody))}
	}

	if out == nil || len(respBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
