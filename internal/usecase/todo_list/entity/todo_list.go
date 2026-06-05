package entity

import "time"

type TodoListReq struct {
	ID          int64  `json:"id,omitempty" swaggerignore:"true"`
	UserID      int64  `json:"user_id,omitempty" swaggerignore:"true"`
	Title       string `json:"title,omitempty" validate:"required" name:"Judul"`
	Description string `json:"description" validate:"required" name:"Deskripsi"`
	DoingAt     string `json:"doing_at" validate:"required" name:"Tanggal Aktifitas"` // Format: YYYY-MM-DD
}

type TodoListResponse struct {
	ID          int64     `json:"id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DoingAt     time.Time `json:"doing_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r *TodoListReq) SetID(ID int64) {
	r.ID = ID
}

func (r *TodoListReq) SetUserID(UserID int64) {
	r.UserID = UserID
}
