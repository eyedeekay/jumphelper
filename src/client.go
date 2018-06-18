package jumphelper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client is a HTTP client that makes jumphelper requests
type Client struct {
	host string
	port string

	client *http.Client
}

func (c *Client) address(s string, m ...string) string {
	if len(m) > 0 {
		u := "http://" + c.host + ":" + c.port + "/" + m[0] + "/" + s
		return u
	}
	u := "http://" + c.host + ":" + c.port + "/" + s
	return u
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
	if strings.HasPrefix(sbytes, "TRUE") {
		fmt.Println(sbytes)
		return true, nil
	}
	return false, nil
}

// Request writes a request for a base32 answer to a jumphelper server
func (c *Client) Request(s string) (string, error) {
	resp, err := c.client.Get(c.address(s))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// NewClient creates a new jumphelper client
func NewClient(Host, Port string) (*Client, error) {
	var c Client
	c.client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	c.host = Host
	c.port = Port
	return &c, nil
}
