package jumphelper

import (
	"fmt"
	"io/ioutil"
    "net/url"
	"strings"
)

// JumpHelper is a struct that prioritizes i2p address sources
type JumpHelper struct {
	addressBookPath string
	addressBook     []string
}

// LoadAddressBook loads an addressbook in csv(name, b32) format
func (j *JumpHelper) LoadAddressBook() error {
	content, err := ioutil.ReadFile(j.addressBookPath)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	j.addressBook = strings.Split(string(content), "\n")
	return nil
}

func (j *JumpHelper) trim(k string) string {
	r := strings.Replace(k, "http://", "", -1)
	if strings.HasSuffix(k, ".i2p") {
		return r
	}
	r = strings.SplitN(r, ".i2p", -1)[0]
	return strings.TrimSpace(strings.TrimSuffix(r, "/"))
}

// SearchAddressBook finds a (name, b32) pair in the addressbook, or returns nil of one is not found
func (j *JumpHelper) SearchAddressBook(pk string) []string {
	k, e := url.Parse(pk)
    if e != nil {
        return nil
    }
	for _, a := range j.addressBook {
		r := strings.SplitN(a, ",", -1)
		if r[0] == j.trim(k.Host) {
			return r
		}
	}
	return nil
}

// CheckAddressBook returns true if an address is present, false if not
func (j *JumpHelper) CheckAddressBook(pk string) bool {
	k, e := url.Parse(pk)
    if e != nil {
        return false
    }
	for _, a := range j.addressBook {
		r := strings.SplitN(a, ",", -1)[0]
		if r == j.trim(k.Host) {
			return true
		}
	}
	return false
}

// NewJumpHelper creates a new JumpHelper object
func NewJumpHelper(addressBookPath string) (*JumpHelper, error) {
	var j JumpHelper
	j.addressBookPath = addressBookPath
	err := j.LoadAddressBook()
	return &j, err
}
