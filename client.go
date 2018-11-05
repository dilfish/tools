package tools

import (
	"net/http"
	"time"
    "io/ioutil"
)

type client struct {
	http.Client
	baseURL string
}

type User struct {
	name  string
	email string
}

func New(url string) *client {
	return &client{
		http.Client{
			Timeout: time.Duration(1) * time.Second,
		},
		url,
	}
}


func (c *client) Get (u string) ([]byte, error) {
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
