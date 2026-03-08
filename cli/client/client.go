package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/formatho/agent-todo/cli/config"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func New() *Client {
	cfg := config.Get()
	return &Client{
		BaseURL: cfg.ServerURL,
		HTTPClient: &http.Client{
			// TODO: Add timeout configuration
		},
	}
}

func (c *Client) doRequest(method, path string, body interface{}, token string, apiKey string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}

	return c.HTTPClient.Do(req)
}

func (c *Client) Get(path string, useAuth bool) (*http.Response, error) {
	cfg := config.Get()
	token := ""
	apiKey := ""

	if useAuth {
		token = cfg.Token
		apiKey = cfg.APIKey
	}

	return c.doRequest("GET", path, nil, token, apiKey)
}

func (c *Client) Post(path string, body interface{}, useAuth bool) (*http.Response, error) {
	cfg := config.Get()
	token := ""
	apiKey := ""

	if useAuth {
		token = cfg.Token
		apiKey = cfg.APIKey
	}

	return c.doRequest("POST", path, body, token, apiKey)
}

func (c *Client) Patch(path string, body interface{}, useAuth bool) (*http.Response, error) {
	cfg := config.Get()
	token := ""
	apiKey := ""

	if useAuth {
		token = cfg.Token
		apiKey = cfg.APIKey
	}

	return c.doRequest("PATCH", path, body, token, apiKey)
}

func (c *Client) Delete(path string, useAuth bool) (*http.Response, error) {
	cfg := config.Get()
	token := ""
	apiKey := ""

	if useAuth {
		token = cfg.Token
		apiKey = cfg.APIKey
	}

	return c.doRequest("DELETE", path, nil, token, apiKey)
}
