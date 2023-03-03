package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewMultipleHostsReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		fmt.Println("收到请求", req)
		target := targets[rand.Int()%len(targets)]
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}
	return &httputil.ReverseProxy{Director: director}
}
func main() {
	proxy := NewMultipleHostsReverseProxy([]*url.URL{
		{
			Scheme: "http",
			Host:   "192.168.11.207:3000",
		},
	})
	log.Fatal(http.ListenAndServe(":9091", proxy))
}
