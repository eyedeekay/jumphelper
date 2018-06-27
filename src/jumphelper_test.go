package jumphelper

import (
	"log"
	"testing"
)

func TestJumpHelperLocal(t *testing.T) {
	jh, err := NewJumpHelper("../addresses.csv", "127.0.0.1", "7054")
	if err != nil {
		log.Fatal(err)
	}
	x := jh.SearchAddressBook("i2p-projekt.i2p")
	log.Println("Testing Jumphelper Locally i2p-projekt.i2p", x)
	printKvs(x)
}

func TestJumpHelperLocalBool(t *testing.T) {
	jh, err := NewJumpHelper("../addresses.csv", "127.0.0.1", "7054")
	if err != nil {
		log.Fatal(err)
	}
	x := jh.CheckAddressBook("i2p-projekt.i2p")
	log.Println("Testing Jumphelper Locally i2p-projekt.i2p", x)
}

func TestJumpHelperLocalA(t *testing.T) {
	jh, err := NewJumpHelper("../addresses.csv", "127.0.0.1", "7054")
	if err != nil {
		log.Fatal(err)
	}
	x := jh.SearchAddressBook("http://i2p-projekt.i2p")
	log.Println("Testing Jumphelper Locally http://i2p-projekt.i2p", x)
}

func TestJumpHelperLocalBoolA(t *testing.T) {
	jh, err := NewJumpHelper("../addresses.csv", "127.0.0.1", "7054")
	if err != nil {
		log.Fatal(err)
	}
	x := jh.CheckAddressBook("http://i2p-projekt.i2p")
	log.Println("Testing Jumphelper Locally http://i2p-projekt.i2p", x)
}
