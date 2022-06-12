package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	"shop-srvs/goods_srv/proto/proto"
)

//轮播图
func (s *GoodsServer) BannerList(context.Context, *emptypb.Empty) (*proto.BannerListResponse, error) {
	return nil, nil
}
func (s *GoodsServer) CreateBanner(context.Context, *proto.BannerRequest) (*proto.BannerResponse, error) {
	return nil, nil
}
func (s *GoodsServer) DeleteBanner(context.Context, *proto.BannerRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *GoodsServer) UpdateBanner(context.Context, *proto.BannerRequest) (*emptypb.Empty, error) {
	return nil, nil
}
