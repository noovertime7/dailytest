package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	ConnectionCount = 0
	EmptyRegistry   = prometheus.NewRegistry()
)

func init() {
	EmptyRegistry.MustRegister(cc)
}

// 带动态标签的counter
var cc = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "test",
		Name:      "connection_count_with_label",
	},
	[]string{"app", "namespace"},
)

func main() {
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		ConnectionCount++
		cc.With(prometheus.Labels{
			"app":       "simple-counter",
			"namespace": "test",
		}).Set(float64(ConnectionCount))

		c := fmt.Sprintf("%d\n", ConnectionCount)
		writer.Write([]byte("count: " + c))
	})

	// 以下两种写法均可
	////写法一
	//http.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
	//	promhttp.HandlerFor(EmptyRegistry,
	//		promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}).
	//		ServeHTTP(writer, request)
	//
	//})

	// 写法二
	http.Handle("/metrics", promhttp.HandlerFor(EmptyRegistry,
		promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}))
	//
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}

}
