package auth_test

import (
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/rahmatrdn/go-skeleton/internal/http/auth"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// projectRoot resolves the repo root from this test file's location.
// internal/http/auth/jwt_test.go → 4 levels up
func projectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "../../../..")
}

// setBaseEnv sets all required envdecode fields so config.NewConfig() does not panic.
func setBaseEnv(t *testing.T) {
	t.Helper()
	t.Setenv("MYSQL_POOL", "10")
	t.Setenv("MYSQL_SLOW_LOG_THRESHOLD", "300")
	t.Setenv("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/")
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	t.Setenv("MONGODB_DATABASE_NAME", "test")
	t.Setenv("REDIS_HOST", "127.0.0.1:6379")
	t.Setenv("REDIS_READ_TIMEOUT", "600")
	t.Setenv("REDIS_WRITE_TIMEOUT", "600")
	t.Setenv("JWT_EXPIRE_DAYS_COUNT", "3")
}

func testUser() *mentity.User {
	return &mentity.User{ID: 42, Email: "user@test.com", Role: 1}
}

// --- GenerateToken ---

func TestGenerateToken_RSA(t *testing.T) {
	setBaseEnv(t)
	root := projectRoot()
	t.Setenv("JWT_SIGNING_METHOD", "rsa")
	t.Setenv("JWT_PRIVATE_KEY_PATH", filepath.Join(root, "private_key.pem"))
	t.Setenv("JWT_PUBLIC_KEY_PATH", filepath.Join(root, "public_key.pem"))

	j := auth.NewJWTAuth()
	token, err := j.GenerateToken(testUser())

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateToken_HMAC(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("JWT_SIGNING_METHOD", "hmac")
	t.Setenv("JWT_SECRET_KEY", "test-secret-key-at-least-32-chars-long")

	j := auth.NewJWTAuth()
	token, err := j.GenerateToken(testUser())

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateToken_HMAC_EmptySecret(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("JWT_SIGNING_METHOD", "hmac")
	t.Setenv("JWT_SECRET_KEY", "")

	j := auth.NewJWTAuth()
	_, err := j.GenerateToken(testUser())

	assert.EqualError(t, err, "JWT_SECRET_KEY is required for hmac mode")
}

func TestGenerateToken_UnknownMethod(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("JWT_SIGNING_METHOD", "ecdsa")

	j := auth.NewJWTAuth()
	_, err := j.GenerateToken(testUser())

	assert.EqualError(t, err, "unsupported JWT signing method: ecdsa")
}

// Fiber + httptest helpers — used by VerifyToken and RefreshToken tests (added in Task 4)
var _ = fiber.New
var _ = httptest.NewRequest
