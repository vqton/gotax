# AUTH Module — Workflows & Data Flows

---

## WF-1: Login Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Client  │     │  Router  │     │  Auth MW │     │ Handler  │     │ Service  │
└────┬─────┘     └────┬─────┘     └────┬─────┘     └────┬─────┘     └────┬─────┘
     │ POST /login    │                │                │                │
     │ ───────────────>                │                │                │
     │                │                │ (no auth MW   │                │
     │                │                │  on login)    │                │
     │                │ ──> ──────> ──>│                │                │
     │                │                │                │                │
     │                │                │ ──> ──────> ──>│                │
     │                │                │                │ Authenticate() │
     │                │                │                │───────────────>
     │                │                │                │                │
     │                │                │                │                │──> GetByUsername()
     │                │                │                │                │<── user
     │                │                │                │                │──> check IsActive
     │                │                │                │                │──> check locked
     │                │                │                │                │──> comparePassword()
     │                │                │                │                │──> generateAccessToken()
     │                │                │                │                │──> generateRefreshToken()
     │                │                │                │                │──> auditLogin()
     │                │                │                │                │
     │                │                │                │<── tokens+user│
     │                │                │<── 200        │                │
     │                │<── 200        │                │                │
     │<── 200        │                │                │                │
```

## WF-2: Token Refresh Flow

```
Client          Router          Handler         Service           DB/Redis
  │ POST /refresh │              │                │                  │
  │ ──────────────>              │                │                  │
  │               │ ──> ──> ──> │                │                  │
  │               │              │ RefreshToken() │                  │
  │               │              │───────────────>                   │
  │               │              │                │── get refresh    │
  │               │              │                │   token from DB  │
  │               │              │                │<── token         │
  │               │              │                │── verify expiry  │
  │               │              │                │── verify not     │
  │               │              │                │   revoked        │
  │               │              │                │── rotate token   │
  │               │              │                │── revoke old     │
  │               │              │                │── issue new      │
  │               │              │                │── store new      │
  │               │              │                │   token in DB    │
  │               │              │                │── audit refresh  │
  │               │              │<── tokens      │                  │
  │<── 200        │              │                │                  │
```

## WF-3: Auth Middleware (Request Guard)

```
Request
  │
  ▼
[All /api/v1/* routes]
  │
  ▼
AuthMiddleware()
  │
  ├── No Authorization header → 401
  │
  ├── Not Bearer format → 401
  │
  ├── Invalid/expired token → 401
  │
  ├── Token blacklisted → 401
  │
  └── Valid token
        │
        ▼
  Set context: user_id, username, role
        │
        ▼
  [Route handler]
        │
        ▼
  RoleMiddleware() [on admin routes]
        │
        ├── Missing role → 401
        ├── Wrong role → 403
        └── OK → next
```

## WF-4: Password Reset Flow

```
User            Client          Server          Email Service
  │ "Forgot PW"  │               │               │
  │ ─────────────>               │               │
  │               │ POST         │               │
  │               │ /forgot-pw   │               │
  │               │ {email}      │               │
  │               │──────────────>               │
  │               │               │── validate   │
  │               │               │   email      │
  │               │               │── generate   │
  │               │               │   reset token│
  │               │               │   (15min)    │
  │               │               │── store hash │
  │               │               │── send email │
  │               │               │──────────────>
  │               │ 200 OK        │               │
  │               │<──────────────               │
  │<── email      │               │               │
  │    with link  │               │               │
  │               │               │               │
  │ Click link    │               │               │
  │ ─────────────>               │               │
  │               │ GET /reset-pw│               │
  │               │ ?token=xxx   │               │
  │               │ + new PW body│               │
  │               │──────────────>               │
  │               │               │── verify token│
  │               │               │── hash new PW│
  │               │               │── check hist │
  │               │               │── update user│
  │               │               │── revoke all │
  │               │               │   sessions   │
  │               │               │── audit      │
  │               │               │── delete     │
  │               │               │   reset token│
  │               │ 200 OK        │               │
  │<── redirect   │<──────────────               │
  │    to login   │               │               │
```

## Data Flow: Auth Entities

```
┌─────────────────────────────────────────────────────────┐
│                     users table                          │
├─────────────────────────────────────────────────────────┤
│ id (UUID PK) │ username │ password_hash │ full_name     │
│ email │ role │ is_active │ failed_attempts │ locked_until │
│ password_changed_at │ password_history │ last_login_at │
│ totp_secret │ totp_enabled │ backup_codes                 │
│ created_at │ updated_at                                  │
└─────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────┐
│              refresh_tokens table               │
├────────────────────────────────────────────────┤
│ id │ user_id (FK) │ token_hash │ device_info   │
│ ip_address │ expires_at │ created_at │ revoked_at│
└────────────────────────────────────────────────┘

┌────────────────────────────────────────────────┐
│              password_reset_tokens table        │
├────────────────────────────────────────────────┤
│ id │ user_id (FK) │ token_hash │ expires_at    │
│ used_at │ created_at                            │
└────────────────────────────────────────────────┘

┌────────────────────────────────────────────────┐
│              audit_log table                    │
├────────────────────────────────────────────────┤
│ ... existing fields ...                        │
│ NEW: action can be: LOGIN, LOGIN_FAIL, LOGOUT  │
│ PASSWORD_CHANGE, PASSWORD_RESET, 2FA_ENABLE,   │
│ 2FA_DISABLE, ROLE_CHANGE, LOCKOUT, SESSION_KILL│
└────────────────────────────────────────────────┘
```

## WF-5: Account Lockout State Machine

```
                ┌─────────────────┐
                │    ACTIVE       │
                │ failed=0        │
                └────────┬────────┘
                         │ login failure
                         ▼
                ┌─────────────────┐
                │    WARNING      │
                │ failed=1-4      │
                └────────┬────────┘
                         │ 5th failure
                         ▼
                ┌─────────────────┐
                │    LOCKED       │  ← 30min timer
                │ locked_until=T  │
                └────────┬────────┘
                         │ timer expires
                         ▼
                ┌─────────────────┐
                │    ACTIVE       │
                │ failed=0        │
                └─────────────────┘
```
