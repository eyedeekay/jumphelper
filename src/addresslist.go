package jumphelper

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/eyedeekay/gosam"
    "github.com/eyedeekay/i2pasta/convert"
)

type addresslist struct {
	addressBookURL string
	samHost        string
	samPort        string
    Lock           bool

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
    a.Lock = false
	return nil
}

func (a *addresslist) trim(k string) string {
	r := strings.Replace(k, "http://", "", -1)
	if strings.HasSuffix(k, ".i2p") {
		return r
	}
	r = strings.SplitN(r, ".i2p", -1)[0]
	return strings.TrimSpace(strings.TrimSuffix(r, "/"))
}

func (a *addresslist) SearchAddressList(host string) []string {
    if a.Lock == false {
        for _, addresspair := range a.RemoteAddressBook {
            r := strings.SplitN(addresspair, ",", 2)
            if len(r) == 2 {
                if r[0] == a.trim(host) {
                    //j.printKvs(r)
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

func newAddressList(u, samhost, samport string) (*addresslist, error) {
	var a addresslist
	var err error
	a.samHost = samhost
	a.samPort = samport
	a.addressBookURL = u
    a.Lock = true
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
