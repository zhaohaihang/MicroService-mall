package model

import (
	"gorm.io/gorm"
	"time"
)

// User 用户表
type User struct {
	gorm.Model
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null comment '用户手机号'"`
	Password string     `gorm:"type:varchar(100);not null comment '用户密码'"`
	NickName string     `gorm:"type:varchar(20) comment '用户名'"`
	Birthday *time.Time `gorm:"type:datetime comment '用户出生日期'"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女,male表示男'"`
	Role     int        `gorm:"column:role;default:1;type:int comment '1表示普通用户,2表示管理员'"`
}
