package models

import "time"

type Base struct {
	Id         int       `json:"id" Gorm:"primary_key;AUTO_INCREMENT"`
	CreateTime time.Time `json:"create_time" Gorm:"auto_now_add;type(datetime)"`
	UpdateTime time.Time `json:"update_time" Gorm:"auto_now;type(datetime)"`
}
