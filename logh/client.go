// Copyright 2018 Sean.ZH

package logh

import (
	"io"
	"net/http"
	"time"
)

// Cli package a http client and baseurl
type Cli struct {
	http.Client
	baseURL string
}

// New create an cli object
func New(url string, sec int) *Cli {
	return &Cli{
		http.Client{
			Timeout: time.Duration(sec) * time.Second,
		},
		url,
	}
}

// SetBaseURL change base url for the client
func (c *Cli) SetBaseURL(u string) {
	c.baseURL = u
}

// GetBaseURL returns current baseUrl
func (c *Cli) GetBaseURL() string {
	return c.baseURL
}

// Get do a get for client
func (c *Cli) Get(u string) ([]byte, error) {
	req, _ := http.NewRequest("GET", c.baseURL+u, nil)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
