package main

import (
	"log"
	"os"
	"shop-srvs/goods_srv/model"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	dsn := "root:8!yJRAJwhH6t2xaK@tcp(114.116.88.12)/goods?charset=utf8mb4&parseTime=True&loc=Local"

	logger := logger2.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger2.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			Colorful:      true,        //禁用彩色打印
			LogLevel:      logger2.Info,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tb_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true,  // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}

	_ = db.Set(model.TableOptions, model.GetOptions("商品分类表")).AutoMigrate(&model.Category{})
	_ = db.Set(model.TableOptions, model.GetOptions("品牌")).AutoMigrate(&model.Brands{})
	_ = db.Set(model.TableOptions, model.GetOptions("商品品牌分类中间表")).AutoMigrate(&model.GoodsCategoryBrand{})
	_ = db.Set(model.TableOptions, model.GetOptions("banner")).AutoMigrate(&model.Banner{})
	_ = db.Set(model.TableOptions, model.GetOptions("商品表")).AutoMigrate(&model.Goods{})

}
