package sse

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
