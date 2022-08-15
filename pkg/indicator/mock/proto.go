package mock

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"static-analyze/internal/app/proto"
)

type Proto struct {
}

func (m Proto) GetHistory(ctx context.Context, in *proto.RequestMessage, opts ...grpc.CallOption) (*proto.ResponseMessage, error) {
	return &proto.ResponseMessage{
		P: []*proto.Pair{
			{Time: timestamppb.Now(), Value: 0.2},
			{Time: timestamppb.Now(), Value: 0.112},
		},
	}, nil
}
