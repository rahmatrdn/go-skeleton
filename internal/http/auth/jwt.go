package auth

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rahmatrdn/go-skeleton/config"
	"github.com/rahmatrdn/go-skeleton/entity"
)

const (
	privateKeyPath = "private_key.pem"
	publicKeyPath  = "public_key.pem"
)

type JWT struct{}

func NewJWTAuth() *JWT {
	return &JWT{}
}

type JWTAuth interface {
	GenerateToken(user entity.LoginResponse) (string, error)
}

func (j *JWT) GenerateToken(user entity.LoginResponse) (string, error) {
	cfg := config.NewConfig()

	privateKeyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return "", err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", err
	}

	claims := &entity.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JwtExpireDaysCount) * 24 * time.Hour)),
		},
		Email:      user.Email,
		UserID:     user.UserID,
		RoleAccess: user.RoleAccess,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fmt.Errorf("EMPTY TOKEN")
	}

	token := authHeader[7:]

	publicKeyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return err
	}

	claims := &entity.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil || !tkn.Valid {
		return err
	}

	// Set data in Local Context
	c.Locals("user_id", claims.UserID)

	return nil
}
