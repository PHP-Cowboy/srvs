package main

import (
	"srvs/goods_srv/global"
	"srvs/goods_srv/initialize"
	"srvs/goods_srv/model"
)

func main() {
	initialize.InitConfig()
	initialize.InitMysql()

	db := global.DB
	_ = db.Set(model.TableOptions, model.GetOptions("商品分类表")).AutoMigrate(&model.Category{})
	_ = db.Set(model.TableOptions, model.GetOptions("品牌")).AutoMigrate(&model.Brands{})
	_ = db.Set(model.TableOptions, model.GetOptions("商品品牌分类中间表")).AutoMigrate(&model.GoodsCategoryBrand{})
	_ = db.Set(model.TableOptions, model.GetOptions("banner")).AutoMigrate(&model.Banner{})
	_ = db.Set(model.TableOptions, model.GetOptions("商品表")).AutoMigrate(&model.Goods{})

}
