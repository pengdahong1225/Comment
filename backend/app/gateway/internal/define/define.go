package define

type CommentData struct {
	CommentID  string `json:"comment_id"`
	Content    string `json:"content"`
	UserID     string `json:"user_id"`
	UserName   string `json:"user_name"`
	CreateTime string `json:"create_time"`
}
