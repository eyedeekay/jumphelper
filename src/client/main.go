package main

import (
	"flag"
	"log"

	"github.com/eyedeekay/jumphelper/src"
)

func main() {
	host := flag.String("host", "127.0.0.1", "Host address to listen on.")
	port := flag.String("port", "7054", "Port to listen on.")
	url := flag.String("url", "false", "URL to check.")
	//addr := flag.Bool("url", "false", "URL to check.")
	flag.Parse()

	s, err := jumphelper.NewClient(*host, *port)
	if err != nil {
		log.Fatal(err, "Error starting client")
	}

	if *url != "false" {
		if !*addr {
			if b, e := s.Check("i2p-projekt.i2p"); b {
				log.Println("true")
				if e != nil {
					log.Fatal(e)
				}
			} else {
				log.Println("false")
				if e != nil {
					log.Fatal(e)
				}
			}
		} else {
			if s, e := s.Request("i2p-projekt.i2p"); b {
				log.Println(s)
				if e != nil {
					log.Fatal(e)
				}
			} else {
				log.Println("false")
				if e != nil {
					log.Fatal(e)
				}
			}
		}
	}
}
