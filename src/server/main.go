package main

import (
	"flag"
	"log"

	//"github.com/eyedeekay/jumphelper/src"
	".."
)

func main() {
	log.Println("Starting server:")
	host := flag.String("host", "0.0.0.0", "Host address to listen on.")
	port := flag.String("port", "7054", "Port to listen on.")
	samhost := flag.String("samhost", "127.0.0.1", "Host address to listen on.")
	samport := flag.String("samport", "7656", "Port to listen on.")
	book := flag.String("hostfile", "./addresses.csv", "Local address book")
    //useremote := flag.Bool("useremote", false, "Use external address books")

	flag.Parse()

	s, err := jumphelper.NewServer(*host, *port, *book, *samhost, *samport)
	if err != nil {
		log.Fatal(err, "Error starting server")
	}
	s.Serve()
}
