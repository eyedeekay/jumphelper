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
func SetServerPort(s interface{}) func(*Server) error {
	return func(c *Server) error {
		switch v := s.(type) {
		case string:
			port, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("Invalid port; non-number")
			}
			if port < 65536 && port > -1 {
				c.port = v
				return nil
			}
			return fmt.Errorf("Invalid port")
		case int:
			if v < 65536 && v > -1 {
				c.port = strconv.Itoa(v)
				return nil
			}
			return fmt.Errorf("Invalid port")
		default:
			return fmt.Errorf("Invalid port")
		}
	}
}
