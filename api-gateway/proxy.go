package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func newReverseProxy(target string) *httputil.ReverseProxy {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", url.Host)
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		req.URL.Path = req.URL.Path
	}

	return proxy
}
