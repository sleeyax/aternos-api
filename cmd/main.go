package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	aternos "github.com/sleeyax/aternos-api"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	// Parse CLI flags.
	session := flag.String("session", "", "ATERNOS_SESSION")
	lang := flag.String("lang", "en", "ATERNOS_LANGUAGE")
	server := flag.String("server", "", "ATERNOS_SERVER")
	proxy := flag.String("proxy", "", "optional proxy to connect to")
	flag.Parse()

	// Check if all required flags are specified.
	if *session == "" || *server == "" {
		fmt.Println("Missing cookie values.")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println()
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Example: aternos-api -lang en -proxy http://127.0.0.1:8888 -server a6liGFpEBjF1LvIS -session a6HBs6rMiHJvDHxQQfCDKQQmutFSWgIJhQVgzaPOSvr6TCXtdZMPKIlGjqubur80l2AYp9r9seRMilWHznPv05U9mVPCiPm8we2e")
		os.Exit(0)
	}

	// Parse proxy (if specified).
	var p *url.URL
	if *proxy != "" {
		p, _ = url.Parse(*proxy)
	}

	// Create a new api instance, with all required authentication cookies set.
	api := aternos.New(&aternos.Options{
		Cookies: []*http.Cookie{
			{
				Name:  "ATERNOS_LANGUAGE",
				Value: *lang,
			},
			{
				Name:  "ATERNOS_SESSION",
				Value: *session,
			},
			{
				Name:  "ATERNOS_SERVER",
				Value: *server,
			},
		},
		Proxy:              p,
		InsecureSkipVerify: true,
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

				// Run a goroutine in the background that sends a bunch of keep-alive requests at intervals.
				go wss.SendHearthBeats(ctx)
			case "status":
				// Current server status, containing a bunch of other useful info such as IP address/Dyn IP to connect to, amount of active players, detected problems etc.
				var serverInfo aternos.ServerInfo
				json.Unmarshal(msg.MessageBytes, &serverInfo)

				log.Printf("Server status: %s\n", serverInfo.StatusLabel)

				switch serverInfo.Status {
				case aternos.Starting:
					// Start streaming console logs.
					wss.StartConsoleLogStream()
				case aternos.Online:
					// Stop streaming console logs
					// and start streaming heap & tick info.
					wss.StopConsoleLogStream()
					wss.StartHeapInfoStream()
					wss.StartTickStream()
				case aternos.Stopping:
					// Stop heap & tick stream..
					wss.StopHeapInfoStream()
					wss.StopTickStream()
				case aternos.Offline:
					cancel() // stops sending heartbeats.
					return   // closes the connection.
				}
			case "queue_reduced":
				// Waiting queue.
				// Your server will start once it's your turn in this queue.
				var queue aternos.QueueReduction
				json.Unmarshal(msg.MessageBytes, &queue)
				log.Printf("Reduced queue %d: %d people in queue (max time %d)\n", queue.Number, queue.Total, queue.MaxTime)
			case "line":
				// Console stream.
				// Usually contains verbose output that is useful for logging/debugging purposes.
				log.Printf("%s: %s\n", strings.Title(msg.Stream), thisOrThat(msg.Data.Content, msg.Type))
			case "heap":
				// Current memory usage.
				var heap aternos.Heap
				json.Unmarshal(msg.Data.ContentBytes, &heap)
				log.Printf("Heap usage in bytes: %d\n", heap.Usage)
			case "tick":
				var tick aternos.Tick
				json.Unmarshal(msg.Data.ContentBytes, &tick)
				log.Printf("Tick time: %f\n", tick.AverageTickTime)
			case "started":
				// Started sending a specific stream (console, heap or tick).
				log.Printf("Started stream: %s\n", msg.Stream)
			case "stopped":
				// Stopped sending a specific stream (console, heap or tick).
				log.Printf("Stopped stream: %s\n", msg.Stream)
			case "connected":
				// Stream has connected.
				// Usually this message is sent after starting a stream.
				// It may contain no additional data.
				log.Printf("Connected: %s\n", msg.Data)
			case "disconnected":
				// Stream or server has disconnected.
				log.Printf("Disconnected: %s\n", msg.Data)
			case "backup_progress":
				// Backup status.
				// Backups are automatically made while playing or when the server is stopped.
				var backupProgress aternos.BackupProgress
				json.Unmarshal(msg.MessageBytes, &backupProgress)
				log.Printf("Backup: %d/100%% complete (%s)\n", backupProgress.Progress, backupProgress.Action)
			default:
				data, _ := json.Marshal(msg)
				fmt.Printf("Unhandled message '%s', skipping\nDump:%s\n", msg.Type, data)
			}
		// Stop server & quit when CTRL + C is pressed.
		case <-interruptSignal:
			if err = api.StopServer(); err != nil {
				return
			}
		}
	}
}

// thisOrThat returns the first item that isn't an empty string.
func thisOrThat(left, right string) string {
	if left != "" {
		return left
	}
	return right
}
