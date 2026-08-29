# JWT Signing Toggle Design

**Date:** 2026-06-05  
**Branch:** feat/upgrade-go-fiber-version  
**Approach:** Conditional logic within existing functions (Approach B)

---

## Overview

Add env-driven toggle to switch JWT signing between RSA (RS512) and HMAC (HS512). Backward compatible — default stays RSA. No new files, no new interfaces. All changes confined to `config/config.go` and `internal/http/auth/jwt.go`.

---

## Environment Variables

| Variable | Values | Default | Notes |
|---|---|---|---|
| `JWT_SIGNING_METHOD` | `rsa` / `hmac` | `rsa` | Toggle |
| `JWT_PRIVATE_KEY_PATH` | file path | `private_key.pem` | RSA mode only |
| `JWT_PUBLIC_KEY_PATH` | file path | `public_key.pem` | RSA mode only |
| `JWT_SECRET_KEY` | string (min 32 chars recommended) | — | HMAC mode only, required if method=hmac |

Existing `JWT_EXPIRE_DAYS_COUNT` unchanged.

---

## Config Changes (`config/config.go`)

Add to `Config` struct:

```go
JwtSigningMethod  string `env:"JWT_SIGNING_METHOD,default=rsa"`
JwtPrivateKeyPath string `env:"JWT_PRIVATE_KEY_PATH,default=private_key.pem"`
JwtPublicKeyPath  string `env:"JWT_PUBLIC_KEY_PATH,default=public_key.pem"`
JwtSecretKey      string `env:"JWT_SECRET_KEY"`
```

Remove hardcoded constants `privateKeyPath` and `publicKeyPath` from `jwt.go` — replaced by config fields.

---

## jwt.go Changes (`internal/http/auth/jwt.go`)

### GenerateToken

```
read cfg.JwtSigningMethod
if "rsa":
    load privateKey from cfg.JwtPrivateKeyPath
    jwt.NewWithClaims(RS512, claims).SignedString(privateKey)
else if "hmac":
    validate cfg.JwtSecretKey not empty
    jwt.NewWithClaims(HS512, claims).SignedString([]byte(secretKey))
else:
    return error "unsupported JWT signing method: <value>"
```

### VerifyToken

```
read cfg.JwtSigningMethod
build keyFunc:
    if "rsa"  → return publicKey (loaded from cfg.JwtPublicKeyPath)
    if "hmac" → return []byte(cfg.JwtSecretKey)
    else      → return error "unsupported JWT signing method"
jwt.ParseWithClaims(token, claims, keyFunc)
```

### RefreshToken

Same branch pattern as VerifyToken for parse, then same as GenerateToken for re-sign.

---

## Error Handling

| Condition | Error |
|---|---|
| `JWT_SIGNING_METHOD` not `rsa` or `hmac` | `"unsupported JWT signing method: <value>"` |
| `JWT_SECRET_KEY` empty when method=`hmac` | `"JWT_SECRET_KEY is required for hmac mode"` |
| File not found (RSA mode) | propagate `os.ReadFile` error as-is |

Validation happens inside each function (consistent with existing code pattern — no startup validation).

---

## Backward Compatibility

- `JWT_SIGNING_METHOD` unset → defaults to `rsa` → identical behavior to current code
- `JWT_PRIVATE_KEY_PATH` / `JWT_PUBLIC_KEY_PATH` default to existing filenames
- No changes to JWT claims structure (`entity.Claims`)
- No changes to middleware (`verify_token.go`)
- No changes to handler or usecase layers

---

## Out of Scope

- EdDSA or other algorithms
- Inline PEM via env var (file path only for RSA)
- Key rotation
- Startup validation of config
