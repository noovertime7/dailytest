package main

import (
	"fmt"
	"github.com/noovertime7/dailytest/informer_demo"
	"github.com/noovertime7/dailytest/proxy-apiserver/transport"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func main() {
	// 从kubeconfig文件加载配置
	_, config := informer_demo.GetClientSet()

	//// 获取API服务器的URL
	apiServerURL, err := url.Parse(config.Host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse API server URL: %v\n", err)
		os.Exit(1)
	}

	apiServerURL = &url.URL{
		Scheme: "https",
		Host:   "tjyw-k8s-api:6443",
		Path:   "/",
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(apiServerURL)

	// 设置反向代理的Transport为Kubernetes的RoundTripper
	http2configCopy := *config
	http2configCopy.WrapTransport = transport.NewDynamicImpersonatingRoundTripper
	http2configCopy.Host = apiServerURL.Host
	ts, err := rest.TransportFor(&http2configCopy)
	if err != nil {
		panic(err)
	}
	proxy.Transport = ts

	http.HandleFunc("/proxy/", func(writer http.ResponseWriter, request *http.Request) {
		request.URL.Path = request.URL.Path[len("/proxy/"):]

		sourceIPs := utilnet.SourceIPs(request)
		// 计算latency
		startTime := time.Now()

		proxy.ServeHTTP(writer, request)
		latency := time.Since(startTime)

		log.Printf("verb=%q host=%q endpoint=%q URI=%q latency=%v userAgent=%q srcIP=%v",
			request.Method,
			request.Host,
			request.URL.Path,
			request.RequestURI,
			latency, // Set the desired response value
			request.UserAgent(),
			sourceIPs, // Set the desired sourceIPs value
		)
	})

	// 启动反向代理服务器
	log.Fatal(http.ListenAndServe(":8082", nil))
}
