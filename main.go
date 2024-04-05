package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/koding/websocketproxy"
)

var (
	xremote, xremoteWs, xremoteHost, xurl, xenv string
	xtimeout, xrBuf, xwBuf, xlistenPort         int64
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
			HandshakeTimeout: time.Duration(xtimeout) * time.Millisecond,
		}
		proxy.Upgrader = &websocket.Upgrader{
			ReadBufferSize:  int(xrBuf),
			WriteBufferSize: int(xwBuf),
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
		IdleConnTimeout: time.Duration(xtimeout) * time.Millisecond,
	}

	proxy.ServeHTTP(w, r)
}

func main() {

	switch xenv {
	case "bd":
		//todo
	default: //ali

		http.HandleFunc("/", HandleHttpRequest)
		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", xlistenPort), nil); err != nil {
			fmt.Printf("ListenAndServe fail:%v\n", err)
		}
	}
}

func init() {
	xremote, xremoteWs, xremoteHost, xurl, xenv = os.Getenv("XREMOTE"),
		os.Getenv("XREMOTEWS"), os.Getenv("XREMOTE_HOST"), os.Getenv("XURL"), os.Getenv("XENV")
	var err error
	xtimeout, err = strconv.ParseInt(os.Getenv("XTIMEOUT"), 10, 64)
	if err != nil {
		xtimeout = 3000
	}
	xrBuf, err = strconv.ParseInt(os.Getenv("XRBUF"), 10, 64)
	if err != nil {
		xrBuf = 8192
	}
	xwBuf, err = strconv.ParseInt(os.Getenv("XWBUF"), 10, 64)
	if err != nil {
		xwBuf = 8192
	}
	xlistenPort, err = strconv.ParseInt(os.Getenv("XLISTEN_PORT"), 10, 64)
	if err != nil {
		xlistenPort = 9000
	}
}
