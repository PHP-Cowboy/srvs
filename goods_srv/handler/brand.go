package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	"shop-srvs/goods_srv/proto/proto"
)

//品牌和轮播图
func (s *GoodsServer) BrandList(context.Context, *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	return nil, nil
}
func (s *GoodsServer) CreateBrand(context.Context, *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	return nil, nil
}
func (s *GoodsServer) DeleteBrand(context.Context, *proto.BrandRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *GoodsServer) UpdateBrand(context.Context, *proto.BrandRequest) (*emptypb.Empty, error) {
	return nil, nil
}
