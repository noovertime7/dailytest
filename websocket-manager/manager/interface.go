package manager

import "github.com/gin-gonic/gin"

type SocketHandler interface {
	Read()
	Write()
	Close() error
}

type WSHandler interface {
	SocketHandler
	GetHandlerName() string
	SetUp(ctx *gin.Context) error
}
