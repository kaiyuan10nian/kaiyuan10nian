package model

type Article struct {
	ID        uint `gorm:"primary_key"`
	UserId     uint      `json:"user_id" gorm:"not null"`
	Title      string `json:"title"  gorm:"type:text;not null"`
	Content    string `json:"content" gorm:"type:text;not null"`
	CreatedAt  Time   `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt  Time   `json:"updated_at" gorm:"type:timestamp"`
}
