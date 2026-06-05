package entity

import "time"

type Post struct {
	ID          int64      `gorm:"column:id"`
	UserID      int64      `gorm:"column:user_id"`
	UserName    string     `gorm:"column:user_name"`
	Title       string     `gorm:"column:title"`
	Slug        string     `gorm:"column:slug"`
	Body        string     `gorm:"column:body"`
	Status      int        `gorm:"column:status"`
	PublishedAt *time.Time `gorm:"column:published_at"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
}

func (Post) TableName() string {
	return "posts"
}
