package sse

import (
	"Comment/module/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

var (
	instance *Server
	once     sync.Once
)

func Instance() *Server {
	once.Do(func() {
		instance = &Server{
			eventChan:  make(chan Event),
			register:   make(chan *Client),
			unregister: make(chan string),
			hub:        NewHub(),
		}
	})
	return instance
}

type Server struct {
	eventChan  chan Event
	register   chan *Client
	unregister chan string

	hub *Hub
}

func (s *Server) Start() {
	go func() {
		for {
			select {
			// 注册
			case client := <-s.register:
				logrus.Infof("SSE Client %s Connected", client.ID)
				s.hub.AddClient(client)
			// 注销
			case clientID := <-s.unregister:
				logrus.Infof("SSE Client %s Disconnected", clientID)
				s.hub.RemoveClient(clientID)
			// todo 服务器事件
			case event := <-s.eventChan:

			}
		}
	}()
}

func (s *Server) AppendEvent(e *Event) {
	s.eventChan <- *e
}

func (s *Server) Handler(ctx *gin.Context) {
	// 确保是SSE请求
	if ctx.Request.Header.Get("Accept") != "text/event-stream" {
		logrus.Infof("Invalid Accept header: %s", ctx.Request.Header.Get("Accept"))
		ctx.JSON(http.StatusBadRequest, "This endpoint only supports EventStream")
		return
	}

	// 创建新客户端
	clientID, err := utils.GenerateUUID()
	if err != nil {
		logrus.Errorf("uuid生成失败:%s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "服务器内部错误",
		})
	}
	client := &Client{
		ID:      clientID,
		Message: make(chan Event, 10), // 缓冲防止阻塞
	}

	// 注册客户端
	s.register <- client
	go func() {
		client.Loop(ctx)
		//  退出注销该客户端
		s.unregister <- client.ID
	}()
}
