# JWT Signing Toggle Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `JWT_SIGNING_METHOD` env toggle to switch JWT signing between RSA (RS512) and HMAC (HS512) while preserving full backward compatibility.

**Architecture:** Conditional `switch` on `cfg.JwtSigningMethod` inside each of `GenerateToken`, `VerifyToken`, and `RefreshToken`. No new files, no new interfaces. Two config fields replace hardcoded PEM file path constants; two new fields add HMAC support.

**Tech Stack:** Go, `github.com/golang-jwt/jwt/v4`, `github.com/gofiber/fiber/v3`, `github.com/joeshaw/envdecode`

---

## File Map

| Action | File | What changes |
|---|---|---|
| Modify | `config/config.go` | Add 4 JWT fields to `Config` struct |
| Modify | `.env` | Add `JWT_SIGNING_METHOD`, `JWT_PRIVATE_KEY_PATH`, `JWT_PUBLIC_KEY_PATH`, `JWT_SECRET_KEY` |
| Modify | `internal/http/auth/jwt.go` | Replace hardcoded constants; add RSA/HMAC switch to all 3 functions |
| Create | `internal/http/auth/jwt_test.go` | Unit tests for all branches of all 3 functions |

---

## Task 1: Config Fields + .env

**Files:**
- Modify: `config/config.go`
- Modify: `.env`

- [ ] **Step 1: Add 4 fields to `Config` struct in `config/config.go`**

Current struct ends at `JwtExpireDaysCount`. Add after it:

```go
JwtExpireDaysCount int    `env:"JWT_EXPIRE_DAYS_COUNT"`

// New fields — add these:
JwtSigningMethod  string `env:"JWT_SIGNING_METHOD,default=rsa"`
JwtPrivateKeyPath string `env:"JWT_PRIVATE_KEY_PATH,default=private_key.pem"`
JwtPublicKeyPath  string `env:"JWT_PUBLIC_KEY_PATH,default=public_key.pem"`
JwtSecretKey      string `env:"JWT_SECRET_KEY"`
```

- [ ] **Step 2: Add vars to `.env`**

Append to the `# JWT` section:

```env
# JWT
JWT_EXPIRE_DAYS_COUNT=3
JWT_SIGNING_METHOD=rsa
JWT_PRIVATE_KEY_PATH=private_key.pem
JWT_PUBLIC_KEY_PATH=public_key.pem
# JWT_SECRET_KEY=your-secret-key-min-32-chars   # uncomment for hmac mode
```

- [ ] **Step 3: Verify compilation**

```bash
cd /path/to/go-skeleton
go build ./...
```

Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add config/config.go .env.example
git commit -m "feat(config): add JWT signing method toggle env vars"
```

---

## Task 2: Write Tests for `GenerateToken` (TDD)

**Files:**
- Create: `internal/http/auth/jwt_test.go`

- [ ] **Step 1: Create test file with helpers and GenerateToken tests**

Create `internal/http/auth/jwt_test.go`:

```go
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
```

- [ ] **Step 2: Run tests — expect FAIL (GenerateToken not yet updated)**

```bash
cd /path/to/go-skeleton
go test ./internal/http/auth/... -run TestGenerateToken -v
```

Expected: `FAIL` — tests compile but RSA/HMAC switch not implemented yet, so RSA test passes (old code), HMAC tests fail with wrong error or wrong behavior. Confirm failures before proceeding.

- [ ] **Step 3: Commit test file**

```bash
git add internal/http/auth/jwt_test.go
git commit -m "test(auth): add failing tests for JWT GenerateToken RSA/HMAC toggle"
```

---

## Task 3: Implement `GenerateToken` — Make Tests Pass

**Files:**
- Modify: `internal/http/auth/jwt.go`

- [ ] **Step 1: Replace `GenerateToken` in `internal/http/auth/jwt.go`**

Remove the two `const` at the top of the file:
```go
// DELETE these two lines:
const (
    privateKeyPath = "private_key.pem"
    publicKeyPath  = "public_key.pem"
)
```

Replace the entire `GenerateToken` method with:

```go
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
```

- [ ] **Step 2: Run GenerateToken tests — expect PASS**

```bash
go test ./internal/http/auth/... -run TestGenerateToken -v
```

Expected:
```
--- PASS: TestGenerateToken_RSA
--- PASS: TestGenerateToken_HMAC
--- PASS: TestGenerateToken_HMAC_EmptySecret
--- PASS: TestGenerateToken_UnknownMethod
```

- [ ] **Step 3: Commit**

```bash
git add internal/http/auth/jwt.go
git commit -m "feat(auth): implement JWT signing method toggle in GenerateToken (RSA/HMAC)"
```

---

## Task 4: Write Tests for `VerifyToken` + `RefreshToken`

**Files:**
- Modify: `internal/http/auth/jwt_test.go` (append to existing file)

- [ ] **Step 1: Append helpers + VerifyToken + RefreshToken tests to `jwt_test.go`**

```go
// callVerifyToken builds a Fiber test request with the given Authorization header,
// executes VerifyToken inside a handler, and returns the captured user_id and error.
func callVerifyToken(t *testing.T, authHeader string) (int64, error) {
	t.Helper()
	app := fiber.New()
	var capturedUID int64
	var capturedErr error

	app.Get("/test", func(c fiber.Ctx) error {
		capturedErr = auth.VerifyToken(c)
		if capturedErr == nil {
			if uid, ok := c.Locals("user_id").(int64); ok {
				capturedUID = uid
			}
		}
		return nil
	})

	req := httptest.NewRequest("GET", "/test", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	_, err := app.Test(req)
	require.NoError(t, err)
	return capturedUID, capturedErr
}

// callRefreshToken builds a Fiber test request and returns the new token string and error.
func callRefreshToken(t *testing.T, authHeader string) (string, error) {
	t.Helper()
	app := fiber.New()
	var capturedToken string
	var capturedErr error

	app.Get("/test", func(c fiber.Ctx) error {
		capturedToken, capturedErr = auth.RefreshToken(c)
		return nil
	})

	req := httptest.NewRequest("GET", "/test", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	_, err := app.Test(req)
	require.NoError(t, err)
	return capturedToken, capturedErr
}

// --- VerifyToken ---

func TestVerifyToken_RSA(t *testing.T) {
	setBaseEnv(t)
	root := projectRoot()
	t.Setenv("JWT_SIGNING_METHOD", "rsa")
	t.Setenv("JWT_PRIVATE_KEY_PATH", filepath.Join(root, "private_key.pem"))
	t.Setenv("JWT_PUBLIC_KEY_PATH", filepath.Join(root, "public_key.pem"))

	j := auth.NewJWTAuth()
	token, err := j.GenerateToken(testUser())
	require.NoError(t, err)

	uid, err := callVerifyToken(t, "Bearer "+token)
	require.NoError(t, err)
	assert.Equal(t, int64(42), uid)
}

func TestVerifyToken_HMAC(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("JWT_SIGNING_METHOD", "hmac")
	t.Setenv("JWT_SECRET_KEY", "test-secret-key-at-least-32-chars-long")

	j := auth.NewJWTAuth()
	token, err := j.GenerateToken(testUser())
	require.NoError(t, err)

	uid, err := callVerifyToken(t, "Bearer "+token)
	require.NoError(t, err)
	assert.Equal(t, int64(42), uid)
}

func TestVerifyToken_EmptyHeader(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("JWT_SIGNING_METHOD", "rsa")

	_, err := callVerifyToken(t, "")
	assert.EqualError(t, err, "EMPTY TOKEN")
}

func TestVerifyToken_UnknownMethod(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("JWT_SIGNING_METHOD", "ecdsa")

	_, err := callVerifyToken(t, "Bearer sometoken")
	assert.EqualError(t, err, "unsupported JWT signing method: ecdsa")
}

// --- RefreshToken ---

func TestRefreshToken_RSA(t *testing.T) {
	setBaseEnv(t)
	root := projectRoot()
	t.Setenv("JWT_SIGNING_METHOD", "rsa")
	t.Setenv("JWT_PRIVATE_KEY_PATH", filepath.Join(root, "private_key.pem"))
	t.Setenv("JWT_PUBLIC_KEY_PATH", filepath.Join(root, "public_key.pem"))

	j := auth.NewJWTAuth()
	token, err := j.GenerateToken(testUser())
	require.NoError(t, err)

	newToken, err := callRefreshToken(t, "Bearer "+token)
	require.NoError(t, err)
	assert.NotEmpty(t, newToken)
}

func TestRefreshToken_HMAC(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("JWT_SIGNING_METHOD", "hmac")
	t.Setenv("JWT_SECRET_KEY", "test-secret-key-at-least-32-chars-long")

	j := auth.NewJWTAuth()
	token, err := j.GenerateToken(testUser())
	require.NoError(t, err)

	newToken, err := callRefreshToken(t, "Bearer "+token)
	require.NoError(t, err)
	assert.NotEmpty(t, newToken)
}

func TestRefreshToken_EmptyHeader(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("JWT_SIGNING_METHOD", "rsa")

	_, err := callRefreshToken(t, "")
	assert.EqualError(t, err, "EMPTY TOKEN")
}

func TestRefreshToken_UnknownMethod(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("JWT_SIGNING_METHOD", "ecdsa")

	_, err := callRefreshToken(t, "Bearer sometoken")
	assert.EqualError(t, err, "unsupported JWT signing method: ecdsa")
}
```

- [ ] **Step 2: Run new tests — expect FAIL**

```bash
go test ./internal/http/auth/... -run "TestVerifyToken|TestRefreshToken" -v
```

Expected: FAIL on HMAC tests (VerifyToken/RefreshToken still use old RSA-only code).

- [ ] **Step 3: Commit failing tests**

```bash
git add internal/http/auth/jwt_test.go
git commit -m "test(auth): add failing tests for VerifyToken and RefreshToken RSA/HMAC toggle"
```

---

## Task 5: Implement `VerifyToken` + `RefreshToken` — Make Tests Pass

**Files:**
- Modify: `internal/http/auth/jwt.go`

- [ ] **Step 1: Replace `VerifyToken` in `jwt.go`**

```go
func VerifyToken(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fmt.Errorf("EMPTY TOKEN")
	}

	token := authHeader[7:]
	cfg := config.NewConfig()

	claims := &entity.Claims{}

	var keyFunc jwt.Keyfunc
	switch cfg.JwtSigningMethod {
	case "rsa":
		publicKeyBytes, err := os.ReadFile(cfg.JwtPublicKeyPath)
		if err != nil {
			return err
		}
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
		if err != nil {
			return err
		}
		keyFunc = func(t *jwt.Token) (interface{}, error) { return publicKey, nil }
	case "hmac":
		if cfg.JwtSecretKey == "" {
			return fmt.Errorf("JWT_SECRET_KEY is required for hmac mode")
		}
		keyFunc = func(t *jwt.Token) (interface{}, error) { return []byte(cfg.JwtSecretKey), nil }
	default:
		return fmt.Errorf("unsupported JWT signing method: %s", cfg.JwtSigningMethod)
	}

	tkn, err := jwt.ParseWithClaims(token, claims, keyFunc)
	if err != nil || !tkn.Valid {
		return err
	}

	c.Locals("user_id", claims.UserID)
	return nil
}
```

- [ ] **Step 2: Replace `RefreshToken` in `jwt.go`**

```go
func RefreshToken(c fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("EMPTY TOKEN")
	}

	oldToken := authHeader[7:]
	cfg := config.NewConfig()

	claims := &entity.Claims{}

	var keyFunc jwt.Keyfunc
	switch cfg.JwtSigningMethod {
	case "rsa":
		publicKeyBytes, err := os.ReadFile(cfg.JwtPublicKeyPath)
		if err != nil {
			return "", err
		}
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
		if err != nil {
			return "", err
		}
		keyFunc = func(t *jwt.Token) (interface{}, error) { return publicKey, nil }
	case "hmac":
		if cfg.JwtSecretKey == "" {
			return "", fmt.Errorf("JWT_SECRET_KEY is required for hmac mode")
		}
		keyFunc = func(t *jwt.Token) (interface{}, error) { return []byte(cfg.JwtSecretKey), nil }
	default:
		return "", fmt.Errorf("unsupported JWT signing method: %s", cfg.JwtSigningMethod)
	}

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
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		return newToken.SignedString([]byte(cfg.JwtSecretKey))
	default:
		return "", fmt.Errorf("unsupported JWT signing method: %s", cfg.JwtSigningMethod)
	}
}
```

- [ ] **Step 3: Run all auth tests — expect all PASS**

```bash
go test ./internal/http/auth/... -v
```

Expected: all 12 tests pass.

- [ ] **Step 4: Run full build check**

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 5: Commit**

```bash
git add internal/http/auth/jwt.go
git commit -m "feat(auth): implement JWT signing method toggle in VerifyToken and RefreshToken (RSA/HMAC)"
```

---

## Self-Review Notes

- **Spec coverage:**
  - `JWT_SIGNING_METHOD` toggle → Task 1 + Tasks 3 + 5 ✓
  - `JWT_PRIVATE_KEY_PATH` / `JWT_PUBLIC_KEY_PATH` config → Task 1 ✓
  - `JWT_SECRET_KEY` config → Task 1 ✓
  - `GenerateToken` RSA/HMAC/unknown branches → Tasks 2 + 3 ✓
  - `VerifyToken` RSA/HMAC/unknown branches → Tasks 4 + 5 ✓
  - `RefreshToken` RSA/HMAC/unknown branches → Tasks 4 + 5 ✓
  - Error: unsupported method → all 3 functions ✓
  - Error: empty secret for hmac → GenerateToken + VerifyToken ✓
  - Backward compat (default=rsa) → envdecode `default=rsa` tag ✓
  - Remove hardcoded `privateKeyPath` / `publicKeyPath` constants → Task 3 ✓
