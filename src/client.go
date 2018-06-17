package jumphelper

import (
	"io/ioutil"
	"net"
	//"os"
)

// Client is a TCP client that responds to jumphelper requests
type Client struct {
	host string
	port string

	serverAddress *net.TCPAddr
	connection    net.Conn
	err           error
}

func (c *Client) address() string {
	return c.host + ":" + c.port
}

// Echo writes a request to a jumphelper server
func (c *Client) Echo(s string) ([]byte, error) {
	c.connection.Write([]byte(s))
	if c.err != nil {
		return nil, c.err
	}
	return ioutil.ReadAll(c.connection)
	//return nil
}

// NewClient creates a new jumphelper client
func NewClient(Host, Port string) (*Client, error) {
	var c Client
	c.host = Host
	c.port = Port
	c.serverAddress, c.err = net.ResolveTCPAddr("tcp", c.address())
	if c.err != nil {
		return nil, c.err
	}
	c.connection, c.err = net.DialTCP("tcp", nil, c.serverAddress)
	if c.err != nil {
		return nil, c.err
	}
	return &c, nil
}
