package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Mobile    string `gorm:"varchar(11);not null;unique"`
	Password  string `gorm:"size:255;not null"`
	InviteCode string      `json:"InviteCode" gorm:"not null"`
}