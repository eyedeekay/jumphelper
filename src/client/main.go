package main

import (
	"flag"
	"fmt"
	"log"

	".."
)

func main() {
	host := flag.String("host", "127.0.0.1", "Host address to listen on.")
	port := flag.String("port", "7854", "Port to listen on.")
	url := flag.String("url", "false", "URL to check.")
	addr := flag.Bool("addr", false, "Show base32 URL?.")
    addr64 := flag.Bool("addr64", false, "Show base32 URL?.")
	verbose := flag.Bool("verbose", false, "Verbose?.")
	flag.Parse()

	c, err := jumphelper.NewClient(*host, *port, *verbose)
	if err != nil {
		log.Fatal(err, "Error starting client")
	}

	if *url != "false" {
		if !*addr {
			if b, e := c.Check(*url); b {
				fmt.Println("true")
				if e != nil {
					log.Fatal(e)
				}
			} else {
				fmt.Println("false")
				if e != nil {
					log.Fatal(e)
				}
			}
		} else {
			if s, e := c.Request(*url); s != "FALSE" {
				fmt.Println("true", s)
				if e != nil {
					log.Fatal(e)
				}
			} else {
				fmt.Println("false")
				if e != nil {
					log.Fatal(e)
				}
			}
		}
	}
}
