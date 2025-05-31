package sse

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	clients    map[string]*Client
	mu         sync.Mutex
	eventChan  chan Event
	register   chan *Client
	unregister chan string
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[string]*Client),
		eventChan:  make(chan Event),
		register:   make(chan *Client),
		unregister: make(chan string),
	}
}

func (s *Server) AppendEvent(e *Event) {
	s.eventChan <- *e
}

func (s *Server) Run() {
	go func() {
		for {
			select {
			// 注册
			case client := <-s.register:
				s.mu.Lock()
				s.clients[client.ID] = client
				s.mu.Unlock()
				log.Printf("Client %s connected", client.ID)
			// 注销
			case clientID := <-s.unregister:
				s.mu.Lock()
				if client, ok := s.clients[clientID]; ok {
					close(client.Message)
					delete(s.clients, clientID)
					log.Printf("Client %s disconnected", clientID)
				}
				s.mu.Unlock()
			// 服务器事件
			case event := <-s.eventChan:
				s.mu.Lock()
				// 广播所有客户端
				for _, client := range s.clients {
					select {
					case client.Message <- event:
					default:
						// 如果无法发送，认为客户端已断开
						close(client.Message)
						delete(s.clients, client.ID)
					}
				}
				s.mu.Unlock()
			}
		}
	}()
}

func (s *Server) Handler(ctx *gin.Context) {
	w := ctx.Writer
	r := ctx.Request
	// 确保是SSE请求
	if r.Header.Get("Accept") != "text/event-stream" {
		log.Printf("Invalid Accept header: %s", r.Header.Get("Accept"))
		ctx.JSON(http.StatusBadRequest, "This endpoint only supports EventStream")
		return
	}

	// 设置SSE响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Printf("Streaming not supported")
		ctx.JSON(http.StatusInternalServerError, "Streaming not supported")
		return
	}

	// 创建新客户端
	clientID := fmt.Sprintf("%d", time.Now().UnixNano())
	client := &Client{
		ID:      clientID,
		Message: make(chan Event, 10), // 缓冲防止阻塞
	}

	// 注册客户端
	s.register <- client
	defer func() {
		s.unregister <- client.ID
	}()

	// 监听客户端断开
	notify := ctx.Done()

	// 初始欢迎消息
	client.Message <- Event{
		Event: "welcome",
		Data:  "Connected to SSE server",
		ID:    clientID,
		Retry: 5000,
	}

	// 持续发送事件
	for {
		select {
		case <-notify:
			return // 退出，触发defer注销该客户端
		case event := <-client.Message:
			// 将事件转换为SSE格式
			if event.Event != "" {
				fmt.Fprintf(w, "event: %s\n", event.Event)
			}
			if event.ID != "" {
				fmt.Fprintf(w, "id: %s\n", event.ID)
			}
			if event.Retry != 0 {
				fmt.Fprintf(w, "retry: %d\n", event.Retry)
			}

			data, err := json.Marshal(event.Data)
			if err != nil {
				log.Printf("Error marshaling event data: %v", err)
				continue
			}

			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
