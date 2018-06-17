package jumphelper

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// Client is a HTTP client that makes jumphelper requests
type Client struct {
	host string
	port string

	connection net.Conn
}

func (c *Client) address(s string) string {
	return c.host + ":" + c.port + "/" + s
}

// Check writes a request for a true-false answer to a jumphelper server
func (c *Client) Check(s string) (bool, error) {
	resp, err := http.Get("http://" + c.address(s))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if strings.HasPrefix("TRUE", string(bytes)) {
		return true, nil
	}
	return false, nil
}

// Request writes a request for a base32 answer to a jumphelper server
func (c *Client) Request(s string) (string, error) {
	resp, err := http.Get(c.address(s))
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
	c.host = Host
	c.port = Port
	return &c, nil
}
