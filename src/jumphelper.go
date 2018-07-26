package jumphelper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"github.com/eyedeekay/i2pasta/convert"
)

// JumpHelper is a struct that prioritizes i2p address sources
type JumpHelper struct {
	addressBookPath  string
	subscriptionURLs []string

	samHost string
	samPort string

	ext     bool
	verbose bool

	addressBook       []string
	remoteAddressBook []*addresslist
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
	var kv string
	log.Println("Seeking Address", pk)
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
				j.printKvs(r)
				return r
			}
		}
	}
	for _, r := range j.remoteAddressBook {
		for _, a := range r.RemoteAddressBook {
			r := strings.SplitN(a, ",", 2)
			if len(r) == 2 {
				if r[0] == j.trim(k.Host) {
					j.printKvs(r)
					i := i2pconv.I2pconv{}
					s, e := i.I2p64to32(r[1])
					if e != nil {
						return nil
					}
					v := []string{r[0], s, r[1]}
					return v
				}
			}
		}
	}
	return nil
}

// CheckAddressBook returns true if an address is present, false if not
func (j *JumpHelper) CheckAddressBook(pk string) bool {
	var kv string
	log.Println("Seeking Address", pk)
	if !strings.HasPrefix(pk, "http://") {
		kv = "http://" + pk
	} else {
		kv = pk
	}
	k, e := url.Parse(kv)
	log.Println("test", k.Host)
	if e != nil {
		return false
	}
	for _, a := range j.addressBook {
		r := strings.SplitN(a, ",", 2)
		if len(r) == 2 {
			if r[0] == j.trim(k.Host) {
				j.printKvs(r)
				return true
			}
		}
	}
	return false
}

// NewJumpHelper creates a new JumpHelper object
func NewJumpHelper(addressBookPath, host, port string, subs []string, use bool) (*JumpHelper, error) {
	return NewJumpHelperFromOptions(
		SetJumpHelperAddressBookPath(addressBookPath),
		SetJumpHelperHost(host),
		SetJumpHelperPort(port),
		SetJumpHelperUseHelper(use),
		SetJumpHelperSubscription(subs),
	)
}

// NewJumpHelperFromOptions creates a new JumpHelper object
func NewJumpHelperFromOptions(opts ...func(*JumpHelper) error) (*JumpHelper, error) {
	var j JumpHelper
	j.addressBookPath = "/var/lib/i2pd/addressbook/addresses.csv"
	j.samHost = "127.0.0.1"
	j.samPort = "7656"
	j.ext = false
	j.subscriptionURLs = []string{"http://joajgazyztfssty4w2on5oaqksz6tqoxbduy553y34mf4byv6gpq.b32.i2p/export/alive-hosts.txt"}
	for _, o := range opts {
		if err := o(&j); err != nil {
			return nil, fmt.Errorf("Service configuration error: %s", err)
		}
	}
	err := j.LoadAddressBook()
	if err != nil {
		return nil, err
	}
	if len(j.subscriptionURLs) < 1 {
		j.ext = false
	}
	log.Println("Creating a jumphelper")
	if j.ext {
		for _, u := range j.subscriptionURLs {
			nab, e := newAddressList(u, j.samHost, j.samPort)
			if e != nil {
				return nil, e
			}
			j.remoteAddressBook = append(j.remoteAddressBook, nab)
		}
	}
	return &j, err
}

func (j *JumpHelper) printKvs(kv []string) {
	if j.verbose {
		for i, s := range kv {
			log.Println("Key-value Pair", i, s)
		}
	}
}

// Subs Lists all the known address pairs
func (j *JumpHelper) Subs() []string {
	var r []string
	r = append(r, j.addressBook...)
	for _, x := range j.remoteAddressBook {
		r = append(r, x.RemoteAddressBook...)
	}
	return r
}
