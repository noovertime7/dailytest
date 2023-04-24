package main

import (
	"github.com/gin-gonic/gin"
	"github.com/noovertime7/dailytest/prometheus_proxy_demo/handler"
	"log"
)

func main() {
	r := gin.Default()

	prometheusHandler, err := handler.NewPrometheusHandler("http://39.106.52.41:30010")
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/proxy/thanos/matrix", prometheusHandler.Matrix)

	r.Run(":9090")
}
