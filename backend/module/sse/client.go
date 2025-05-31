package sse

type Client struct {
	ID      string
	Message chan Event
}
