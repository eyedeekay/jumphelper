package jumphelper

import (
	"log"
	"testing"
)

func TestJumpHelperLocal(t *testing.T) {
	jh, err := NewJumpHelper("../addresses.csv", "127.0.0.1", "7054", nil, false)
	if err != nil {
		log.Fatal(err)
	}
	x := jh.SearchAddressBook("i2p-projekt.i2p")
	y := jh.SearchAddressBook("fireaxe.i2p.i2p")
	log.Println("Testing Jumphelper Locally i2p-projekt.i2p", x, "fireaxe.i2p", y)
	printKvs(x)
}

func TestJumpHelperLocalBool(t *testing.T) {
	jh, err := NewJumpHelper("../addresses.csv", "127.0.0.1", "7054", nil, false)
	if err != nil {
		log.Fatal(err)
	}
	x := jh.CheckAddressBook("i2p-projekt.i2p")
	y := jh.CheckAddressBook("fireaxe.i2p")
	log.Println("Testing Jumphelper check-only Locally i2p-projekt.i2p", x, "fireaxe.i2p", y)
}

func TestJumpHelperLocalA(t *testing.T) {
	jh, err := NewJumpHelper("../addresses.csv", "127.0.0.1", "7054", nil, false)
	if err != nil {
		log.Fatal(err)
	}
	x := jh.SearchAddressBook("http://i2p-projekt.i2p")
	y := jh.SearchAddressBook("http://fireaxe.i2p")
	log.Println("Testing Jumphelper Locally http://i2p-projekt.i2p", x, "http://fireaxe.i2p", y)
}

func TestJumpHelperLocalBoolA(t *testing.T) {
	jh, err := NewJumpHelper("../addresses.csv", "127.0.0.1", "7054", nil, false)
	if err != nil {
		log.Fatal(err)
	}
	x := jh.CheckAddressBook("http://i2p-projekt.i2p")
	y := jh.CheckAddressBook("http://fireaxe.i2p")
	log.Println("Testing Jumphelper check-only Locally http://i2p-projekt.i2p", x, "http://fireaxe.i2p", y)
}
