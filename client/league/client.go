package league

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/gabriel-ross/lcutil/client"
	cu "github.com/gabriel-ross/lcutil/client/clientutil"
)

type Client struct {
	token string
	Port  string
	Path  string
}

// Both operating systems produce an output where we can find the important pieces for Client
func newClient(output []byte) (client.Client, error) {
	ports := regexp.MustCompile(`--app-port=([0-9]*)`).FindAllSubmatch(output, 1)
	paths := regexp.MustCompile(`--install-directory=([\w//-_]*)`).FindAllSubmatch(output, 1)
	tokens := regexp.MustCompile(`--remoting-auth-token=([\w-_]*)`).FindAllSubmatch(output, 1)

	if len(ports) < 0 && len(tokens) < 0 {
		return &Client{}, NotRunningErr
	}

	port := string(ports[0][1])
	token := string(tokens[0][1])
	path := string(paths[0][1])

	return &Client{token: token, Port: port, Path: path}, nil
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
	req, err := client.DefaultNewRequest(req_type, u, form)
	if err != nil {
		return req, err
	}

	req.SetBasicAuth(`riot`, c.token)
	return req, nil

}

func (c *Client) Get(u url.URL) (*http.Response, error) {
	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return &http.Response{}, err
	}

	return cu.HttpClient.Do(req)
}

func (c *Client) Post(u url.URL, data []byte) (*http.Response, error) {
	req, err := c.NewRequest("POST", u, data)
	if err != nil {
		return &http.Response{}, err
	}

	return cu.HttpClient.Do(req)
}

var (
	DownloadFailedErr error = fmt.Errorf("Failed to download file.")
	NotRunningErr     error = errors.New("League of legends is not currently running!")
)
