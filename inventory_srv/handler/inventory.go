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

//用户下单扣减库存
func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin()

	for _, info := range req.GoodsInvInfo {
		var inv model.Inventory
		result := global.DB.Where("goods = ?", info.GoodsId).First(&inv)
		if result != nil || result.RowsAffected <= 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "数据未找到")
		}

		if inv.Stocks < info.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "库存不足")
		}

		inv.Stocks -= info.Num

		tx.Save(&inv)
	}

	tx.Commit()
	return &emptypb.Empty{}, nil
}

//库存归还 订单超时归还
func (i *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin()

	for _, info := range req.GoodsInvInfo {
		var inv model.Inventory
		result := global.DB.Where("goods = ?", info.GoodsId).First(&inv)
		if result != nil || result.RowsAffected <= 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "数据未找到")
		}

		inv.Stocks += info.Num

		tx.Save(&inv)
	}

	tx.Commit()
	return &emptypb.Empty{}, nil
}
