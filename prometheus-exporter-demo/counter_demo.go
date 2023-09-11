package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	ConnectionCount = 0
)

func init() {
	prometheus.MustRegister(cc)
	prometheus.MustRegister(cf)

}

// 带动态标签的counter
var cc = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "test",
		Name:      "connection_count_with_label",
	},
	[]string{"app", "namespace"},
)

// 不带标签
var cf = prometheus.NewCounterFunc(prometheus.CounterOpts{
	Namespace: "test",
	Name:      "connection_count",
}, func() float64 {
	return float64(ConnectionCount)
})

func main() {
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		cc.With(prometheus.Labels{
			"app":       "simple-counter",
			"namespace": "test",
		}).Inc()

		ConnectionCount++
		c := fmt.Sprintf("%d\n", ConnectionCount)
		writer.Write([]byte("count: " + c))
	})

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}

}
