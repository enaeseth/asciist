// Package client implements communication with the asciist service.
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

// A Client provides convenient interaction with the asciist service.
type Client struct {
	URL  string
	HTTP *http.Client
}

// New creates a new Client for a service rooted at the given url.
// It uses the default net/http Client to talk to the service.
func New(url string) *Client {
	return NewWithHTTP(url, http.DefaultClient)
}

// NewWithHTTP creates a new Client for a service rooted at the given url,
// using the given http Client for communication.
func NewWithHTTP(url string, httpClient *http.Client) *Client {
	return &Client{
		URL:  url,
		HTTP: httpClient,
	}
}

// Convert asks the service to create ASCII art based on the image whose
// contents can be read out of the given reader.
// The art returned will be `width` characters wide, with a height chosen
// by the service to match the input image's aspect ratio when drawn with
// typical monospace fonts at a typical line height.
func (c *Client) Convert(reader io.Reader, width uint) (string, error) {
	image, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return c.ConvertBytes(image, width)
}

// ConvertFile asks the service to create ASCII art based on the image
// with the given filename.
// The art returned will be `width` characters wide, with a height chosen
// by the service to match the input image's aspect ratio when drawn with
// typical monospace fonts at a typical line height.
func (c *Client) ConvertFile(filename string, width uint) (string, error) {
	image, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return c.ConvertBytes(image, width)
}

// ConvertBytes asks the service to create ASCII art based on the image
// whose contents are contained in the given byte slice.
// The art returned will be `width` characters wide, with a height chosen
// by the service to match the input image's aspect ratio when drawn with
// typical monospace fonts at a typical line height.
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
