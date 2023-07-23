package main

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/koding/websocketproxy"
)

var (
	xremote, xremoteWs, xremoteHost, xurl, xenv string
)

func HandleHttpRequest(w http.ResponseWriter, r *http.Request) {

	r.Host = xremoteHost
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
	}

	if len(xurl) != 0 && r.URL.Path != xurl { //非指定url拒绝访问
		return
	}

	if strings.ToLower(r.Header.Get("Upgrade")) == "websocket" {

		remote, err := url.Parse(xremoteWs)
		if err != nil {
			return
		}
		proxy := websocketproxy.NewProxy(remote)
		proxy.Dialer = &websocket.Dialer{
			TLSClientConfig:  tlsConf,
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
		}
		proxy.Upgrader = &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		proxy.ServeHTTP(w, r)
		return
	}

	remote, err := url.Parse(xremote)
	if err != nil {
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Transport = &http.Transport{
		TLSClientConfig: tlsConf,
	}
	proxy.ServeHTTP(w, r)
	return
}

func main() {
	switch xenv {
	case "bd":
		//todo
	default: //ali
		http.HandleFunc("/", HandleHttpRequest)
		http.ListenAndServe("0.0.0.0:9000", nil)
	}

}

func init() {
	xremote, xremoteWs, xremoteHost, xurl, xenv = os.Getenv("XREMOTE"),
		os.Getenv("XREMOTEWs"), os.Getenv("XREMOTE_HOST"), os.Getenv("XURL"), os.Getenv("XENV")
}
