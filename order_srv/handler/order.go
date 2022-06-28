package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"shop-srvs/order_srv/global"
	"shop-srvs/order_srv/model"
	"shop-srvs/order_srv/proto/proto"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

func (*OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var (
		shopCart []*model.ShoppingCart
		rsp      proto.CartItemListResponse
	)

	result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCart)

	if result.Error != nil {
		return nil, result.Error
	}

	rsp.Total = int32(result.RowsAffected)

	for _, cart := range shopCart {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      cart.ID,
			UserId:  cart.User,
			GoodsId: cart.Goods,
			Nums:    cart.Nums,
			Checked: false,
		})
	}

	return &rsp, nil
}

func (*OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	var (
		shopCart model.ShoppingCart
	)

	result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Goods: req.GoodsId}).First(&shopCart)

	if result != nil {
		return nil, result.Error
	}

	//查询到更新 未查到新增
	if result.RowsAffected >= 1 {
		shopCart.Nums += req.Nums
	} else {
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}

	result = global.DB.Save(&shopCart)

	if result != nil {
		return nil, result.Error
	}

	return &proto.ShopCartInfoResponse{Id: shopCart.ID}, nil
}

func (*OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	var (
		shopCart model.ShoppingCart
	)

	result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Goods: req.GoodsId}).First(&shopCart)

	if result != nil {
		return nil, result.Error
	}

	if result.RowsAffected < 1 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}

	shopCart.Checked = req.Checked

	result = global.DB.Save(&shopCart)

	if result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}

func (*OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		新建订单
		1. 从购物车中获取到选中的商品
		2. 商品的价格自己查询 - 访问商品服务 (跨微服务)
		3. 库存的扣减 - 访问库存服务 (跨微服务)
		4. 订单的基本信息表 - 订单的商品信息表
		5. 从购物车中删除已购买的记录
	*/
	return nil, nil
}

func (*OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var (
		orders []model.OrderInfo
		total  int64
		rsp    proto.OrderListResponse
	)

	result := global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)

	if result.Error != nil {
		return nil, result.Error
	}

	result = global.DB.Where(&model.OrderInfo{User: req.UserId}).Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	rsp.Total = int32(total)

	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
		})
	}

	return &rsp, nil
}

func (*OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var (
		order      model.OrderInfo
		orderGoods []model.OrderGoods
		rsp        proto.OrderInfoDetailResponse
	)

	result := global.DB.First(&order, req.Id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected < 1 {
		return nil, status.Errorf(codes.NotFound, "订单记录不存在")
	}

	if req.UserId > 0 && order.User != req.UserId {
		return nil, status.Errorf(codes.InvalidArgument, "非法的请求")
	}

	rsp.OrderInfo = &proto.OrderInfoResponse{
		Id:      order.ID,
		UserId:  order.User,
		OrderSn: order.OrderSn,
		PayType: order.PayType,
		Status:  order.Status,
		Post:    order.Post,
		Total:   order.OrderMount,
		Address: order.Address,
		Name:    order.SignerName,
		Mobile:  order.SingerMobile,
	}

	result = global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, goods := range orderGoods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			Id:         goods.ID,
			OrderId:    goods.Order,
			GoodsId:    goods.Goods,
			GoodsName:  goods.GoodsName,
			GoodsImage: goods.GoodsImage,
			GoodsPrice: goods.GoodsPrice,
			Nums:       goods.Nums,
		})
	}

	return &rsp, nil
}

func (*OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	return nil, nil
}
