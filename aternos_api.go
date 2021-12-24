package aternos_api

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sleeyax/gotcha"
	"github.com/sleeyax/gotcha/adapters/fhttp"
	httpx "github.com/useflyent/fhttp"
	"net/http"
	"strings"
	"time"
)

var ServerAlreadyStartedError = errors.New("Server already running!")
var ServerAlreadyStoppedError = errors.New("Server already stopped!")

type AternosApi struct {
	Config
	client *gotcha.Client
	sec    string
}

// Make allocates a new Aternos API instance.
func Make(config Config) AternosApi {
	client, _ := gotcha.NewClient(&gotcha.Options{
		Adapter:   fhttp.NewAdapter(),
		PrefixURL: "https://aternos.org/",
		Headers: http.Header{
			"user-agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"},
			"accept":             {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
			"accept-language":    {"en-US,en;q=0.9"},
			"accept-encoding":    {"gzip, deflate, br"},
			"cookie":             config.Cookies,
			httpx.HeaderOrderKey: []string{"user-agent", "accept", "accept-language", "accept-encoding", "cookie"},
		},
		FollowRedirect: false,
		Retry:          false,
	})

	var sec string
	for _, v := range config.Cookies {
		prefix := "ATERNOS_SEC_"
		if strings.HasPrefix(v, prefix) {
			sec = strings.ReplaceAll(strings.ReplaceAll(v, prefix, ""), "=", "%3A")
		}
	}

	return AternosApi{
		Config: config,
		client: client,
		sec:    sec,
	}
}

// GetPlayers returns all online players.
func (api AternosApi) GetPlayers() ([]string, error) {
	res, err := api.client.Get("players")
	if err != nil {
		return nil, err
	}

	defer res.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	players := make([]string, 0)

	document.Find("div.playername").Each(func(i int, selection *goquery.Selection) {
		players = append(players, strings.TrimSpace(selection.Text()))
	})

	return players, nil
}

// GetServerInfo returns server information.
func (api AternosApi) GetServerInfo() (ServerInfo, error) {
	res, err := api.client.Get("server")
	if err != nil {
		return ServerInfo{}, err
	}

	defer res.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return ServerInfo{}, err
	}

	info := ServerInfo{
		Status:   strings.TrimSpace(document.Find("span.statuslabel-label").First().Text()),
		Software: strings.TrimSpace(document.Find("span#software").First().Text()),
		Address:  strings.TrimSpace(document.Find("span#ip").First().Text()),
		// TODO: port
		Version: strings.TrimSpace(document.Find("span#version").First().Text()),
	}

	return info, nil
}

// TODO: start & stop server over websockets (wss://aternos.org/hermes/)

// StartServer starts your Minecraft server.
func (api AternosApi) StartServer() error {
	info, err := api.GetServerInfo()
	if err != nil {
		return err
	}

	if info.Status == "Online" {
		return ServerAlreadyStartedError
	}

	res, err := api.client.Get(fmt.Sprintf("panel/ajax/start.php?headstart=0&access-credits=0&SEC=%s&TOKEN=%s", api.sec, api.Token))
	if err != nil {
		return err
	}

	defer res.Close()

	json, err := res.Json()
	if err != nil {
		return err
	}
	if json["success"] == false {
		return errors.New("Aternos failed to start the server.")
	}

	return nil
}

// ConfirmServer sends a confirmation that 'you're still active' while waiting for the server to start.
// You should call this function right after starting the server.
// This might not always be required anymore though.
//
// delay specifies the amount of seconds to wait before submitting the confirmation.
func (api AternosApi) ConfirmServer(delay time.Duration) error {
	for {
		time.Sleep(delay)

		info, err := api.GetServerInfo()
		if err != nil {
			return err
		}

		status := info.Status

		if status != "Preparing" && status != "Online" {
			res, err := api.client.Get(fmt.Sprintf("panel/ajax/confirm.php?headstart=0&access-credits=0&SEC=%s&TOKEN=%s", api.sec, api.Token))
			if err != nil {
				return err
			}

			res.Close()
		} else {
			break
		}
	}

	return nil
}

// StopServer stops the minecraft server.
// This function doesn't wait until the server is fully stopped, it only requests a shutdown.
func (api AternosApi) StopServer() error {
	info, err := api.GetServerInfo()
	if err != nil {
		return err
	}

	if info.Status == "Offline" {
		return ServerAlreadyStoppedError
	}

	_, err = api.client.Get(fmt.Sprintf("panel/ajax/stop.php?SEC=%s&TOKEN=%s", api.sec, api.Token))

	return err
}
