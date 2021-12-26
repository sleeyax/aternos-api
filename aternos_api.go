package aternos_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/sleeyax/gotcha"
	"github.com/sleeyax/gotcha/adapters/fhttp"
	httpx "github.com/useflyent/fhttp"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type AternosApi struct {
	*Options
	client *gotcha.Client
	// ajax security token.
	sec string
	// ajax token.
	token string
}

// New allocates a new Aternos API instance.
func New(options *Options) *AternosApi {
	jar, _ := cookiejar.New(&cookiejar.Options{})

	var adapter gotcha.Adapter
	if options.Proxy != nil {
		adapter = &fhttp.Adapter{Transport: &httpx.Transport{
			Proxy: httpx.ProxyURL(options.Proxy),
		}}
	} else {
		adapter = fhttp.NewAdapter()
	}

	client, _ := gotcha.NewClient(&gotcha.Options{
		Adapter:   adapter,
		PrefixURL: "https://aternos.org/",
		Headers: http.Header{
			"user-agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"},
			"accept":             {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
			"accept-language":    {"en-US,en;q=0.9"},
			"accept-encoding":    {"gzip, deflate, br"},
			"cookie":             nil,
			httpx.HeaderOrderKey: []string{"user-agent", "accept", "accept-language", "accept-encoding", "cookie"},
		},
		CookieJar:      jar,
		FollowRedirect: false,
		Retry:          false,
	})

	u, _ := url.Parse(client.Options.PrefixURL)
	jar.SetCookies(u, options.Cookies)

	return &AternosApi{
		Options: options,
		client:  client,
	}
}

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

// GetServerInfo returns server information.
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

// TODO: start & stop server over websockets (wss://aternos.org/hermes/)

// StartServer starts your Minecraft server.
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

// ConfirmServer sends a confirmation that 'you're still active' while waiting for the server to start.
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

		if status != Preparing && status != Online {
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

// StopServer stops the minecraft server.
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
