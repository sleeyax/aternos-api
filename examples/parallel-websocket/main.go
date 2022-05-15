// This example demonstrates how to start and stop multiple servers over multiple websocket connections.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	aternos "github.com/sleeyax/aternos-api"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	var wg sync.WaitGroup

	apis, err := readFile("/home/quinten/Programming/go/src/github.com/sleeyax/aternos-api/examples/parallel-websocket/sessions.txt")
	if err != nil {
		log.Fatalf("main: failed to read file (%s)\n", err)
	}

	for i, api := range apis {
		wg.Add(1)
		go start(i, api, &wg)
	}

	wg.Wait()
	log.Println("main: all goroutines finished")
}

func readFile(textFile string) ([]*aternos.Api, error) {
	file, err := os.Open(textFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var apis []*aternos.Api

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		split := strings.Split(line, ",")
		if len(split) < 2 {
			return nil, errors.New("invalid sessions.txt")
		}

		o := &aternos.Options{
			Cookies: []*http.Cookie{
				{
					Name:  "ATERNOS_LANGUAGE",
					Value: "en",
				},
				{
					Name:  "ATERNOS_SESSION",
					Value: split[0],
				},
				{
					Name:  "ATERNOS_SERVER",
					Value: split[1],
				},
			},
			InsecureSkipVerify: true,
		}

		if len(split) >= 3 {
			p, _ := url.Parse(split[2])
			o.Proxy = p
		}

		api := aternos.New(o)

		apis = append(apis, api)
	}

	return apis, nil
}

func start(id int, api *aternos.Api, wg *sync.WaitGroup) {
	// Connect to the Aternos websocket server.
	wss, err := api.ConnectWebSocket()
	if err != nil {
		log.Printf("main %d: failed to connect to websocket server (%s)\n", id, err)
		wg.Done()
		return
	}

	log.Printf("main %d: connected to websocket\n", id)

	ctx, cancel := context.WithCancel(context.Background())

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	defer func() {
		cancel()
		if err = wss.Close(); err != nil {
			log.Printf("main %d: failed to close websocket connection (%s)\n", id, err)
		}
		wg.Done()
	}()

	for {
		select {
		case msg, ok := <-wss.Message:
			if !ok {
				log.Printf("main %d: stopped reading messages from websocket \n", id)
				return
			}

			switch msg.Type {
			case "ready":
				// Start the server over HTTP.
				if err = api.StartServer(); err != nil {
					log.Printf("main %d: failed to start server (%e) \n", id, err)
					return
				}

				// Run a goroutine in the background that sends a bunch of keep-alive requests.
				go wss.StartHearthBeat(ctx)
			case "status":
				// Current server status, containing a bunch of other useful info such as IP address/Dyn IP to connect to, amount of active players, detected problems etc.
				var info aternos.ServerInfo
				json.Unmarshal(msg.MessageBytes, &info)

				log.Printf("main %d: server is %s\n", id, info.StatusLabel)

				if info.StatusLabel == "waiting" {
					api.ConfirmServer(nil, time.Second*10)
				}

				if info.Status == aternos.Offline {
					return
				}
			}
		// Stop the server, close the connection & quit the app when CTRL + C is pressed.
		case <-interruptSignal:
			if err = api.StopServer(); err != nil {
				log.Printf("main %d: failed to stop server (%s)\n", id, err)
				return
			}
		}
	}
}
