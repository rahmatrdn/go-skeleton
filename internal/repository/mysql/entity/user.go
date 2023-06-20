package entity

type User struct {
	ID       int64  `gorm:"primaryKey"`
	Email    string `gorm:"index:email_k"`
	Phone    string `gorm:"index:phone_k"`
	Password string
	Name     string
	Role     int8 `gorm:"index:role_k"`
}

func (User) TableName() string {
	return "users"
}
