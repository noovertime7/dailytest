package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Err  string      `json:"err"`
	Data interface{} `json:"data"`
}

func NotOK(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, &Response{Err: err.Error()})
}

func OK(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = "success"
	}
	ctx.JSON(http.StatusOK, &Response{Data: data})
}
