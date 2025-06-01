package sse

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Event SSE标准事件类
type Event struct {
	// 事件类型，默认为message
	Event string `json:"event"`
	// 数据
	Data any `json:"data"`
	// 事件ID
	ID string `json:"id,omitempty"`
	// 重试时间
	Retry int `json:"retry,omitempty"`
}

type Client struct {
	ID      string
	RoomID  string
	Message chan Event
}

func (c *Client) Loop(ctx *gin.Context) {
	// 监听客户端断开
	notify := ctx.Done()
	writer := ctx.Writer

	// 设置SSE响应头
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := writer.(http.Flusher)
	if !ok {
		log.Printf("Streaming not supported")
		ctx.JSON(http.StatusInternalServerError, "Streaming not supported")
		return
	}

	// 持续发送
	for {
		select {
		case <-notify:
			return // 退出事件循环
		case event := <-c.Message:
			// 将事件转换为SSE格式
			if event.Event != "" {
				fmt.Fprintf(writer, "event: %s\n", event.Event)
			}
			if event.ID != "" {
				fmt.Fprintf(writer, "id: %s\n", event.ID)
			}
			if event.Retry != 0 {
				fmt.Fprintf(writer, "retry: %d\n", event.Retry)
			}

			data, err := json.Marshal(event.Data)
			if err != nil {
				log.Printf("Error marshaling event data: %v", err)
				continue
			}

			fmt.Fprintf(writer, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
