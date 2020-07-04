package main

import (
	"flag"
	"fmt"
	"github.com/alexmeli100/remit/lib/mtn"
	"log"
	"os"
)

var fs = flag.NewFlagSet("mtnmomo", flag.ExitOnError)
var host = fs.String("host", "localhost:3000", "callback host url")
var pk = fs.String("primary-key", "", "user api key")

func main() {
	fs.Parse(os.Args[1:])

	if *pk == "" {
		fs.Usage()
		os.Exit(1)
	}

	user, err := mtn.NewUser(*host, *pk)

	if err != nil {
		log.Fatal(err)
	}

	apiKey, err := user.Login()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Momo SandBox Credentials")
	fmt.Printf("user secret: %s\n", apiKey)
	fmt.Printf("user ID: %s\n", user.UserId)
}
