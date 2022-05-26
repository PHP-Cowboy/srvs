package model

import (
	"time"
)

type User struct {
	Base
	Mobile   string     `gorm:"type:varchar(11);comment:手机号;unique;not null"`
	Password string     `gorm:"type:varchar(100);not null;comment:密码"`
	Nickname string     `gorm:"type:varchar(20);comment:昵称"`
	Birthday *time.Time `gorm:"type:datetime;comment:生日"`
	Gender   int32      `gorm:"type:tinyint unsigned;comment:1:男,2:女,3:未知;default:3"`
	Role     int32      `gorm:"type:tinyint unsigned;comment:1:普通用户,2:管理员;default:2"`
}
