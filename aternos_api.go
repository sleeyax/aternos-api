package aternos_api

import (
	tls "github.com/refraction-networking/utls"
	"github.com/sleeyax/aternos-api/internal/tlsadapter"
	"github.com/sleeyax/gotcha"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type Api struct {
	options *Options
	client  *gotcha.Client
	// ajax security token.
	sec string
	// ajax token.
	token string
}

// New allocates a new Aternos API instance.
func New(options *Options) *Api {
	jar, _ := cookiejar.New(&cookiejar.Options{})

	adapter := tlsadapter.New(&tls.Config{ServerName: "aternos.org", InsecureSkipVerify: options.InsecureSkipVerify})

	client, _ := gotcha.NewClient(&gotcha.Options{
		Adapter:   adapter,
		PrefixURL: "https://aternos.org/",
		Headers: http.Header{
			"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"},
			"Accept":          {"*/*"},
			"Accept-Language": {"en-US,en;q=0.9"},
		},
		CookieJar:      jar,
		FollowRedirect: false,
		Retry:          false,
		Hooks: gotcha.Hooks{
			AfterResponse: []gotcha.AfterResponseHook{
				func(response *gotcha.Response, retry gotcha.RetryFunc) (*gotcha.Response, error) {
					if location := response.Header.Get("location"); strings.Contains(location, "go") {
						return response, UnauthenticatedError
					}
					return response, nil
				},
			},
		},
		Proxy: options.Proxy,
	})

	u, _ := url.Parse(client.Options.PrefixURL)
	jar.SetCookies(u, options.Cookies)

	return &Api{
		options: options,
		client:  client,
	}
}
