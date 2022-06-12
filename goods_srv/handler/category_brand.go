package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	"shop-srvs/goods_srv/proto/proto"
)

//品牌分类
func (s *GoodsServer) CategoryBrandList(context.Context, *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	return nil, nil
}

//通过category获取brands
func (s *GoodsServer) GetCategoryBrandList(context.Context, *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	return nil, nil
}
func (s *GoodsServer) CreateCategoryBrand(context.Context, *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	return nil, nil
}
func (s *GoodsServer) DeleteCategoryBrand(context.Context, *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *GoodsServer) UpdateCategoryBrand(context.Context, *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	return nil, nil
}
