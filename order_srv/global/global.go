package global

import (
	"gorm.io/gorm"
	"srvs/order_srv/config"
	"srvs/order_srv/proto/proto"
)

var (
	DB              *gorm.DB
	ServerConfig    = config.ServerConfig{}
	NacosConfig     = config.NacosConfig{}
	GoodsServer     proto.GoodsServer
	InventoryServer proto.InventoryServer
)
