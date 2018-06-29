package jumphelper

import (
	"log"
	"testing"
	"time"
)

func ServiceStart() *Client {
	NewService("127.0.0.1", "7054", "../addresses.csv", "127.0.0.1", "7656", nil, false)
	log.Printf("testing service start")
	c, err := NewClient("127.0.0.1", "7054", true)
	if err != nil {
		log.Fatal(err, "Error connecting to service")
	}
	return c
}

func ServiceCheck(c *Client) {
	log.Printf("testing *Client Lookup")
	time.Sleep(2 * time.Second)
	if b, e := c.Check("i2p-projekt.i2p"); b {
		log.Println("Found i2p-projekt.i2p in addressbook", b)
	} else {
		log.Fatal("obscure error ", e)
	}
	time.Sleep(2 * time.Second)
	if b, _ := c.Check("fireaxe.i2p"); b {
		log.Fatal("Found fireaxe.i2p in addressbook")
	} else {
		log.Println("Subaddress fireaxe.i2p not found, this is correct", b)
	}
}

func ServiceHarderCheck(c *Client) {
	log.Printf("testing *Client Lookup")
	time.Sleep(2 * time.Second)
	if b, e := c.Check("i2p-projekt.i2p/en"); b {
		log.Println("Found i2p-projekt.i2p in addressbook")
	} else {
		log.Fatal(e)
	}
}

func ServiceRequest(c *Client) {
	log.Printf("testing *Client Request")
	time.Sleep(2 * time.Second)
	if b, e := c.Request("i2p-projekt.i2p"); e == nil {
		log.Println("Found", b, "in addressbook")
	} else {
		log.Fatal(e)
	}
	time.Sleep(2 * time.Second)
	if b, _ := c.Check("fireaxe.i2p"); b {
		log.Println("Found fireaxe.i2p in addressbook")
	} else {
		log.Fatal("Subaddress fireaxe.i2p not found, this is incorrect")
	}
}

func ServiceHarderRequest(c *Client) {
	log.Printf("testing *Client Request")
	time.Sleep(2 * time.Second)
	if b, e := c.Request("i2p-projekt.i2p/en"); e == nil {
		log.Println("Found", b, "in addressbook")
	} else {
		log.Fatal(e)
	}
}

func TestService(t *testing.T) {
	c := ServiceStart()
	ServiceCheck(c)
	time.Sleep(2 * time.Second)
	ServiceHarderCheck(c)
	time.Sleep(2 * time.Second)
	ServiceRequest(c)
	time.Sleep(2 * time.Second)
	ServiceHarderRequest(c)
}
