package jumphelper

import (
	//"log"
	"testing"
)

func TestAddressbookLocal(t *testing.T) {
	newAddressList(
		"http://joajgazyztfssty4w2on5oaqksz6tqoxbduy553y34mf4byv6gpq.b32.i2p/export/alive-hosts.txt",
		"127.0.0.1",
		"7656",
	)
}
