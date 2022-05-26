package model

import (
	"gorm.io/gorm"
	"time"
)

const TableOptions string = "gorm:table_options"

func GetOptions(tableName string) string {
	return "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci comment '" + tableName + "'"
}

type Base struct {
	Id       int32          `gorm:"primaryKey;type:int(11) unsigned AUTO_INCREMENT;comment:id"`
	CreateAt time.Time      `gorm:"type:datetime;not null;comment:创建时间"`
	UpdateAt time.Time      `gorm:"type:datetime;not null;comment:更新时间"`
	DeleteAt gorm.DeletedAt `gorm:"type:datetime;index;comment:删除时间"`
}
