package models

// Response 统一返回格式
type Response struct {
	Code    int    `json:"code"` // 业务状态码
	Message string `json:"message"`
	Data    any    `json:"data"`
}

const (
	Success = 0
	Failed  = 1
)
