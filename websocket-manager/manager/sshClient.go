package manager

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/noovertime7/dailytest/utils"
	"github.com/noovertime7/dailytest/websocket-manager/manager/termical"
	"log"
	"strconv"
	"sync"
	"time"
)

const sshWsHandlerName = "sshWsHandler"

func SSHWsHandlerRegistration(manager *Manager) {
	manager.Registration(&sshWsHandler{})
}

type RecordData struct {
	Event string  `json:"event"` // 输入输出事件
	Time  float64 `json:"time"`  // 时间差
	Data  []byte  `json:"data"`  // 数据
}

type Meta struct {
	TERM      string
	Width     int
	Height    int
	UserName  string
	ConnectId string
	HostId    uint
	HostName  string
}

type sshWsHandler struct {
	sync.RWMutex
	Terminal    *termical.Terminal // ssh客户端
	Conn        *websocket.Conn    // socket 连接
	messageType int                // 发送的数据类型
	recorder    []*RecordData      // 操作记录
	CreatedAt   time.Time          // 创建时间
	UpdatedAt   time.Time          // 最新的更新时间
	Meta        Meta               // 元信息
	written     bool               // 是否已写入记录, 一个流只允许写入一次
}

func (s *sshWsHandler) Read() {
	return
}

func (s *sshWsHandler) Write() {
	return
}

func (s *sshWsHandler) Close() error {
	return s.Conn.Close()
}

func (s *sshWsHandler) GetHandlerName() string {
	return sshWsHandlerName
}

func (s *sshWsHandler) SetUp(ctx *gin.Context) error {
	// 设置默认xterm窗口大小
	cols, _ := strconv.Atoi(ctx.DefaultQuery("cols", "188"))
	rows, _ := strconv.Atoi(ctx.DefaultQuery("rows", "42"))
	terminalConfig := termical.Config{
		IpAddress: "yunxue521.top",
		Port:      "22",
		UserName:  "root",
		Password:  "1qaz@WSXchenteng@",
		Width:     cols,
		Height:    rows,
	}
	ws, err := UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return err
	}
	terminal, err := termical.NewTerminal(terminalConfig)
	if err != nil {
		_ = ws.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
		_ = ws.Close()
		return err
	}
	resizeCh := make(chan interface{}, 10)
	stream, err := termical.NewTerminalSession(resizeCh, ws)
	go func() {
		for {
			if terminal.IsClosed() {
				return
			}
			select {
			case data := <-resizeCh:
				msg := data.(termical.TerminalMessage)
				fmt.Printf("监听到resize信号%v\n", msg)
				terminal.SetWinSize(msg.Cols, msg.Rows)
			}
		}
	}()
	if err != nil {
		log.Printf("NewTerminalSession error %v\n", err)
		return err
	}

	err = terminal.Connect(stream, stream, stream)

	if err != nil {
		_ = ws.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
		_ = ws.Close()
		return err
	}

	// 断开ws和ssh的操作
	stream.WsConn.SetCloseHandler(func(code int, text string) error {
		if err := stream.Close(); err != nil {
			return err
		}
		if err := terminal.Close(); err != nil {
			return err
		}
		log.Printf("ws断开成功")
		return nil
	})

	go func() {
		for {
			// 每5秒
			timer := time.NewTimer(5 * time.Second)
			<-timer.C

			if stream.IsClosed() || terminal.IsClosed() {
				_ = timer.Stop()
				break
			}
			// 如果有 10 分钟没有数据流动，则断开连接
			if time.Now().Unix()-stream.UpdateAt.Unix() > 60*10 {
				stream.WsConn.WriteMessage(websocket.TextMessage, utils.Str2Bytes("检测到终端闲置，已断开连接...\r\n"))
				stream.WsConn.WriteMessage(websocket.BinaryMessage, utils.Str2Bytes("检测到终端闲置，已断开连接..."))
				stream.Close()
				_ = timer.Stop()
				break
			}
		}
	}()

	return nil
}
