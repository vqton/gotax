# AUTH Module — Templates

**Version:** 1.0  

---

## Login Request/Response

```json
// POST /api/v1/auth/login
// Request
{
  "username": "admin",
  "password": "P@ssword2026!"
}

// Response — Success (no 2FA)
// 200
{
  "access_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleS0wMDEifQ...",
  "refresh_token": "a1b2c3d4e5f6...32byte-encoded",
  "expires_in": 900,
  "token_type": "Bearer",
  "user": {
    "id": "U-001",
    "username": "admin",
    "full_name": "Admin User",
    "role": "admin",
    "totp_enabled": false
  }
}

// Response — 2FA Required
// 200
{
  "requires_2fa": true,
  "temp_token": "temp-abc123...single-use",
  "expires_in": 300,
  "user": {
    "id": "U-001",
    "username": "admin",
    "full_name": "Admin User",
    "role": "admin",
    "totp_enabled": true
  }
}

// Response — Password Expired
// 200
{
  "password_expired": true,
  "message": "Password must be changed",
  "temp_token": "temp-def456...single-use",
  "user": {
    "id": "U-001",
    "username": "admin",
    "full_name": "Admin User",
    "role": "admin"
  }
}

// Response — Invalid Credentials
// 401
{
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid username or password"
  }
}

// Response — Account Locked
// 423
{
  "error": {
    "code": "ACCOUNT_LOCKED",
    "message": "Account temporarily locked",
    "retry_after_minutes": 30
  }
}

// Response — Rate Limited
// 429
{
  "error": {
    "code": "RATE_LIMITED",
    "message": "Too many login attempts",
    "retry_after_seconds": 900
  }
}
```

## 2FA Verification

```json
// POST /api/v1/auth/2fa/verify
// Request
{
  "temp_token": "temp-abc123",
  "totp_code": "123456"
}

// Response — Success
// 200
{
  "access_token": "eyJ...",
  "refresh_token": "a1b2c3...",
  "expires_in": 900,
  "token_type": "Bearer"
}

// Response — Invalid Code
// 401
{
  "error": {
    "code": "INVALID_TOTP",
    "message": "Invalid verification code"
  }
}
```

## Token Refresh

```json
// POST /api/v1/auth/refresh
// Request
{
  "refresh_token": "a1b2c3d4e5f6..."
}

// Response — Success
// 200
{
  "access_token": "eyJ...",
  "refresh_token": "f6e5d4c3b2a1...",
  "expires_in": 900,
  "token_type": "Bearer"
}

// Response — Expired
// 401
{
  "error": {
    "code": "TOKEN_EXPIRED",
    "message": "Refresh token expired, please login again"
  }
}

// Response — Revoked / Replay
// 401
{
  "error": {
    "code": "TOKEN_REVOKED",
    "message": "Session expired, please login again"
  }
}
```

## Password Change

```json
// PUT /api/v1/auth/password
// Request
{
  "current_password": "P@ssword2026!",
  "new_password": "N3wP@ssword2026!"
}

// Response — Success
// 200
{
  "message": "Password changed successfully"
}

// Response — Wrong Current
// 401
{
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Current password is incorrect"
  }
}

// Response — Weak Password
// 422
{
  "error": {
    "code": "WEAK_PASSWORD",
    "message": "Password does not meet requirements",
    "details": [
      "Minimum 12 characters",
      "At least 1 uppercase letter",
      "At least 1 lowercase letter",
      "At least 1 digit",
      "At least 1 special character"
    ]
  }
}

// Response — Recently Used
// 422
{
  "error": {
    "code": "PASSWORD_REUSED",
    "message": "Password was used recently"
  }
}
```

## Forgot/Reset Password

```json
// POST /api/v1/auth/forgot-password
// Request
{
  "email": "user@company.com"
}

// Response — Always 200 (prevent enumeration)
// 200
{
  "message": "If email exists, reset link sent"
}

// POST /api/v1/auth/reset-password
// Request
{
  "token": "reset-abc123...",
  "new_password": "N3wP@ssword2026!"
}

// Response — Success
// 200
{
  "message": "Password reset successfully"
}
```

## 2FA Setup

```json
// POST /api/v1/auth/2fa/enable
// Response
// 200
{
  "secret": "JBSWY3DPEHPK3PXP",
  "qr_code_url": "otpauth://totp/GoTax:admin?secret=...&issuer=GoTax",
  "backup_codes": [
    "1234-5678",
    "2345-6789",
    // ... 10 codes
  ]
}

// POST /api/v1/auth/2fa/confirm
// Request
{
  "totp_code": "123456"
}

// Response — Success
// 200
{
  "message": "2FA enabled successfully"
}

// POST /api/v1/auth/2fa/disable
// Request
{
  "current_password": "P@ssword2026!",
  "totp_code": "123456"
}

// Response — Success
// 200
{
  "message": "2FA disabled successfully"
}
```

## Session Management

```json
// GET /api/v1/auth/sessions
// Response
// 200
[
  {
    "id": "sess-abc123",
    "device": "Mozilla/5.0 Chrome/120 Windows NT 10.0",
    "ip": "192.168.1.100",
    "created_at": "2026-07-27T08:00:00Z",
    "last_activity": "2026-07-27T10:30:00Z",
    "is_current": true
  },
  {
    "id": "sess-def456",
    "device": "GoTax Mobile App/2.0 Android 14",
    "ip": "10.0.0.5",
    "created_at": "2026-07-26T14:00:00Z",
    "last_activity": "2026-07-27T09:15:00Z",
    "is_current": false
  }
]

// DELETE /api/v1/auth/sessions/sess-def456
// Response
// 200
{
  "message": "Session revoked"
}

// DELETE /api/v1/auth/sessions
// Response
// 200
{
  "message": "All sessions revoked"
}
```

## User Management (Admin)

```json
// POST /api/v1/users
// Request
{
  "username": "newaccountant",
  "password": "P@ssword2026!",
  "full_name": "New Accountant",
  "email": "new@company.com",
  "role": "accountant"
}

// Response — Success
// 201
{
  "message": "user created",
  "data": {
    "id": "U-002",
    "username": "newaccountant",
    "full_name": "New Accountant",
    "email": "new@company.com",
    "role": "accountant",
    "is_active": true,
    "created_at": "2026-07-27T10:00:00Z"
  }
}

// POST /api/v1/users/U-002/force-reset
// Response
// 200
{
  "message": "Password reset successfully",
  "temporary_password": "Temp-abc-123!"
}
```

## Common Error Format

```json
// All error responses follow this shape:
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable description"
  }
}

// Validation errors include details:
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": [
      "username is required",
      "password minimum 12 characters"
    ]
  }
}
```

## Table: Auth Audit Log Migration

```sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS failed_attempts INT DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS locked_until TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_changed_at TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_history TEXT[] DEFAULT '{}';
ALTER TABLE users ADD COLUMN IF NOT EXISTS totp_secret TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS totp_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS backup_codes TEXT[] DEFAULT '{}';
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_ip VARCHAR(45);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(64) NOT NULL,
    device_info TEXT,
    ip_address  VARCHAR(45),
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at  TIMESTAMPTZ
);
CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_hash ON refresh_tokens(token_hash);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(64) NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_reset_tokens_hash ON password_reset_tokens(token_hash);
```
