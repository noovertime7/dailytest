package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/noovertime7/dailytest/prometheus_proxy_demo/response"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"strconv"
	"time"
)

type prometheusHandler struct {
	client api.Client
}

func NewPrometheusHandler(server string) (*prometheusHandler, error) {
	client, err := api.NewClient(api.Config{Address: server})
	if err != nil {
		return nil, err
	}
	return &prometheusHandler{client: client}, nil
}

func (p *prometheusHandler) Matrix(c *gin.Context) {
	query := c.Query("query")
	start := c.Query("start")
	end := c.Query("end")
	step, _ := strconv.Atoi(c.Query("step"))

	v1api := v1.NewAPI(p.client)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	s, _ := time.Parse("2006-01-02T15:04:05Z", start)
	e, _ := time.Parse("2006-01-02T15:04:05Z", end)
	r := v1.Range{
		Start: s,
		End:   e,
	}

	if step > 0 {
		r.Step = time.Duration(step) * time.Second
	} else {
		// 不传step就动态控制
		r.Step = dynamicTimeStep(r.Start, r.End)
	}

	obj, _, err := v1api.QueryRange(ctx, query, r)
	if err != nil {
		log.Error(err)
		response.NotOK(c, err)
		return
	}

	response.OK(c, obj)
}

func dynamicTimeStep(start time.Time, end time.Time) time.Duration {
	interval := end.Sub(start)
	if interval < 30*time.Minute {
		return 30 * time.Second // 30 分钟以内，step为30s, 返回60个点以内
	} else {
		return interval / 60 // 返回60个点，动态step
	}
}
