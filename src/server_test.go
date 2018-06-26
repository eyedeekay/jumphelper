package jumphelper

import (
	"log"
	"testing"
)

func ServiceStart() *Client {
    log.Printf("testing service start")
    c, err := NewClient("127.0.0.1", "7054")
	if err != nil {
		log.Fatal(err, "Error connecting to service")
	}
    return c
}

func ServiceCheck(c *Client) {
    log.Printf("testing *Client Lookup")
    if b, e := c.Check("i2p-projekt.i2p"); b {
        log.Println("Found i2p-projekt.i2p in addressbook")
    }else{
        log.Fatal(e)
    }
}

func ServiceHarderCheck(c *Client) {
    log.Printf("testing *Client Lookup")
    if b, e := c.Check("i2p-projekt.i2p/en"); b {
        log.Println("Found i2p-projekt.i2p in addressbook")
    }else{
        log.Fatal(e)
    }
}

func TestService(t *testing.T) {
    Service()
    c := ServiceStart()
    ServiceCheck(c)
    ServiceHarderCheck(c)
}
