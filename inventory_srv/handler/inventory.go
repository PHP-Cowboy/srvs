package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"shop-srvs/inventory_srv/proto/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (i *InventoryServer) SetInv(context.Context, *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	return nil, nil
}

func (i *InventoryServer) InvDetail(context.Context, *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	return nil, nil
}

func (i *InventoryServer) Sell(context.Context, *proto.SellInfo) (*emptypb.Empty, error) {
	return nil, nil
}

func (i *InventoryServer) Reback(context.Context, *proto.SellInfo) (*emptypb.Empty, error) {
	return nil, nil
}
