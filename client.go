// Copyright 2018 Sean.ZH

package tools

import (
	"net/http"
	"time"
    "io/ioutil"
)

// Cli package a http client and baseurl
type Cli struct {
	http.Client
	baseURL string
}

// New create an cli object
func New(url string) *Cli {
	return &Cli {
		http.Client{
			Timeout: time.Duration(1) * time.Second,
		},
		url,
	}
}


// Get do a get for client
func (c *Cli) Get (u string) ([]byte, error) {
    req, err := http.NewRequest("GET", c.baseURL + u, nil)
    if err != nil {
        return nil, err
    }
    resp, err := c.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)
}
