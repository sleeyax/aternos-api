// This example demonstrates how to start and stop a server over HTTP.
// It's the simplest but not the most performant way to use this package.
package main

import (
	aternos "github.com/sleeyax/aternos-api"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// session & server id are set by environment variables.
	session := os.Getenv("session")
	server := os.Getenv("server")
	if session == "" || server == "" {
		log.Fatalln("not all environment variables are set!")
	}

	// Create new Aternos api instance, providing the required authentication cookies.
	api := aternos.New(&aternos.Options{
		Cookies: []*http.Cookie{
			{
				Name:  "ATERNOS_LANGUAGE",
				Value: "en",
			},
			{
				Name:  "ATERNOS_SESSION",
				Value: session,
			},
			{
				Name:  "ATERNOS_SERVER",
				Value: server,
			},
		},
	})

	// Start the server.
	if err := api.StartServer(); err != nil {
		log.Fatalln(err)
	}

	log.Println("starting server...")

	// Keep confirming the server every 10 seconds.
	// Please note that these confirmations might not be required anymore in the latest version of Aternos.
	log.Println("confirming...")
	if err := api.ConfirmServer(10 * time.Second); err != nil {
		log.Fatalln(err)
	}
	log.Println("confirmed!")

	// Get the current server status & info.
	info, err := api.GetServerInfo()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("server is", info.StatusLabel)
	log.Println("name:", info.Name)
	log.Println("dyn IP:", info.DynIP)
	log.Println("address:", info.Address)
	log.Println("port:", info.Port)

	// Stop the server right after it came online.
	// Normally in a production app you wouldn't do this, of course.
	// This is only for demonstration purposes.
	if info.Status == aternos.Online {
		if err = api.StopServer(); err != nil {
			log.Fatalln(err)
		}
		log.Println("server is stopping...")
		// To check whether the server has actually fully stopped, you can periodically check the status with:
		// info, err := api.GetServerInfo()
		// ...
	}
}
