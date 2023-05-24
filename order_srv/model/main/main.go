package main

import (
	"srvs/order_srv/global"
	"srvs/order_srv/initialize"
	"srvs/order_srv/model"
)

func main() {
	initialize.InitConfig()
	initialize.InitMysql()

	db := global.DB

	_ = db.Set(model.TableOptions, model.GetOptions("购物车表")).AutoMigrate(&model.ShoppingCart{})
	_ = db.Set(model.TableOptions, model.GetOptions("订单信息表")).AutoMigrate(&model.OrderInfo{})
	_ = db.Set(model.TableOptions, model.GetOptions("订单商品信息表")).AutoMigrate(&model.OrderGoods{})

}
