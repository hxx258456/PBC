package service

import (
	"context"
	"pyramid/pyramid-manage/backend/app/backend/internal/biz"

	pb "pyramid/pyramid-manage/backend/api/backend/v1"
)

// GreeterService is a greeter service.
type GreeterService struct {
	pb.UnimplementedGreeterServer

	uc *biz.GreeterUseCase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUseCase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements backend.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &pb.HelloReply{Message: "Hello " + g.Hello}, nil
}
