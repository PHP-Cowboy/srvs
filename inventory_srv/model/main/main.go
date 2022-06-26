package main

import (
	"shop-srvs/inventory_srv/global"
	"shop-srvs/inventory_srv/initialize"
	"shop-srvs/inventory_srv/model"
)

func main() {
	initialize.InitConfig()
	initialize.InitMysql()

	db := global.DB

	_ = db.Set(model.TableOptions, model.GetOptions("商品库存表")).AutoMigrate(&model.Inventory{})

}
