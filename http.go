package aternos_api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// getDocument sends a GET request to the specified url and reads the response as a goquery.Document.
func (api *Api) getDocument(url string) (*goquery.Document, error) {
	res, err := api.client.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return document, nil
}

// genSec generates an AJAX security token called SEC.
func (api *Api) genSec() {
	key := randomString(11) + "00000"
	value := randomString(11) + "00000"

	api.sec = fmt.Sprintf("%s:%s", key, value)
	api.client.Options.CookieJar.SetCookies(api.client.Options.FullUrl, []*http.Cookie{
		{
			Name:  fmt.Sprintf("ATERNOS_SEC_%s", key),
			Value: value,
		},
	})
}

// extractToken extracts and unpacks the AJAX TOKEN from given HTML document.
func (api *Api) extractToken(document *goquery.Document) error {
	var script string

	document.Find("script[type='text/javascript']").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		script = strings.TrimSpace(selection.Text())
		return !strings.Contains(script, "window")
	})

	if script == "" {
		return errors.New("failed to find token")
	}

	vm := goja.New()

	if err := vm.Set("atob", atob); err != nil {
		return err
	}

	script = fmt.Sprintf("window = {}; %s window['AJAX_TOKEN'];", script)

	v, err := vm.RunString(script)
	if err != nil {
		return err
	}

	api.token = v.String()

	return nil
}

// GetServerInfo fetches all server information over HTTP.
func (api *Api) GetServerInfo() (ServerInfo, error) {
	document, err := api.getDocument("server")
	if err != nil {
		return ServerInfo{}, err
	}

	var script string
	var info ServerInfo
	prefix := "var lastStatus ="
	suffix := ";"

	document.Find("script:not([src])").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		script = strings.TrimSpace(selection.Text())
		return !strings.HasPrefix(script, prefix)
	})

	if script == "" {
		return ServerInfo{}, errors.New("failed to find server info")
	}

	data := strings.TrimSuffix(strings.ReplaceAll(script, prefix, ""), suffix)
	if err = json.Unmarshal([]byte(data), &info); err != nil {
		return ServerInfo{}, err
	}

	api.genSec()
	err = api.extractToken(document)

	return info, err
}

// StartServer starts your Minecraft server over HTTP.
func (api *Api) StartServer() error {
	info, err := api.GetServerInfo()
	if err != nil {
		return err
	}

	if info.Status == Online {
		return ServerAlreadyStartedError
	}

	res, err := api.client.Get(fmt.Sprintf("panel/ajax/start.php?headstart=0&access-credits=0&SEC=%s&TOKEN=%s", api.sec, api.token))
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

// ConfirmServer sends a confirmation over HTTP to claim that 'you're still active'.
// You should call this function once it's your turn in queue, after the server has started.
//
// This function will run synchronously if no context is given (nil value).
//
// Delay specifies the amount of seconds to wait before submitting the next confirmation.
// Recommended to wait time is around 10 seconds.
func (api *Api) ConfirmServer(ctx context.Context, delay time.Duration) error {
	isAsync := ctx != nil

	if !isAsync {
		ctx = context.Background()
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			info := ServerInfo{Status: Preparing}

			// When the function is called asynchronously (aka as a go routine like `go ConfirmServer()`),
			// we'll just assume the server is known to be in a preparing state already.
			// This reduces additional and unnecessary overhead.
			if !isAsync {
				var err error
				info, err = api.GetServerInfo()
				if err != nil {
					if isAsync {
						log.Println("Failed to get server info while confirming server:", err)
					}
					return err
				}
			}

			if info.Status != Preparing {
				time.Sleep(delay)
				break
			}

			res, err := api.client.Get(fmt.Sprintf("panel/ajax/confirm.php?headstart=0&access-credits=0&SEC=%s&TOKEN=%s", api.sec, api.token))
			if err != nil {
				if isAsync {
					log.Println("Failed to confirm server:", err)
				}
				return err
			}
			res.Close()

			return nil
		}
	}
}

// StopServer stops the Minecraft server over HTTP.
// This function doesn't wait until the server is fully stopped, it only requests a shutdown.
func (api *Api) StopServer() error {
	info, err := api.GetServerInfo()
	if err != nil {
		return err
	}

	if info.Status == Offline {
		return ServerAlreadyStoppedError
	}

	_, err = api.client.Get(fmt.Sprintf("panel/ajax/stop.php?SEC=%s&TOKEN=%s", api.sec, api.token))

	return err
}

// GetCookies returns the current authentication cookies that are being used.
//
// You can use this function to export them (to for example a .txt file) so you can resume the session later.
func (api *Api) GetCookies() []*http.Cookie {
	u, _ := url.Parse(api.client.Options.PrefixURL)
	return api.client.Options.CookieJar.Cookies(u)
}
