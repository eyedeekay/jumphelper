package jumphelper

import (
	"fmt"
	"strconv"
)

// ServerOption is a server option
type ServerOption func(*Server) error

//SetServerAddressBookPath sets the host of the Server client's SAM bridge
func SetServerAddressBookPath(s string) func(*Server) error {
	return func(c *Server) error {
		c.addressBookPath = s
		return nil
	}
}

//SetServerHost sets the host of the Server client's SAM bridge
func SetServerHost(s string) func(*Server) error {
	return func(c *Server) error {
		c.host = s
		return nil
	}
}

//SetServerPort sets the port of the Server client's SAM bridge
func SetServerPort(s string) func(*Server) error {
	return func(c *Server) error {
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

//SetServerPortInt sets the port of the Server client's SAM bridge with an int
func SetServerPortInt(s int) func(*Server) error {
	return func(c *Server) error {
		if s < 65536 && s > -1 {
			c.port = strconv.Itoa(s)
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetServerRate sets the host of the Server client's SAM bridge
func SetServerRate(s int) func(*Server) error {
	return func(c *Server) error {
		c.rate = s
		return nil
	}
}

//SetServerBurst sets the host of the Server client's SAM bridge
func SetServerBurst(s int) func(*Server) error {
	return func(c *Server) error {
		c.burst = s
		return nil
	}
}

//SetServerUseHelper sets the host of the Server client's SAM bridge
func SetServerUseHelper(s bool) func(*Server) error {
	return func(c *Server) error {
		c.ext = s
		return nil
	}
}
