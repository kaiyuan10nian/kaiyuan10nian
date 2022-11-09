package model

import "github.com/jinzhu/gorm"

type Tags struct {
	gorm.Model
	TagName  string      `json:"tagname" gorm:"not null"`
}
