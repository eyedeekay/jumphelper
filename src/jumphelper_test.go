package jumphelper

import (
	"log"
	"testing"
)

func TestJumpHelperLocal(t *testing.T) {
	jh, err := NewJumpHelper("../addresses.csv")
	if err != nil {
		log.Fatal(err)
	}
	x := jh.SearchAddressBook("i2p-projekt.i2p")
	log.Println(x)
}

func TestJumpHelperLocalBool(t *testing.T) {
	jh, err := NewJumpHelper("../addresses.csv")
	if err != nil {
		log.Fatal(err)
	}
	x := jh.CheckAddressBook("i2p-projekt.i2p")
	log.Println(x)
}

func TestJumpHelperLocalA(t *testing.T) {
	jh, err := NewJumpHelper("../addresses.csv")
	if err != nil {
		log.Fatal(err)
	}
	x := jh.SearchAddressBook("http://i2p-projekt.i2p")
	log.Println(x)
}

func TestJumpHelperLocalBoolA(t *testing.T) {
	jh, err := NewJumpHelper("../addresses.csv")
	if err != nil {
		log.Fatal(err)
	}
	x := jh.CheckAddressBook("http://i2p-projekt.i2p")
	log.Println(x)
}
