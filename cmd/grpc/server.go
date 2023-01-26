package grpc_server

import (
	"context"

	pb "github.com/emika-team/grpc-proto/line-oa/go"
	"github.com/emika-team/line-oa-manager/pkg/grpc/message"
)

type server struct {
	pb.UnimplementedLineOAMessageServer
}

func NewGRPCHandler() *server {
	return &server{}
}

func (s *server) SendMessage(ctx context.Context, in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	return message.SendMessage(ctx, in)
}

func (s *server) GetMessageContent(ctx context.Context, in *pb.GetMessageContentRequest) (*pb.GetMessageContentResponse, error) {
	return message.GetContent(ctx, in)
}
