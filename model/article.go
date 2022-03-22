package model

import (
	uuid "github.com/satori/go.uuid"
)

type Article struct {
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	UserId     uint      `json:"user_id" gorm:"not null"`
	Title      string `json:"title"  gorm:"type:text;not null"`
	Content    string `json:"content" gorm:"type:text;not null"`
	CreatedAt  Time   `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt  Time   `json:"updated_at" gorm:"type:timestamp"`
}
