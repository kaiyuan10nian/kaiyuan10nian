package model

import "github.com/jinzhu/gorm"

type InviteCode struct {
	gorm.Model
	Userid  uint      `json:"user_id" gorm:"not null"`
	Code    string      `json:"code" gorm:"not null"`
	Status  uint      `json:"status" gorm:"not null"`
}
