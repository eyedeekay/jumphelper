package main

import (
	"flag"
	"log"
	"strings"
)

import (
	"github.com/eyedeekay/jumphelper/src"
    "github.com/eyedeekay/sam-forwarder"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	var r string
	for _, s := range *i {
		r += s + ","
	}
	return strings.TrimSuffix(r, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	log.Println("Starting server:")
	var subscriptions arrayFlags
	host := flag.String("host", "127.0.0.1", "Host address to listen on.")
	port := flag.String("port", "7054", "Port to listen on.")
	samhost := flag.String("samhost", "127.0.0.1", "Host connect to SAM on.")
	samport := flag.String("samport", "7656", "SAM port.")
	book := flag.String("hostfile", "./addresses.csv", "Local address book.")
	useremote := flag.Bool("useremote", true, "Use external address books.")
	share := flag.Bool("share", false, "Repeat concatenated listings as subscription list")
    forward := flag.Bool("i2p", false, "Forward service to an i2p destination over a SAM connection")
    tunname := flag.String("tunname", "jumphelper", "Tunnel name for SAM forwarding.")
	verbose := flag.Bool("verbose", false, "Verbose?")

	flag.Var(&subscriptions, "subs", "Subscription URL(Can be specified multiple times)")

	flag.Parse()

	if len(subscriptions) < 1 {
		subscriptions = append(subscriptions, "http://joajgazyztfssty4w2on5oaqksz6tqoxbduy553y34mf4byv6gpq.b32.i2p/export/alive-hosts.txt")
	}
    var forwarder *samforwarder.SAMForwarder
    var err error
    if *forward {
        if forwarder, err = samforwarder.NewSAMForwarderFromOptions(
            samforwarder.SetHost(*host),
            samforwarder.SetPort(*port),
            samforwarder.SetSAMHost(*samhost),
            samforwarder.SetSAMPort(*samport),
            samforwarder.SetName(*tunname),
        ); err == nil {
            go forwarder.Serve()
        }else{
            log.Fatal(err, "Error starting forwarder")
        }
    }

	s, err := jumphelper.NewServer(*host, *port, *book, *samhost, *samport,
		subscriptions,
		*useremote, *verbose, *share)
	if err != nil {
		log.Fatal(err, "Error starting server")
	}
	s.Serve()
}
