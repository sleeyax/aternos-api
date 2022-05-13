package tlsadapter

import (
	"bufio"
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"net/url"
)

type httpProxyDialer struct {
	// HTTP proxy URL to connect to.
	proxyURL *url.URL
}

// Dial tunnels traffic through the proxy to the destination network:address.
func (hpd *httpProxyDialer) Dial(network string, address string) (net.Conn, error) {
	conn, err := net.Dial(network, hpd.proxyURL.Host)
	if err != nil {
		return nil, err
	}

	connectHeader := make(http.Header)
	if user := hpd.proxyURL.User; user != nil {
		proxyUser := user.Username()
		if proxyPassword, passwordSet := user.Password(); passwordSet {
			credential := base64.StdEncoding.EncodeToString([]byte(proxyUser + ":" + proxyPassword))
			connectHeader.Set("Proxy-Authorization", "Basic "+credential)
		}
	}

	connectReq := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Opaque: address},
		Host:   address,
		Header: connectHeader,
	}

	if err = connectReq.Write(conn); err != nil {
		conn.Close()
		return nil, err
	}

	res, err := http.ReadResponse(bufio.NewReader(conn), connectReq)
	if err != nil {
		conn.Close()
		return nil, err
	}

	if res.StatusCode != 200 {
		conn.Close()
		return nil, errors.New(res.Status)
	}

	return conn, nil
}
