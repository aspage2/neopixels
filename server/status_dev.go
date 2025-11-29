//go:build !arm

package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
)

// When not building on arm, proxies all status requests
// to another server.
type StatusHandler struct {
	forwardaddr *url.URL
}

func (this *StatusHandler) InitializeServer(sm *http.ServeMux) {
	if sm == nil {
		sm = http.DefaultServeMux
	}
	proxy := &httputil.ReverseProxy{
		Director: this.proxydirector,
	}
	http.Handle("/lights/", proxy)
	
	fmt.Printf("Initialized DEV proxy to the Raspberry Pi at %s. All /lights/* calls will be proxied do this server.", this.forwardaddr)
}

func (this *StatusHandler) proxydirector(req *http.Request) {
	oldPath := req.URL.Path
	*req.URL = *this.forwardaddr
	req.URL.Path = path.Join(this.forwardaddr.Path, oldPath)
	if strings.HasSuffix(oldPath, "/") {
		req.URL.Path += "/"
	}
}

func NewStatusHandler(
	numPixels int,
	pixelOrder string,
	brightness float32,
	xfade int,
) (*StatusHandler, error) {
	forwardaddr := os.Getenv("NEOPIXEL_FORWARD_ADDR")
	if forwardaddr == "" {
		return nil, errors.New("on a development machine, must specify host with NEOPIXEL_FORWARD_ADDR")
	}
	u, err := url.Parse(forwardaddr)
	if err != nil {
		return nil, err
	}
	checkU := new(url.URL)
	*checkU = *u
	// Have to add the ending slash like this because path.Join cleans the url
	checkU.Path = path.Join(checkU.Path, "/lights/status")
	resp, err := http.Get(checkU.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("could not verify response from backend")
	}
	return &StatusHandler{u}, nil
}

