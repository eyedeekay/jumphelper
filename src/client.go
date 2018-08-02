package jumphelper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

import "log"

// Client is a HTTP client that makes jumphelper requests
type Client struct {
	host    string
	port    string
	verbose bool

	client *http.Client
}

// Log wraps Println to control verbosity.
func (c *Client) Log(s string) string {
	if c.verbose {
		log.Println(s)
	}
	return s
}

func (c *Client) address(s string, m ...string) string {
	if len(m) > 0 {
		u := "http://" + c.host + ":" + c.port + "/" + m[0] + "/" + s + "/"
		return c.Log(u)
	}
	u := "http://" + c.host + ":" + c.port + "/" + s + "/"
	return c.Log(u)
}

// Check writes a request for a true-false answer to a jumphelper server
func (c *Client) Check(s string) (bool, error) {
	resp, err := c.client.Get(c.address(s, "check"))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	sbytes := strings.TrimSpace(string(bytes))
	c.Log("Log: " + sbytes)
	if strings.HasPrefix(sbytes, "TRUE") {
		return true, nil
	}
	return false, nil
}

// Request writes a request for a base32 answer to a jumphelper server
func (c *Client) Request(s string) (string, error) {
	resp, err := c.client.Get(c.address(s, "request"))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	c.Log(resp.Header.Get("Location"))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Jump writes a request for a base64 address to a jumphelper server
func (c *Client) Jump(s string) (string, error) {
	resp, err := c.client.Get(c.address(s, "jump"))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	c.Log(resp.Header.Get("Location"))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
// NewClient creates a new jumphelper client
func NewClient(Host, Port string, verbose bool) (*Client, error) {
	return NewClientFromOptions(SetClientHost(Host), SetClientPort(Port), SetClientVerbose(verbose))
}

// NewClientFromOptions creates a new jumphelper client from functional arguments
func NewClientFromOptions(opts ...func(*Client) error) (*Client, error) {
	var c Client
	c.client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	c.host = "127.0.0.1"
	c.port = "7854"
	c.verbose = false
	for _, o := range opts {
		if err := o(&c); err != nil {
			return nil, fmt.Errorf("Client configuration error: %s", err)
		}
	}
	return &c, nil
}
