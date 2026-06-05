package entity

import "time"

type CommentReq struct {
	ID       int64  `json:"id,omitempty" swaggerignore:"true"`
	UserID   int64  `json:"user_id,omitempty" swaggerignore:"true"`
	PostID   int64  `json:"post_id" validate:"required" name:"Post ID"`   // ID of the post being commented on
	ParentID *int64 `json:"parent_id,omitempty" name:"Parent Comment ID"` // Optional, for replies to other comments
	Body     string `json:"body" validate:"required" name:"Komentar"`
}

type CommentResponse struct {
	ID        int64     `json:"id,omitempty"`
	PostID    int64     `json:"post_id"`
	UserID    *int64    `json:"user_id,omitempty"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *CommentReq) SetID(ID int64) {
	r.ID = ID
}

func (r *CommentReq) SetUserID(UserID int64) {
	r.UserID = UserID
}
