package biz

import (
	"Comment/app/gateway/internal/models"
	"Comment/module/services"
	"Comment/proto/pb"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CommentHandler struct {
}

// 新增评论接口
func (receiver *CommentHandler) HandleAddComment(ctx *gin.Context) {
	// 参数校验
	form, ok := validate(ctx, &models.AddCommentForm{})
	if !ok {
		return
	}
	resp := models.Response{}
	defer ctx.JSON(http.StatusOK, resp)

	// 调用评论服务
	conn, err := services.Instance.GetConnection("comment")
	if err != nil {
		logrus.Errorf("获取评论服务连接失败:%s", err.Error())
		resp.Code = models.Failed
		resp.Message = "评论服务连接失败"
		return
	}
	client := pb.NewCommentServiceClient(conn)
	req := &pb.CreateCommentRequest{
		Comment: &pb.Comment{
			UserId:    form.UserId,
			UserName:  form.UserName,
			Content:   form.Content,
			PubStamp:  form.Stamp,
			PubRegion: "",
			RoomId:    form.RoomId,
		},
	}
	_, err = client.CreateComment(ctx, req)
	if err != nil {
		resp.Code = models.Failed
		resp.Message = err.Error()
		return
	}
	resp.Code = models.Success
	resp.Message = "评论发送成功"
}
