package models

// binding的逗号之间不能有空格

type CommentData struct {
	CommentID  string `json:"comment_id"`
	Content    string `json:"content"`
	UserID     string `json:"user_id"`
	UserName   string `json:"user_name"`
	CreateTime string `json:"create_time"`
}

// 新评论表单
type AddCommentForm struct {
	RoomId   int64  `json:"room_id" form:"room_id" binding:"required"`
	UserId   int64  `json:"user_id" form:"user_id" binding:"required"`
	UserName string `json:"user_name" form:"user_name"`
	Content  string `json:"content" form:"content" binding:"required"`
	Stamp    int64  `json:"stamp" form:"stamp"`
}
