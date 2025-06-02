package biz

import (
	"Comment/app/comment/internal/mysql"
	"Comment/app/comment/internal/redis"
	"Comment/module/api/IpQuery"
	"Comment/module/services"
	"Comment/proto/pb"
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommentHandler struct {
	pb.UnimplementedCommentServiceServer
}

func (receiver *CommentHandler) CreateComment(ctx context.Context, in *pb.CreateCommentRequest) (*emptypb.Empty, error) {
	// 查询地区
	newComment := mysql.RoomComment{
		RoomId:   in.Comment.RoomId,
		UserId:   in.Comment.UserId,
		UserName: in.Comment.UserName,
		Content:  in.Comment.Content,
		PubStamp: in.Comment.PubStamp,
	}
	ipInfo, err := IpQuery.QueryIpGeolocation(in.PubIp)
	if err != nil {
		logrus.Infof("查询ip归属地失败,ip:%s,err:%s", in.PubIp, err.Error())
		newComment.PubRegion = "未知地区"
	} else {
		newComment.PubRegion = ipInfo.RegionName
	}

	// 存储
	db := mysql.DBSession
	result := db.Create(&newComment)
	if result.Error != nil {
		logrus.Errorf("存储评论失败,err:%s", result.Error.Error())
		return nil, result.Error
	}

	// 异步推送
	go receiver.pushComment(&newComment)

	return &emptypb.Empty{}, nil
}
func (receiver *CommentHandler) pushComment(comment *mysql.RoomComment) {
	// 查询相关的sse服
	sseServiceIDList, err := redis.QuerySseServiceByRoomId(comment.RoomId)
	if err != nil {
		logrus.Errorf("查询sse服务错误,err:%s", err.Error())
		return
	}
	// 调用sse服务推送评论
	for _, nodeId := range sseServiceIDList {
		conn, err := services.Instance.GetConnectionByNodeId("sse-service", nodeId)
		if err != nil {
			logrus.Errorf("获取sse服务连接失败,err:%s", err.Error())
			continue
		}
		client := pb.NewSSEServiceClient(conn)
		_, err = client.Report(context.Background(), &pb.ReportRequest{
			Comment: &pb.Comment{
				UserId:    comment.UserId,
				UserName:  comment.UserName,
				Content:   comment.Content,
				PubStamp:  comment.PubStamp,
				PubRegion: comment.PubRegion,
				RoomId:    comment.RoomId,
			},
		})
		if err != nil {
			logrus.Errorf("调用sse服务错误,err:%s", err.Error())
			continue
		}
	}
}
