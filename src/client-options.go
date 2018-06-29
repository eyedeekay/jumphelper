package jumphelper

import (
	"fmt"
	"strconv"
)

// ClientOption is a server option
type ClientOption func(*Client) error

//SetClientHost sets the host of the Server to connect to
func SetClientHost(s string) func(*Client) error {
	return func(c *Client) error {
		c.host = s
		return nil
	}
}

//SetClientPort sets the port of the Server to connect to
func SetClientPort(s string) func(*Client) error {
	return func(c *Client) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Invalid port; non-number")
		}
		if port < 65536 && port > -1 {
			c.port = s
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetClientPortInt sets the port of the Server to connect to with an int
func SetClientPortInt(s int) func(*Client) error {
	return func(c *Client) error {
		if s < 65536 && s > -1 {
			c.port = strconv.Itoa(s)
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetClientVerbose sets the verbosity of the client
func SetClientVerbose(s bool) func(*Client) error {
	return func(c *Client) error {
		c.verbose = s
		return nil
	}
}
