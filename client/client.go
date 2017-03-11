package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/enaeseth/asciist/service"
)

type Client struct {
	URL  string
	HTTP *http.Client
}

func New(url string) *Client {
	return NewWithHTTP(url, http.DefaultClient)
}

func NewWithHTTP(url string, httpClient *http.Client) *Client {
	return &Client{
		URL:  url,
		HTTP: httpClient,
	}
}

func (c *Client) Convert(reader io.Reader, width uint) (string, error) {
	image, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return c.ConvertBytes(image, width)
}

func (c *Client) ConvertFile(filename string, width uint) (string, error) {
	image, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return c.ConvertBytes(image, width)
}

func (c *Client) ConvertBytes(image []byte, width uint) (string, error) {
	req := service.Request{
		Width: width,
		Image: image,
	}

	payload, err := json.Marshal(&req)
	if err != nil {
		return "", err
	}

	resp, err := c.HTTP.Post(c.URL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := resp.Body.Close(); err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		defaultErr := errors.New(resp.Status)

		var failure service.Failure
		if err := json.Unmarshal(body, &failure); err == nil {
			return "", errors.New(failure.Error)
		}

		return "", defaultErr
	}

	var success service.Success
	if err := json.Unmarshal(body, &success); err != nil {
		return "", err
	}

	return success.Art, nil
}
