package jumphelper

import (
	"fmt"
	"strconv"
)

// JumpHelperOption is a jumphelper option
type JumpHelperOption func(*JumpHelper) error

//SetJumpHelperAddressBookPath sets the host of the JumpHelper client's SAM bridge
func SetJumpHelperAddressBookPath(s string) func(*JumpHelper) error {
	return func(c *JumpHelper) error {
		c.addressBookPath = s
		return nil
	}
}

//SetJumpHelperHost sets the host of the JumpHelper client's SAM bridge
func SetJumpHelperHost(s string) func(*JumpHelper) error {
	return func(c *JumpHelper) error {
		c.samHost = s
		return nil
	}
}

//SetJumpHelperPort sets the port of the JumpHelper client's SAM bridge
func SetJumpHelperPort(s string) func(*JumpHelper) error {
	return func(c *JumpHelper) error {
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

//SetJumpHelperPortInt sets the port of the JumpHelper client's SAM bridge with an int
func SetJumpHelperPortInt(s int) func(*JumpHelper) error {
	return func(c *JumpHelper) error {
		if s < 65536 && s > -1 {
			c.samPort = strconv.Itoa(s)
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetJumpHelperUseHelper sets the host of the JumpHelper client's SAM bridge
func SetJumpHelperUseHelper(s bool) func(*JumpHelper) error {
	return func(c *JumpHelper) error {
		c.ext = s
		return nil
	}
}