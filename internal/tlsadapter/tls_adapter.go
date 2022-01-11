package tlsadapter

import (
	"context"
	tls "github.com/refraction-networking/utls"
	"github.com/sleeyax/gotcha"
	fhttpadapter "github.com/sleeyax/gotcha/adapters/fhttp"
	fhttp "github.com/useflyent/fhttp"
	"net"
	"net/url"
)

// TLSAdapter implements a custom gotcha.Adapter with advanced TLS options.
type TLSAdapter struct {
	// TLS fingerprint to use.
	// Defaults to tls.HelloCustom.
	Fingerprint tls.ClientHelloID

	// Optional custom client hello specification to use.
	// When specified, Fingerprint should be set to tls.HelloCustom.
	// Defaults to a custom Chrome 96 spec.
	ClientHelloSpec *tls.ClientHelloSpec

	// Optional TLS configuration to use.
	Config *tls.Config
}

// New creates a new gotcha adapter configured with a Chrome 96 browser TLS fingerprint.
func New(config *tls.Config) *TLSAdapter {
	return &TLSAdapter{Fingerprint: tls.HelloCustom, ClientHelloSpec: CustomClientHelloSpec, Config: config}
}

// DoRequest executes a HTTP 1 request and returns its response.
func (ua *TLSAdapter) DoRequest(options *gotcha.Options) (*gotcha.Response, error) {
	transport := &fhttp.Transport{
		Proxy: func(*fhttp.Request) (*url.URL, error) {
			return options.Proxy, nil
		},
		// TODO: fix
		//DialTLSContext:      ua.DialTLSContext,
		MaxConnsPerHost:     1,
		MaxIdleConns:        1,
		MaxIdleConnsPerHost: 1,
	}

	adapter := fhttpadapter.Adapter{Transport: transport}

	return adapter.DoRequest(options)
}

func (ua *TLSAdapter) DialTLSContext(_ context.Context, network string, addr string) (net.Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	uconn := tls.UClient(conn, ua.Config, ua.Fingerprint)

	if ua.Fingerprint == tls.HelloCustom {
		if err = uconn.ApplyPreset(ua.ClientHelloSpec); err != nil {
			return nil, err
		}
	}

	if err = uconn.Handshake(); err != nil {
		return nil, err
	}

	return uconn, err
}
