package sse

import (
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

func TestSSEServer(t *testing.T) {
	sse_server := NewServer()
	sse_server.Start()

	// 模拟事件生产者
	go func() {
		for {
			data := "测试事件"

			event := &Event{
				Event: "update",
				Data:  data,
				ID:    "",
				Retry: 0,
			}
			sse_server.AppendEvent(event)
			time.Sleep(3 * time.Second)
		}
	}()

	r := gin.Default()
	r.GET("/sse", sse_server.Handler)
	r.Run(":8080")
}
