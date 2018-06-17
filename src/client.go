package jumphelper

import (
	"io/ioutil"
	"net"
	"strings"
)

// Client is a HTTP client that makes jumphelper requests
type Client struct {
	host string
	port string

	connection    net.Conn
}

func (c *Client) address() string {
	return c.host + ":" + c.port
}

// Check writes a request for a true-false answer to a jumphelper server
func (c *Client) Check(s string) (bool, error) {
	c.connection.Write([]byte(s))
    bytes, err := ioutil.ReadAll(c.connection)
	if err != nil {
		return false, err
	}
    if strings.HasPrefix("TRUE", string(bytes)){
        return true, nil
    }
	return false, nil
}

// Request writes a request for a base32 answer to a jumphelper server
func (c *Client) Request(s string) (string, error) {
	c.connection.Write([]byte(s))
    bytes, err := ioutil.ReadAll(c.connection)
	if err != nil {
		return nil, err
	}
	return string(bytes)
}

// NewClient creates a new jumphelper client
func NewClient(Host, Port string) (*Client, error) {
	var c Client
	c.host = Host
	c.port = Port
	return &c, nil
}
