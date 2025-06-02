package biz

import (
	sse "Comment/app/sse/internal/sse_server"
	"Comment/proto/pb"
	"context"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SSEServiceHandler struct {
	pb.UnimplementedSSEServiceServer
}

func (receiver *SSEServiceHandler) Report(ctx context.Context, in *pb.ReportRequest) (*emptypb.Empty, error) {
	data, err := proto.Marshal(in.Comment)
	if err != nil {
		return nil, err
	}
	sse.Instance().PushMsg(in.Comment.RoomId, data)

	return &emptypb.Empty{}, nil
}
