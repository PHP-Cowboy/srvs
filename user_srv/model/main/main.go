package main

import (
	"shop-srvs/user_srv/global"
	"shop-srvs/user_srv/initialize"
	"shop-srvs/user_srv/model"
)

func main() {
	initialize.InitConfig()
	initialize.InitMysql()

	db := global.DB

	_ = db.Set(model.TableOptions, model.GetOptions("用户表")).AutoMigrate(&model.User{})

}
