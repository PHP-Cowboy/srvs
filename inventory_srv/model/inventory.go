package model

type Inventory struct {
	BaseModel
	Goods      int32 `gorm:"index;comment:商品"`
	Stock      int32 `gorm:"comment:库存"`
	VersionNum int32 `gorm:"comment:商品"`
}
