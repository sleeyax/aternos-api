package aternos_api

import (
	"net/http"
	"net/url"
)

type Options struct {
	// Initial authorization cookies.
	//
	// It must include at least ATERNOS_SESSION, ATERNOS_SERVER.
	//
	// It's recommended to also specify ATERNOS_LANGUAGE.
	Cookies []*http.Cookie

	// Optional HTTP proxy to use.
	Proxy *url.URL
}
