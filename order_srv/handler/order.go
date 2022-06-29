package handler

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
	"shop-srvs/order_srv/global"
	"shop-srvs/order_srv/model"
	"shop-srvs/order_srv/proto/proto"
	"time"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

func GenderOrderSn(userId int32) string {
	now := time.Now()

	rand.Seed(now.UnixNano())

	return fmt.Sprintf("%d%d%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(), userId, rand.Intn(90)+10)
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
	var (
		carts        []*model.ShoppingCart
		cartNumsMap  = make(map[int32]int32)
		goodsId      []int32
		total        float32
		orderGoods   []*model.OrderGoods
		goodsInvInfo []*proto.GoodsInvInfo
	)

	db := global.DB
	result := db.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Find(&carts)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected <= 0 {
		return nil, status.Errorf(codes.NotFound, "用户购物车数据未找到")
	}

	for _, cart := range carts {
		goodsId = append(goodsId, cart.Goods)
		cartNumsMap[cart.Goods] = cart.Nums
	}

	//商品服务中查询商品信息
	goodsList, err := global.GoodsServer.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsId})
	if err != nil {
		return nil, err
	}

	for _, goods := range goodsList.Data {
		total += goods.ShopPrice * float32(cartNumsMap[goods.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      goods.Id,
			GoodsName:  goods.Name,
			GoodsImage: goods.GoodsFrontImage,
			GoodsPrice: goods.ShopPrice,
			Nums:       cartNumsMap[goods.Id],
		})

		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: goods.Id,
			Num:     cartNumsMap[goods.Id],
		})
	}

	//库存服务扣减库存
	_, err = global.InventoryServer.Sell(context.Background(), &proto.SellInfo{
		GoodsInvInfo: goodsInvInfo,
	})
	if err != nil {
		return nil, err
	}

	orderInfo := &model.OrderInfo{
		User:         req.UserId,
		OrderSn:      GenderOrderSn(req.UserId),
		OrderMount:   total,
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
	}

	tx := db.Begin()

	if result := tx.Save(&orderInfo); result.RowsAffected <= 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "订单数据保存失败")
	}

	for _, good := range orderGoods {
		good.Order = orderInfo.ID
	}

	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected <= 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "订单商品数据保存失败")
	}

	if result := tx.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(model.ShoppingCart{}); result.RowsAffected <= 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "购物车数据删除失败")
	}

	tx.Commit()

	return &proto.OrderInfoResponse{Id: orderInfo.ID, OrderSn: orderInfo.OrderSn, Total: total}, nil
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
	result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected < 1 {
		return nil, status.Errorf(codes.Internal, "订单状态更新失败")
	}

	return &emptypb.Empty{}, nil
}
