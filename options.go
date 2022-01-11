package aternos_api

import (
	"net/http"
	"net/url"
)

type Options struct {
	// Initial authentication cookies.
	//
	// It must include at least ATERNOS_SESSION, ATERNOS_SERVER.
	//
	// It's recommended to also specify ATERNOS_LANGUAGE.
	Cookies []*http.Cookie

	// Optional HTTP proxy to use.
	Proxy *url.URL

	// Disables server SSL certificate checks.
	// It's recommended to enable this only for debugging purposes, such as debugging traffic with a web debugging/HTTP/MITM proxy.
	InsecureSkipVerify bool
}
