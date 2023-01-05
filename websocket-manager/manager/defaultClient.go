package manager

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/noovertime7/dailytest/utils"
	"log"
)

func DefaultWsHandlerRegistration(manager *Manager) {
	manager.Registration(NewDefaultWsHandler())
}

var _ WSHandler = &DefaultWsHandler{}

const defaultWsHandlerName = "DefaultWsHandler"

// DefaultWsHandler 单个 websocket 信息
type DefaultWsHandler struct {
	Context context.Context
	Socket  *websocket.Conn
	Message chan []byte
}

func NewDefaultWsHandler() *DefaultWsHandler {
	return &DefaultWsHandler{Message: make(chan []byte, 128)}
}

func (c *DefaultWsHandler) GetHandlerName() string {
	return defaultWsHandlerName
}

func (c *DefaultWsHandler) Close() error {
	return c.Socket.Close()
}

func (c *DefaultWsHandler) SetUp(ctx *gin.Context) error {
	ws, err := UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return err
	}
	c.Context = ctx
	c.Socket = ws
	// 启动监听
	go c.Read()
	go c.Write()
	return err
}

// 读信息，从 websocket 连接直接读取数据
func (c *DefaultWsHandler) Read() {
	defer func(stopCh <-chan struct{}) {
		<-stopCh
		_ = c.Close()
	}(c.Context.Done())
	for {
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}
		c.Message <- message
	}
}

// 写信息，从 channel 变量 Send 中读取数据写入 websocket 连接
func (c *DefaultWsHandler) Write() {
	for {
		select {
		case <-c.Context.Done():
			break
		case message, ok := <-c.Message:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			data, err := utils.ParseToJsonByte(string(message))
			if err != nil {
				log.Printf("client writemessage err: %s", err)
			}
			err = c.Socket.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("client writemessage err: %s", err)
			}
		}
	}
}
