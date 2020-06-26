package client

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Client struct {
	Url *string
}

func (c *Client) request() (*map[string]interface{}, error) {
	log.Println("making request to", *c.Url)

	result := make(map[string]interface{})

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cl := &http.Client{Transport: tr}
	resp, err := cl.Get(*c.Url)

	if err != nil || resp.StatusCode != 200 {
		return nil, errors.New("failed to get values")
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	return &result, err
}

func (c *Client) Value() (*map[string]interface{}, error) {
	return c.request()
}

func NewClient() *Client {
	return &Client{}
}
