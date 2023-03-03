package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewProxy takes target host and creates a reverse proxy
// NewProxy 拿到 targetHost 后，创建一个反向代理
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

// ProxyRequestHandler handles the http request using proxy
// ProxyRequestHandler 使用 proxy 处理请求
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("接收到请求", r.RequestURI)
		proxy.ServeHTTP(w, r)
	}
}

const test = `/api/logCenter/proxy/grafana/explore?orgId=1&kiosk=true&theme=light&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiNjk5NWRiY2QtNWZmMy00YTIyLTkyMzgtZTQ1Njg5MGVkNTA4IiwiSUQiOjEsIlVzZXJuYW1lIjoieXVud2VpIiwiTmlja05hbWUiOiJ5dW53ZWkiLCJBdXRob3JpdHlJZCI6MTExLCJUZW5hbnRJRCI6IjYwOWQ3MjQ1OWI2NzRmNmViYmFlYTFlNjI1NGUxMjNlIiwiZXhwIjoxNjc3MDMzNjE2LCJpc3MiOiJrdWJlbWFuYWdlIiwibmJmIjoxNjc2OTQ2MjE2fQ.HesceijFJPTpcG3zZxVAeNrUr-_KjplafDamKllR1tc&left={%22datasource%22:%229t4aWZA4k%22,%22queries%22:[{%22refId%22:%22A%22,%22datasource%22:{%22type%22:%22loki%22,%22uid%22:%229t4aWZA4k%22},%22editorMode%22:%22code%22,%22expr%22:%22{job=\%22kuboard/kuboard-promtail\%22,region=\%22ceshi\%22,pod=\%22kuboard-promtail-7h7xq\%22}%22,%22queryType%22:%22range%22}],%22range%22:{%22from%22:%221676873720493%22,%22to%22:%221676960120493%22}}`

func main() {
	// initialize a reverse proxy and pass the actual backend server url here
	// 初始化反向代理并传入真正后端服务的地址
	proxy, err := NewProxy("http://192.168.11.207:3000")
	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	// 使用 proxy 处理所有请求到你的服务
	http.HandleFunc("/api/logCenter/proxy/grafana/explore", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
