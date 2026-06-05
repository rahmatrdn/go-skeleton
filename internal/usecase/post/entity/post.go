package entity

import "time"

// PostStatus constants
const (
	PostStatusDraft     = 1
	PostStatusPublished = 2
	PostStatusArchived  = 3
)

type PostReq struct {
	ID     int64  `json:"id,omitempty" swaggerignore:"true"`
	UserID int64  `json:"user_id,omitempty" swaggerignore:"true"`
	Title  string `json:"title" validate:"required" name:"Judul"`
	Body   string `json:"body" validate:"required" name:"Konten"`
	Status int    `json:"status" validate:"required,oneof=1 2 3" name:"Status"` // 1: Draft, 2: Published, 3: Archived
}

type PostResponse struct {
	ID          int64      `json:"id,omitempty"`
	UserID      int64      `json:"user_id"`
	UserName    string     `json:"user_name"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Body        string     `json:"body"`
	Status      int        `json:"status"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (r *PostReq) SetID(ID int64) {
	r.ID = ID
}

func (r *PostReq) SetUserID(UserID int64) {
	r.UserID = UserID
}
