package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
)

var (
	counter       int
	healthy       int
	lock          sync.Mutex
	emptyRegistry *prometheus.Registry
)

func init() {
	lock = sync.Mutex{}
	emptyRegistry = prometheus.NewRegistry()
	emptyRegistry.MustRegister(NewTestCollector())
}

type TestCollector struct {
	Desc []*prometheus.Desc
}

func NewTestCollector() *TestCollector {
	variableLabels := []string{"ns", "app"}
	constLabels := prometheus.Labels{
		"const_label": "true",
	}

	return &TestCollector{Desc: []*prometheus.Desc{
		// counter
		prometheus.NewDesc(
			"test_app_connection_count",
			"connection count",
			variableLabels,
			constLabels,
		),
		// gauage
		prometheus.NewDesc(
			"test_app_healthy",
			"connection count",
			variableLabels,
			constLabels,
		),
	}}
}

// 描述
func (this *TestCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, d := range this.Desc {
		ch <- d
	}
}

// 收集指标
func (this *TestCollector) Collect(ch chan<- prometheus.Metric) {
	m1, err := prometheus.NewConstMetric(this.Desc[0], prometheus.CounterValue, float64(counter),
		"test", "test-app",
	)
	if err != nil {
		panic(err)
	}
	m2, err := prometheus.NewConstMetric(this.Desc[1], prometheus.GaugeValue, float64(healthy),
		"test", "test-app",
	)
	if err != nil {
		panic(err)
	}
	ch <- m1
	ch <- m2
}

func main() {

	http.HandleFunc("/set-healthy", func(writer http.ResponseWriter, request *http.Request) {
		lock.Lock()
		defer lock.Unlock()
		healthy = 1
		_, _ = writer.Write([]byte(fmt.Sprintf("%d", healthy)))
	})

	http.HandleFunc("/set-unhealthy", func(writer http.ResponseWriter, request *http.Request) {
		lock.Lock()
		defer lock.Unlock()
		healthy = 0
		_, _ = writer.Write([]byte(fmt.Sprintf("%d", healthy)))
	})

	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		lock.Lock()
		defer lock.Unlock()
		counter++
		c := fmt.Sprintf("counnter: %d", counter)
		_, _ = writer.Write([]byte(c))
	})

	http.Handle("/metrics",
		promhttp.HandlerFor(emptyRegistry,
			promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}))

	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
