# AUTH Module — Technical Specifications

**Version:** 1.0  
**Status:** Draft  
**Regulatory Baseline:** Circular 99/2025/TT-BTC Art 28, Decree 23/2025/ND-CP  

---

## 1. API Endpoints

### Authentication

| Method | Path | Description | Auth | Rate Limit |
|--------|------|-------------|------|------------|
| POST | `/api/v1/auth/login` | Login | None | 5/15min per IP+user |
| POST | `/api/v1/auth/refresh` | Refresh tokens | None (use refresh) | 10/15min |
| POST | `/api/v1/auth/logout` | Logout | Bearer | — |
| PUT | `/api/v1/auth/password` | Change password | Bearer | — |
| POST | `/api/v1/auth/forgot-password` | Request reset | None | 1/5min per email |
| POST | `/api/v1/auth/reset-password` | Execute reset | None (use token) | 1/15min |
| POST | `/api/v1/auth/2fa/enable` | Start 2FA setup | Bearer | — |
| POST | `/api/v1/auth/2fa/confirm` | Confirm 2FA setup | Bearer | — |
| POST | `/api/v1/auth/2fa/disable` | Disable 2FA | Bearer | — |
| POST | `/api/v1/auth/2fa/verify` | Verify 2FA during login | None (use temp_token) | 5/15min |
| GET | `/api/v1/auth/sessions` | List active sessions | Bearer | — |
| DELETE | `/api/v1/auth/sessions/:id` | Force logout session | Bearer | — |
| DELETE | `/api/v1/auth/sessions` | Force logout all | Bearer | — |

### Users (Admin)

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | `/api/v1/users` | Create user | Bearer + admin |
| GET | `/api/v1/users` | List users | Bearer + admin |
| GET | `/api/v1/users/:id` | Get user | Bearer |
| PUT | `/api/v1/users/:id` | Update user | Bearer + admin |
| DELETE | `/api/v1/users/:id` | Delete user | Bearer + admin |
| POST | `/api/v1/users/:id/force-reset` | Force password reset | Bearer + admin |

## 2. Data Models

```go
type User struct {
    ID                string      `json:"id"`
    Username          string      `json:"username"`
    PasswordHash      string      `json:"-"`
    FullName          string      `json:"full_name"`
    Email             string      `json:"email,omitempty"`
    Role              UserRole    `json:"role"`
    IsActive          bool        `json:"is_active"`
    FailedAttempts    int         `json:"-"`            // 0-5
    LockedUntil       *time.Time  `json:"-"`
    PasswordChangedAt *time.Time  `json:"password_changed_at"`
    PasswordHistory   []string    `json:"-"`            // last 10 hashes
    TOTPSecret        string      `json:"-"`            // encrypted
    TOTPEnabled       bool        `json:"totp_enabled"`
    BackupCodes       []string    `json:"-"`            // hashed
    LastLoginAt       *time.Time  `json:"last_login_at"`
    LastLoginIP       string      `json:"-"`
    CreatedAt         time.Time   `json:"created_at"`
    UpdatedAt         time.Time   `json:"updated_at"`
}

type RefreshToken struct {
    ID         string
    UserID     string
    TokenHash  string
    DeviceInfo string
    IPAddress  string
    ExpiresAt  time.Time
    CreatedAt  time.Time
    RevokedAt  *time.Time
}

type PasswordResetToken struct {
    ID        string
    UserID    string
    TokenHash string
    ExpiresAt time.Time
    UsedAt    *time.Time
    CreatedAt time.Time
}
```

## 3. Service Interface

```go
type AuthService interface {
    // Authentication
    Login(ctx context.Context, username, password string) (*AuthResult, error)
    Login2FA(ctx context.Context, tempToken, totpCode string) (*AuthResult, error)
    RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
    Logout(ctx context.Context, userID, accessToken, refreshToken string) error

    // Password
    ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error
    ForgotPassword(ctx context.Context, email string) error
    ResetPassword(ctx context.Context, token, newPassword string) error
    ForceResetPassword(ctx context.Context, userID string) (string, error)

    // 2FA
    Enable2FA(ctx context.Context, userID string) (*TOTPSetup, error)
    Confirm2FA(ctx context.Context, userID, totpCode string) error
    Disable2FA(ctx context.Context, userID, password, totpCode string) error

    // Sessions
    GetActiveSessions(ctx context.Context, userID string) ([]Session, error)
    RevokeSession(ctx context.Context, userID, sessionID string) error
    RevokeAllSessions(ctx context.Context, userID string) error

    // Users
    CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
    UpdateUser(ctx context.Context, id string, req *UpdateUserRequest) error
    DeleteUser(ctx context.Context, id string) error
}
```

## 4. Token Format

### Access Token (RS256 JWT)

```json
{
  "sub": "U-001",
  "aud": "gotax-api",
  "iss": "gotax-auth",
  "username": "admin",
  "role": "admin",
  "iat": 1722057600,
  "exp": 1722058500,
  "jti": "unique-token-id"
}
```

### Refresh Token (Opaque)

- 32-byte random value (crypto/rand)
- Stored as SHA-256 hash in DB
- Never transmitted as JWT
- Rotated on each use (old revoked, new issued)

## 5. Password Policy Enforcement

```go
var PasswordPolicy = struct {
    MinLength      int
    RequireUpper   bool
    RequireLower   bool
    RequireDigit   bool
    RequireSpecial bool
    HistorySize    int
    MaxAgeDays     int
    MaxAttempts    int
    LockoutMinutes int
}{
    MinLength:      12,
    RequireUpper:   true,
    RequireLower:   true,
    RequireDigit:   true,
    RequireSpecial: true,
    HistorySize:    10,
    MaxAgeDays:     90,
    MaxAttempts:    5,
    LockoutMinutes: 30,
}
```

## 6. Audit Events

| Action | Entity Type | Description |
|--------|-------------|-------------|
| LOGIN | auth | Successful login |
| LOGIN_FAIL | auth | Failed login attempt |
| LOGOUT | auth | User logout |
| TOKEN_REFRESH | auth | Token refreshed |
| PASSWORD_CHANGE | auth | User changed own password |
| PASSWORD_RESET | auth | Password reset via forgot-password |
| PASSWORD_FORCE_RESET | user | Admin force reset |
| 2FA_ENABLE | auth | 2FA enabled |
| 2FA_DISABLE | auth | 2FA disabled |
| LOCKOUT | auth | Account locked after failures |
| SESSION_KILL | auth | Session force-logout |
| ROLE_CHANGE | user | Admin changed user role |

## 7. Error Codes

| Code | HTTP | Description |
|------|------|-------------|
| INVALID_CREDENTIALS | 401 | Wrong username or password |
| ACCOUNT_INACTIVE | 401 | User account disabled |
| ACCOUNT_LOCKED | 423 | Account temporarily locked |
| TOKEN_EXPIRED | 401 | Access/refresh token expired |
| TOKEN_REVOKED | 401 | Token was revoked |
| TOKEN_REUSED | 401 | Refresh token replayed |
| WEAK_PASSWORD | 422 | Password fails policy |
| PASSWORD_REUSED | 422 | Password in history |
| PASSWORD_EXPIRED | 401 | Password age > 90 days |
| 2FA_REQUIRED | 428 | 2FA verification needed |
| INVALID_TOTP | 401 | Invalid TOTP code |
| RATE_LIMITED | 429 | Too many requests |
| SESSION_NOT_FOUND | 404 | Session ID not found |

## 8. RS256 Key Rotation

```go
type KeyPair struct {
    ID        string
    PublicKey  *rsa.PublicKey
    PrivateKey *rsa.PrivateKey
    CreatedAt time.Time
    ExpiresAt time.Time
}

type KeyStore interface {
    GetSigningKey(ctx context.Context) (*KeyPair, error)   // current key
    GetVerificationKey(ctx context.Context, kid string) (*KeyPair, error)
    RotateKeys(ctx context.Context) error                  // create new, expire old
}
```

- Keys rotated every 90 days
- Previous key kept for verification until all tokens expire (max 15min + 7d overlap)
- JWKS endpoint at `GET /api/v1/auth/.well-known/jwks.json`

## 9. Implementation Priority

| Phase | Items | Dependencies |
|-------|-------|-------------|
| P1: STOP-BLEED | rate limit, lockout, pw policy, bcrypt 12, audit trail | none |
| P2: CORE | refresh token, revocation, password reset, CORS, headers | P1 |
| P3: SECURITY | RS256 key rotation, 2FA/TOTP | P2 |
| P4: UX | session management, backup codes | P3 |
