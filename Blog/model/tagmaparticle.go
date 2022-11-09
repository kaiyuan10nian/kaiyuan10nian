package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type TagMapArticle struct {
	gorm.Model
	ArticleID         uuid.UUID `json:"article_id" gorm:"not null"`
	TagID  uint      `json:"tag_id" gorm:"not null"`
}

