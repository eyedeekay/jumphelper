package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/eyedeekay/jumphelper/src"
)

func main() {
	host := flag.String("host", "127.0.0.1", "Host address to listen on.")
	port := flag.String("port", "7854", "Port to listen on.")
	url := flag.String("url", "false", "URL to check.")
	addr := flag.Bool("addr", false, "Show base32 URL?.")
	addr64 := flag.Bool("addr64", false, "Show base64 address?.")
	verbose := flag.Bool("verbose", false, "Verbose?.")
	//domain := flag.String("domain", ".i2p", "i2p domain name to register")
	flag.Parse()

	c, err := jumphelper.NewClient(*host, *port, *verbose)
	if err != nil {
		log.Fatal(err, "Error starting client")
	}

	if *url != "false" {
		if !*addr {
			if b, e := c.Check(*url); b {
				if e != nil {
					log.Fatal(e)
				}
				fmt.Println("true")
			} else {
				if e != nil {
					log.Fatal(e)
				}
				fmt.Println("false")
			}
		} else {
			if s, e := c.Request(*url); s != "FALSE" {
				if e != nil {
					log.Fatal(e)
				}
				fmt.Println("true", s)
			} else {
				if e != nil {
					log.Fatal(e)
				}
				fmt.Println("false")
			}
		}
		if *addr64 {
			if s, e := c.Jump(*url); len(s) > 40 {
				if e != nil {
					log.Fatal(e)
				}
				fmt.Println(s)
			}
		}
	}
}
