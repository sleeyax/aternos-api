package tlsadapter

import (
	"context"
	"errors"
	"fmt"
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

	// Amount of times to retry the TCP/TLS connection on failure.
	// Defaults to 3 times.
	Retries int
}

// New creates a new gotcha adapter configured with a Chrome 96 browser TLS fingerprint.
func New(config *utls.Config) *TLSAdapter {
	return &TLSAdapter{Fingerprint: utls.HelloCustom, Config: config, Retries: 3}
}

// DoRequest executes a HTTP 1 request and returns its response.
func (ua *TLSAdapter) DoRequest(options *gotcha.Options) (*gotcha.Response, error) {
	transport := &fhttp.Transport{
		// NOTE: setting proxy on the Transport is currently broken, see: https://github.com/sleeyax/gotcha/commit/4b06cd561da906d0a570901e90b5bb5c313c1f1b.
		// We'll use DialTLSContext to connect to the proxy instead.
		// Proxy: fhttp.ProxyURL(options.Proxy),
		DialTLSContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
			return ua.ProxyTLS(network, addr, options.Proxy)
		},
	}

	adapter := fhttpadapter.Adapter{Transport: transport}

	return adapter.DoRequest(options)
}

func (ua *TLSAdapter) connectTLS(conn net.Conn) (net.Conn, error) {
	config := ua.Config.Clone()

	uconn := utls.UClient(conn, config, ua.Fingerprint)

	if ua.Fingerprint == utls.HelloCustom {
		if err := uconn.ApplyPreset(getCustomClientHelloSpec()); err != nil {
			return nil, err
		}
	}

	if err := uconn.Handshake(); err != nil {
		return nil, err
	}

	return uconn, nil
}

func (ua *TLSAdapter) ProxyTLS(network string, addr string, proxy *url.URL) (net.Conn, error) {
	for i := 1; i <= ua.Retries; i += 1 {
		var conn net.Conn
		var err error

		if proxy != nil {
			proxyDialer := HttpProxyDialer{ProxyURL: proxy}
			conn, err = proxyDialer.Dial(network, addr)
		} else {
			conn, err = net.Dial(network, addr)
		}

		if err != nil {
			// return nil, err
			fmt.Printf("adapter: proxy dial error (%v), retrying %d/%d\n", err, i, ua.Retries)
			continue
		}

		tlsConn, err := ua.connectTLS(conn)
		if err != nil {
			// return nil, err
			fmt.Printf("adapter: proxy TLS error (%v), retrying %d/%d\n", err, i, ua.Retries)
			conn.Close() // NOTE:  if the TLS conn failed it's safe to assume the proxy is bad, so we just connect to a new one by going to the beginning of the loop
			continue
		}

		return tlsConn, nil
	}

	return nil, errors.New("adapter: max retries reached")
}
