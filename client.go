package lcutil

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

type Client struct {
	httpClient *http.Client
	token      string
	Port       string
	Path       string
}

// Both operating systems produce an output where we can find the important pieces for Client
func newClient(output []byte) (*Client, error) {
	ports := regexp.MustCompile(`--app-port=([0-9]*)`).FindAllSubmatch(output, 1)
	paths := regexp.MustCompile(`--install-directory=([\w//-_]*)`).FindAllSubmatch(output, 1)
	tokens := regexp.MustCompile(`--remoting-auth-token=([\w-_]*)`).FindAllSubmatch(output, 1)

	if len(ports) < 0 && len(tokens) < 0 {
		return &Client{}, NotRunningErr
	}

	port := string(ports[0][1])
	token := string(tokens[0][1])
	path := string(paths[0][1])

	return &Client{
		httpClient: newHttpClient(),
		token:      token,
		Port:       port, Path: path}, nil
}

// URL returns a url.URL that you can edit further.
func (c *Client) URL(uri string) (u url.URL, err error) {
	urlp, err := url.Parse(fmt.Sprintf("https://127.0.0.1:%s%s", c.Port, uri))
	if err == nil {
		u = *urlp
	}

	return u, err
}

func (c *Client) NewRequest(req_type string, u url.URL, form []byte) (*http.Request, error) {
	req, err := DefaultNewRequest(req_type, u, form)
	if err != nil {
		return req, err
	}

	req.SetBasicAuth(`riot`, c.token)
	return req, nil

}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *Client) Get(u url.URL) (*http.Response, error) {
	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return &http.Response{}, err
	}

	return c.httpClient.Do(req)
}

func (c *Client) Post(u url.URL, data []byte) (*http.Response, error) {
	req, err := c.NewRequest("POST", u, data)
	if err != nil {
		return &http.Response{}, err
	}

	return c.httpClient.Do(req)
}

// Basic way to create a request to most APIs
func DefaultNewRequest(req_type string, u url.URL, data []byte) (req *http.Request, err error) {
	raw := u.String()

	if data != nil {
		req, err = http.NewRequest(req_type, raw, bytes.NewReader(data))
		req.Header.Add("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(req_type, raw, nil)
	}

	if err != nil {
		return &http.Request{}, err
	}

	return req, err

}

var (
	DownloadFailedErr error = fmt.Errorf("Failed to download file.")
	NotRunningErr     error = errors.New("League of legends is not currently running!")
)
