package entity

type User struct {
	ID       int64 `gorm:"primaryKey"`
	Email    string
	Phone    string
	Password string
	Name     string
	Role     int8
}

func (User) TableName() string {
	return "users"
}
