package entity

import "github.com/golang-jwt/jwt/v4"

type LoginReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type LoginResponse struct {
	UserID     int64  `json:"user_id" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required"`
	RoleAccess int8   `json:"role_access" validate:"required"`
	Token      string `json:"access_token" validate:"required"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserID     int64  `json:"user_id"`
	Email      string `json:"email"`
	RoleAccess int8   `json:"role"`
}
