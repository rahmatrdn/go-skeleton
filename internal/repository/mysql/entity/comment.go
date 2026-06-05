package entity

import "time"

type Comment struct {
	ID        int64     `gorm:"column:id"`
	PostID    int64     `gorm:"column:post_id"`
	UserID    *int64    `gorm:"column:user_id"`
	ParentID  *int64    `gorm:"column:parent_id"`
	Body      string    `gorm:"column:body"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Comment) TableName() string {
	return "comments"
}
