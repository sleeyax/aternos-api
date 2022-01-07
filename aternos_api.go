package aternos_api

import (
	"github.com/sleeyax/gotcha"
	"github.com/sleeyax/gotcha/adapters/fhttp"
	httpx "github.com/useflyent/fhttp"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type AternosApi struct {
	options *Options
	client  *gotcha.Client
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
		options: options,
		client:  client,
	}
}
