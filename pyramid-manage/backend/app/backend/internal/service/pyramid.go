package service

import (
	"context"
	"pyramid/pyramid-manage/backend/app/backend/internal/biz"

	pb "pyramid/pyramid-manage/backend/api/backend/v1"
)

type PyramidService struct {
	pb.UnimplementedPyramidServer

	uc *biz.SysChainUseCase
}

func NewPyramidService(uc *biz.SysChainUseCase) *PyramidService {
	return &PyramidService{
		uc: uc,
	}
}
func (s *PyramidService) GetOrganization(ctx context.Context, req *pb.Organization) (*pb.Organization, error) {
	return s.uc.QueryOrganization(ctx, req)
	//return &pb.Organization{}, nil
}
func (s *PyramidService) CreateOrganization(ctx context.Context, req *pb.Organization) (*pb.BaseReply, error) {
	return s.uc.AddOrganization(ctx, req)
	//return &pb.BaseReply{}, nil
}
