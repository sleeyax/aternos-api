// This example demonstrates how to start and stop a server over websockets.
// It's a little more advanced but the most performant and controllable way to use this package.
package main

import (
	"context"
	"encoding/json"
	aternos "github.com/sleeyax/aternos-api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	// Connect to the Aternos websocket server.
	wss, err := api.ConnectWebSocket()
	if err != nil {
		log.Fatal(err)
	}

	defer wss.Close() // closes the connection when the main function ends

	ctx, cancel := context.WithCancel(context.Background())
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	log.Println("Started websocket connection.")

	for {
		select {
		case msg := <-wss.Message:
			switch msg.Type {
			case "ready":
				// Start the server over HTTP.
				if err = api.StartServer(); err != nil {
					log.Println(err)
					return
				}

				// Run a goroutine in the background that sends a bunch of keep-alive requests at default intervals.
				go wss.SendHearthBeats(ctx)
			case "status":
				// Current server status, containing a bunch of other useful info such as IP address/Dyn IP to connect to, amount of active players, detected problems etc.
				var info aternos.ServerInfo
				json.Unmarshal(msg.MessageBytes, &info)

				log.Printf("Server status: %s\n", info.StatusLabel)

				if info.Status == aternos.Online {
					log.Println("Name:", info.Name)
					log.Println("Dyn IP:", info.DynIP)
					log.Println("Address:", info.Address)
					log.Println("Port:", info.Port)
				}

			}
		// Stop the server, close the connection & quit the app when CTRL + C is pressed.
		case <-interruptSignal:
			if err = api.StopServer(); err != nil {
				log.Println(err)
			}
			cancel() // stop sending heartbeats
			return
		}
	}
}
