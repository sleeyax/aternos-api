package aternos_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// getDocument sends a GET request to the specified url and reads the response as a goquery.Document.
func (api *AternosApi) getDocument(url string) (*goquery.Document, error) {
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
func (api *AternosApi) genSec() {
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
func (api *AternosApi) extractToken(document *goquery.Document) error {
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
func (api *AternosApi) GetServerInfo() (ServerInfo, error) {
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
func (api *AternosApi) StartServer() error {
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

// ConfirmServer sends a confirmation over HTTP that says 'you're still active' while waiting for the server to start.
// You should call this function right after starting the server.
//
// delay specifies the amount of seconds to wait before submitting the confirmation.
// Recommended to wait ~10 seconds.
func (api *AternosApi) ConfirmServer(delay time.Duration) error {
	for {
		time.Sleep(delay)

		info, err := api.GetServerInfo()
		if err != nil {
			return err
		}

		status := info.Status

		if status != Online {
			res, err := api.client.Get(fmt.Sprintf("panel/ajax/confirm.php?headstart=0&access-credits=0&SEC=%s&TOKEN=%s", api.sec, api.token))
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

// StopServer stops the Minecraft server over HTTP.
// This function doesn't wait until the server is fully stopped, it only requests a shutdown.
func (api *AternosApi) StopServer() error {
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
func (api *AternosApi) GetCookies() []*http.Cookie {
	u, _ := url.Parse(api.client.Options.PrefixURL)
	return api.client.Options.CookieJar.Cookies(u)
}
