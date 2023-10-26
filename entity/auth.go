package entity

import "github.com/golang-jwt/jwt/v4"

type UserRole int8

const (
	Admin UserRole = 1
	Guest UserRole = 2
)

func GetRoleName(role UserRole) string {
	switch role {
	case Admin:
		return "Admin"
	case Guest:
		return "Guest"
	default:
		return "Unknown"
	}
}

type LoginReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type LoginResponse struct {
	UserID     int64  `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	RoleAccess int8   `json:"role_access"`
	Token      string `json:"access_token"`
}

type CreateUserReq struct {
	Name            string `json:"name" validate:"required" name:"Nama"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ReenterPassword string `json:"reenter_password" validate:"required"`
	Phone           string `json:"phone" validate:"required" name:"Nomor Telepon"`
	RoleAccess      int8   `json:"role_access" validate:"required" name:"Hak Akses"`
}
type CreateUserResponse struct {
	UserID     int64  `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	RoleAccess string `json:"role_access"`
	Phone      string `json:"phone"`
	Token      string `json:"access_token"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserID     int64  `json:"user_id"`
	Email      string `json:"email"`
	RoleAccess int8   `json:"role"`
}
