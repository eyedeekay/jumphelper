package jumphelper

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/eyedeekay/gosam"
)

type addresslist struct {
	addressBookURL string
	samHost        string
	samPort        string

	samBridgeConn *goSam.Client

	tr     *http.Transport
	client *http.Client

	RemoteAddressBook []string
}

// SyncRemoteAddressBooks syncs addressbooks from subscription services to the standalone addressbook
func (a *addresslist) SyncRemoteAddressBooks() error {
	log.Println("Syncing Subscription Contents")
	resp, err := a.client.Get(a.addressBookURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return e
	}
	lines := strings.Split(string(b), "\n")
	for _, l := range lines {
		kv := strings.SplitN(l, "=", 2)
		if len(kv) == 2 {
			a.RemoteAddressBook = append(a.RemoteAddressBook, kv[0]+","+kv[1])
		}
	}
	log.Println("Subscription Contents Synced from", a.addressBookURL)
	return nil
}

func newAddressList(u, samhost, samport string) (*addresslist, error) {
	var a addresslist
	var err error
	a.samHost = samhost
	a.samPort = samport
	a.addressBookURL = u
	a.samBridgeConn, err = goSam.NewClientFromOptions(
		goSam.SetHost(a.samHost),
		goSam.SetPort(a.samPort),
		goSam.SetInLength(2),
		goSam.SetOutLength(2),
		goSam.SetInQuantity(15),
		goSam.SetInBackups(5),
		goSam.SetOutQuantity(5),
		goSam.SetOutBackups(5),
		goSam.SetUnpublished(true),
	)
	if err != nil {
		return nil, err
	}
	a.tr = &http.Transport{
		Dial: a.samBridgeConn.Dial,
	}
	a.client = &http.Client{Transport: a.tr}
	err = a.SyncRemoteAddressBooks()
	if err != nil {
		return nil, err
	}
	return &a, nil
}
