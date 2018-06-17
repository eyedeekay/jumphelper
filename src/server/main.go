package main

import (
    //"github.com/eyedeekay/jumphelper"
    "log"

    ".."
)

func main(){
    log.Println("Starting server:")
    s, err := jumphelper.NewServer("0.0.0.0", "7054", "/var/lib/i2pd/addressbook/addresses.csv")
    if err != nil {
        log.Fatal(err, "Error starting server" )
    }
    s.Serve()
}
