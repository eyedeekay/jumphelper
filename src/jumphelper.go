package jumphelper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cryptix/gosam"
)

// JumpHelper is a struct that prioritizes i2p address sources
type JumpHelper struct {
	addressBookPath string

	samHost string
	samPort string

	samBridgeConn *goSam.Client
	ext           bool

	tr     *http.Transport
	client *http.Client

	addressBook       []string
	remoteAddressBook []string
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

// SyncRemoteAddressBooks syncs addressbooks from subscription services to the standalone addressbook
func (j *JumpHelper) SyncRemoteAddressBooks() error {
	//
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
	var kv string
	if !strings.HasPrefix(pk, "http://") {
		kv = "http://" + pk
	} else {
		kv = pk
	}
	k, e := url.Parse(kv)
	if e != nil {
		return nil
	}
	for _, a := range j.addressBook {
		r := strings.SplitN(a, ",", 2)
		if len(r) == 2 {
			if r[0] == j.trim(k.Host) {
				return r
			}
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
func NewJumpHelper(addressBookPath, host, port string) (*JumpHelper, error) {
	return NewJumpHelperFromOptions(SetJumpHelperAddressBookPath(addressBookPath))
}

// NewJumpHelperFromOptions creates a new JumpHelper object
func NewJumpHelperFromOptions(opts ...func(*JumpHelper) error) (*JumpHelper, error) {
	var j JumpHelper
	j.addressBookPath = "/var/lib/i2pd/addressbook/addresses.csv"
	j.samHost = "127.0.0.1"
	j.samPort = "7056"
	j.ext = false
	for _, o := range opts {
		if err := o(&j); err != nil {
			return nil, fmt.Errorf("Service configuration error: %s", err)
		}
	}
	err := j.LoadAddressBook()
	if err != nil {
		return nil, err
	}
	if j.ext {
		j.samBridgeConn, err = goSam.NewClient(j.samHost + ":" + j.samPort)
	}
	j.tr = &http.Transport{
		Dial: j.samBridgeConn.Dial,
	}
	j.client = &http.Client{Transport: j.tr}
	return &j, err
}

func printKvs(kv []string) {
	for i, s := range kv {
		fmt.Println("Key-value Pair", i, s)
	}
}
