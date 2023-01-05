package manager

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var (
	UpGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024 * 1024 * 10,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// Manager 所有 websocket 信息
type Manager struct {
	Lock     sync.RWMutex
	innerMap map[string]WSHandler
}

func NewManager() *Manager {
	return &Manager{
		Lock:     sync.RWMutex{},
		innerMap: make(map[string]WSHandler, 10),
	}
}

func (m *Manager) Close() error {
	for _, handler := range m.innerMap {
		if err := handler.Close(); err != nil {
			return err
		}
	}
	return nil
}

func WebSocketConnect(ctx *gin.Context) {
	manager := NewManager()
	DefaultWsHandlerRegistration(manager)
	SSHWsHandlerRegistration(manager)
	handlerName := ctx.Param("name")
	handler, err := manager.Get(handlerName)
	if err != nil {
		log.Print(err)
		return
	}
	if err := handler.SetUp(ctx); err != nil {
		log.Print(err)
		return
	}
}

func (m *Manager) Registration(h WSHandler) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	log.Printf("注册%s", h.GetHandlerName())
	m.innerMap[h.GetHandlerName()] = h
}

func (m *Manager) Get(name string) (WSHandler, error) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	handler, ok := m.innerMap[name]
	if !ok {
		return nil, fmt.Errorf("handler not register")
	}
	return handler, nil
}
