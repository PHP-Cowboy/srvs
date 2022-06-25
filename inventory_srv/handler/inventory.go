package handler

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop-srvs/inventory_srv/global"
	"shop-srvs/inventory_srv/model"

	"google.golang.org/protobuf/types/known/emptypb"

	"shop-srvs/inventory_srv/proto/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (i *InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	var inv model.Inventory
	result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	if result.Error != nil {
		return nil, result.Error
	}
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	result = global.DB.Save(&inv)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected <= 0 {
		return nil, status.Errorf(codes.Internal, "保存库存数据失败")
	}
	return &emptypb.Empty{}, nil
}

func (i *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected <= 0 {
		return nil, errors.New("没有库存信息")
	}

	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

func (i *InventoryServer) Sell(context.Context, *proto.SellInfo) (*emptypb.Empty, error) {
	return nil, nil
}

func (i *InventoryServer) Reback(context.Context, *proto.SellInfo) (*emptypb.Empty, error) {
	return nil, nil
}
