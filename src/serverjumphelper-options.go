package jumphelper

import (
	"fmt"
	"strconv"
)

//SetServerJumpHelperHost sets the host of the Server client's SAM bridge
func SetServerJumpHelperHost(s string) func(*Server) error {
	return func(c *Server) error {
		c.samHost = s
		return nil
	}
}

//SetServerJumpHelperPort sets the port of the Server client's SAM bridge
func SetServerJumpHelperPort(s string) func(*Server) error {
	return func(c *Server) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Invalid port; non-number")
		}
		if port < 65536 && port > -1 {
			c.samPort = s
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetServerJumpHelperPortInt sets the port of the Server client's SAM bridge with an int
func SetServerJumpHelperPortInt(s int) func(*Server) error {
	return func(c *Server) error {
		if s < 65536 && s > -1 {
			c.samPort = strconv.Itoa(s)
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}
