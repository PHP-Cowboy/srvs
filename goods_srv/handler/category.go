package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	"shop-srvs/goods_srv/proto/proto"
)

//商品分类
func (s *GoodsServer) GetAllCategorysList(context.Context, *emptypb.Empty) (*proto.CategoryListResponse, error) {
	return nil, nil
}

//获取子分类
func (s *GoodsServer) GetSubCategory(context.Context, *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	return nil, nil
}
func (s *GoodsServer) CreateCategory(context.Context, *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	return nil, nil
}
func (s *GoodsServer) DeleteCategory(context.Context, *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *GoodsServer) UpdateCategory(context.Context, *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	return nil, nil
}
