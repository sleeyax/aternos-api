package tlsadapter

import (
	"context"
	"crypto/tls"
	utls "github.com/refraction-networking/utls"
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
	Fingerprint utls.ClientHelloID

	// Optional TLS configuration to use.
	Config *utls.Config
}

// New creates a new gotcha adapter configured with a Chrome 96 browser TLS fingerprint.
func New(config *utls.Config) *TLSAdapter {
	return &TLSAdapter{Fingerprint: utls.HelloCustom, Config: config}
}

// DoRequest executes a HTTP 1 request and returns its response.
func (ua *TLSAdapter) DoRequest(options *gotcha.Options) (*gotcha.Response, error) {
	transport := &fhttp.Transport{
		Proxy: func(*fhttp.Request) (*url.URL, error) {
			return options.Proxy, nil
		},
		DialTLSContext:      ua.DialTLSContext,
		MaxConnsPerHost:     1,
		MaxIdleConns:        1,
		MaxIdleConnsPerHost: 1,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: ua.Config.InsecureSkipVerify,
		},
	}

	adapter := fhttpadapter.Adapter{Transport: transport}

	return adapter.DoRequest(options)
}

func (ua *TLSAdapter) ConnectTLSContext(_ context.Context, conn net.Conn) (net.Conn, error) {
	config := ua.Config.Clone()

	uconn := utls.UClient(conn, config, ua.Fingerprint)

	if ua.Fingerprint == utls.HelloCustom {
		if err := uconn.ApplyPreset(GetCustomClientHelloSpec()); err != nil {
			return nil, err
		}
	}

	if err := uconn.Handshake(); err != nil {
		return nil, err
	}

	return uconn, nil
}

func (ua *TLSAdapter) DialTLSContext(ctx context.Context, network string, addr string) (net.Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	return ua.ConnectTLSContext(ctx, conn)
}
