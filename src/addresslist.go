package jumphelper

import (
	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"net/url"
	"strings"

	"github.com/eyedeekay/gosam"
	"github.com/eyedeekay/i2pasta/convert"
)

type addresslist struct {
    addressBookURL  string
    samHost string
	samPort string

    samBridgeConn *goSam.Client

    tr     *http.Transport
	client *http.Client
    //Stores comma-separated name,b32 values
    //addressBook       []string
    //Stores comma-separated name,b64 values
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
			i := i2pconv.I2pconv{}
			s, e := i.I2p64to32(kv[1])
			if e != nil {
				return e
			}
			a.RemoteAddressBook = append(a.RemoteAddressBook, kv[0]+","+s)
		}
	}
	log.Println("Subscription Contents Synced from", a.addressBookURL)
	return nil
}

func newAddressList(u, samhost, samport string) (*addresslist, error){
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
