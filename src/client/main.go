package main

import (
    "flag"
    "log"

    "github.com/eyedeekay/jumphelper/src"
)

func main(){
    log.Println("Starting server:")
    host := flag.String("host", "127.0.0.1", "Host address to listen on.")
	port := flag.String("port", "7054", "Port to listen on.")
    url := flag.String("url", "false", "Port to listen on.")
    flag.Parse()

    s, err := jumphelper.NewClient(*host, *port)
    if err != nil {
        log.Fatal(err, "Error starting server" )
    }

    if *url != "false" {
        s.Check("i2p-projekt.i2p")
    }
}
