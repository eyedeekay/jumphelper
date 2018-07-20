package jumphelper

import (
	"fmt"
	"strconv"
)

// addresslistOption is a addresslist option
type addresslistOption func(*addresslist) error

//SetaddresslistAddressBookPath sets the host of the addresslist client's SAM bridge
func SetaddresslistAddressBookPath(s string) func(*addresslist) error {
	return func(c *addresslist) error {
		c.addressBookURL = s
		return nil
	}
}

//SetaddresslistHost sets the host of the addresslist client's SAM bridge
func SetaddresslistHost(s string) func(*addresslist) error {
	return func(c *addresslist) error {
		c.samHost = s
		return nil
	}
}

//SetaddresslistPort sets the port of the addresslist client's SAM bridge
func SetaddresslistPort(s string) func(*addresslist) error {
	return func(c *addresslist) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Invalid port; non-number: %s", s)
		}
		if port < 65536 && port > -1 {
			c.samPort = s
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetaddresslistPortInt sets the port of the addresslist client's SAM bridge with an int
func SetaddresslistPortInt(s int) func(*addresslist) error {
	return func(c *addresslist) error {
		if s < 65536 && s > -1 {
			c.samPort = strconv.Itoa(s)
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}
