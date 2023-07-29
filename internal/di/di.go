package di

import (
	"cgg/api/pb"
	"cgg/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Wire(s *grpc.Server) {
	pb.RegisterChatGTPServiceServer(s, &service.ChatGTP{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
}
