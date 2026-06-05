package auth

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rahmatrdn/go-skeleton/config"
	"github.com/rahmatrdn/go-skeleton/entity"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
)

type JWT struct{}

func NewJWTAuth() *JWT {
	return &JWT{}
}

type JWTAuth interface {
	GenerateToken(user *mentity.User) (string, error)
	RefreshToken(oldToken string) (string, error)
}

func (j *JWT) GenerateToken(user *mentity.User) (string, error) {
	cfg := config.NewConfig()

	claims := &entity.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JwtExpireDaysCount) * 24 * time.Hour)),
		},
		Email:      user.Email,
		UserID:     user.ID,
		RoleAccess: user.Role,
	}

	switch cfg.JwtSigningMethod {
	case "rsa":
		privateKeyBytes, err := os.ReadFile(cfg.JwtPrivateKeyPath)
		if err != nil {
			return "", err
		}
		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
		if err != nil {
			return "", err
		}
		token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
		return token.SignedString(privateKey)
	case "hmac":
		if cfg.JwtSecretKey == "" {
			return "", fmt.Errorf("JWT_SECRET_KEY is required for hmac mode")
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		return token.SignedString([]byte(cfg.JwtSecretKey))
	default:
		return "", fmt.Errorf("unsupported JWT signing method: %s", cfg.JwtSigningMethod)
	}
}

func VerifyToken(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fmt.Errorf("EMPTY TOKEN")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	cfg := config.NewConfig()

	keyFunc, err := buildKeyFunc(cfg)
	if err != nil {
		return err
	}

	claims := &entity.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, keyFunc)
	if err != nil || !tkn.Valid {
		return err
	}

	// Set data in Local Context
	c.Locals("user_id", claims.UserID)

	return nil
}

func (j *JWT) RefreshToken(oldToken string) (string, error) {
	cfg := config.NewConfig()

	keyFunc, err := buildKeyFunc(cfg)
	if err != nil {
		return "", err
	}

	claims := &entity.Claims{}
	tkn, err := jwt.ParseWithClaims(oldToken, claims, keyFunc)
	if err != nil || !tkn.Valid {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JwtExpireDaysCount) * 24 * time.Hour))

	switch cfg.JwtSigningMethod {
	case "rsa":
		privateKeyBytes, err := os.ReadFile(cfg.JwtPrivateKeyPath)
		if err != nil {
			return "", err
		}
		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
		if err != nil {
			return "", err
		}
		newToken := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
		return newToken.SignedString(privateKey)
	case "hmac":
		if cfg.JwtSecretKey == "" {
			return "", fmt.Errorf("JWT_SECRET_KEY is required for hmac mode")
		}
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		return newToken.SignedString([]byte(cfg.JwtSecretKey))
	default:
		return "", fmt.Errorf("unsupported JWT signing method: %s", cfg.JwtSigningMethod)
	}
}

func RefreshToken(c fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("EMPTY TOKEN")
	}

	oldToken := strings.TrimPrefix(authHeader, "Bearer ")
	cfg := config.NewConfig()

	keyFunc, err := buildKeyFunc(cfg)
	if err != nil {
		return "", err
	}

	claims := &entity.Claims{}
	tkn, err := jwt.ParseWithClaims(oldToken, claims, keyFunc)
	if err != nil || !tkn.Valid {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// Update expiry
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JwtExpireDaysCount) * 24 * time.Hour))

	switch cfg.JwtSigningMethod {
	case "rsa":
		privateKeyBytes, err := os.ReadFile(cfg.JwtPrivateKeyPath)
		if err != nil {
			return "", err
		}
		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
		if err != nil {
			return "", err
		}
		newToken := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
		return newToken.SignedString(privateKey)
	case "hmac":
		if cfg.JwtSecretKey == "" {
			return "", fmt.Errorf("JWT_SECRET_KEY is required for hmac mode")
		}
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		return newToken.SignedString([]byte(cfg.JwtSecretKey))
	default:
		return "", fmt.Errorf("unsupported JWT signing method: %s", cfg.JwtSigningMethod)
	}
}

// buildKeyFunc returns a jwt.Keyfunc based on cfg.JwtSigningMethod.
func buildKeyFunc(cfg *config.Config) (jwt.Keyfunc, error) {
	switch cfg.JwtSigningMethod {
	case "rsa":
		publicKeyBytes, err := os.ReadFile(cfg.JwtPublicKeyPath)
		if err != nil {
			return nil, err
		}
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
		if err != nil {
			return nil, err
		}
		return func(t *jwt.Token) (interface{}, error) {
			return publicKey, nil
		}, nil
	case "hmac":
		if cfg.JwtSecretKey == "" {
			return nil, fmt.Errorf("JWT_SECRET_KEY is required for hmac mode")
		}
		return func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecretKey), nil
		}, nil
	default:
		return nil, fmt.Errorf("unsupported JWT signing method: %s", cfg.JwtSigningMethod)
	}
}
