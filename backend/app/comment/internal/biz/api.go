package biz

import (
	"Comment/proto/pb"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CommentHandler struct {
	pb.UnimplementedCommentServiceServer
}

func (receiver *CommentHandler) CreateComment(ctx context.Context, in *pb.CreateCommentRequest) (*emptypb.Empty, error) {
	// 存储

	// 调用sse服务，广播

	return &emptypb.Empty{}, nil
}
