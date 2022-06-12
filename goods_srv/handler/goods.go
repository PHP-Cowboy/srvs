package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	"shop-srvs/goods_srv/proto/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

func (s *GoodsServer) GoodsList(context.Context, *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	return nil, nil
}

//现在用户提交订单有多个商品，你得批量查询商品的信息吧
func (s *GoodsServer) BatchGetGoods(context.Context, *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	return nil, nil
}
func (s *GoodsServer) CreateGoods(context.Context, *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	return nil, nil
}
func (s *GoodsServer) DeleteGoods(context.Context, *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *GoodsServer) UpdateGoods(context.Context, *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *GoodsServer) GetGoodsDetail(context.Context, *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	return nil, nil
}
