package model

import "time"

type Article struct {
	ID       int       `gorm:"primaryKey" json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Category string    `json:"category"`
	Author   string    `json:"author"`
	Image    *string   `json:"image"`
	Created  time.Time `json:"created" gorm:"column:created_at"`
	Featured bool      `json:"featured"`
}

// TableName specifies the table name for GORM
func (Article) TableName() string {
	return "articles"
}

type Response struct {
	Message string      `json:"message"`
	Article *Article    `json:"article,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ImageUploadResponse struct {
	Message string `json:"message"`
	URL     string `json:"url"`
}
