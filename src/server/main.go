package main

import (
	"flag"
	"log"
    "strings"

	"github.com/eyedeekay/jumphelper/src"
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
	samhost := flag.String("samhost", "127.0.0.1", "Host address to listen on.")
	samport := flag.String("samport", "7656", "Port to listen on.")
	book := flag.String("hostfile", "./addresses.csv", "Local address book.")
	useremote := flag.Bool("useremote", true, "Use external address books.")
	share := flag.Bool("share", false, "Repeat concatenated listings as subscription list")
	verbose := flag.Bool("verbose", false, "Verbose?")

	flag.Var(&subscriptions, "subs", "Subscription URL(Can be specified multiple times)")

	flag.Parse()

	if len(subscriptions) < 1 {
		subscriptions = append(subscriptions, "http://joajgazyztfssty4w2on5oaqksz6tqoxbduy553y34mf4byv6gpq.b32.i2p/export/alive-hosts.txt")
	}

	s, err := jumphelper.NewServer(*host, *port, *book, *samhost, *samport,
        subscriptions,
        *useremote, *verbose, *share)
	if err != nil {
		log.Fatal(err, "Error starting server")
	}
	s.Serve()
}
