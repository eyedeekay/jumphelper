package main

import (
	"flag"
	"log"

	//"github.com/eyedeekay/jumphelper/src"
	".."
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	log.Println("Starting server:")
	var subscriptions arrayFlags
	host := flag.String("host", "0.0.0.0", "Host address to listen on.")
	port := flag.String("port", "7054", "Port to listen on.")
	samhost := flag.String("samhost", "127.0.0.1", "Host address to listen on.")
	samport := flag.String("samport", "7656", "Port to listen on.")
	book := flag.String("hostfile", "./addresses.csv", "Local address book")
	useremote := flag.Bool("useremote", false, "Use external address books")
	flag.Var(&subscriptions, "subs", "Subscription URL(Can be specified multiple times)")

	flag.Parse()

    if len(subscriptions) < 1 {
        subscriptions = append(subscriptions, "http://joajgazyztfssty4w2on5oaqksz6tqoxbduy553y34mf4byv6gpq.b32.i2p/export/alive-hosts.txt")
    }

	s, err := jumphelper.NewServer(*host, *port, *book, *samhost, *samport, subscriptions, *useremote)
	if err != nil {
		log.Fatal(err, "Error starting server")
	}
	s.Serve()
}
