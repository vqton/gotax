# GoTax AUTH Module - Technical Specifications

## Version: 1.0 | Stack: Go + Gin + PostgreSQL + Redis

---

## 1. Architecture

```
Client Layer (Web/Mobile/VNeID)
       | HTTPS/TLS 1.3
Auth Service (Go/Gin)
  |-- Auth Handler   (OIDC, Password, 2FA)
  |-- RBAC Engine    (Policy-based access)
  |-- Digital Sig    (PKCS#11, Remote API)
  |-- Session Mgr    (JWT + Redis)
  |-- Audit Logger   (Immutable append-log)
  |-- Tenant Mgr     (Multi-org isolation)
       |
Data Layer: PostgreSQL | Redis | Object Store
```

## 2. Auth Flow

### 2.1 Password Auth
```
POST /api/v1/auth/login
{ "username": "...", "password": "...", "tenant_id": "..." }
  -> Validate credentials (bcrypt compare)
  -> Check 2FA required? If yes, return 2FA challenge
  -> Issue access_token (15min) + refresh_token (7d)
  -> Log audit event
  <- { access_token, refresh_token, expires_in, user_info }
```

### 2.2 VNeID OIDC Auth
```
GET /api/v1/auth/vneid/login
  -> Redirect to VNeID OIDC provider
  -> User authenticates on VNeID
  -> VNeID returns auth code
POST /api/v1/auth/vneid/callback
  { "code": "...", "state": "..." }
  -> Exchange code for ID token + access token
  -> Validate token signature (JWKS)
  -> Map VNeID claims to local user
  -> Create/link user account
  -> Issue local JWT session
  -> Log audit event
  <- { access_token, refresh_token, user_info }
```

### 2.3 2FA (TOTP)
```
POST /api/v1/auth/2fa/setup
  -> Generate TOTP secret
  -> Return QR URI for Google Authenticator
POST /api/v1/auth/2fa/verify
  { "totp_code": "123456" }
  -> Validate TOTP (RFC 6238)
  -> Enable 2FA for user
POST /api/v1/auth/2fa/login
  { "temp_token": "...", "totp_code": "..." }
  -> Verify temp token + TOTP
  -> Issue full session tokens
```

## 3. RBAC Model

### Roles
```
SystemAdmin    - Full access, manage everything
SecurityAdmin  - Manage users, roles, permissions
DataAdmin      - Manage per-company data access
Accountant     - CRUD accounting data, no admin functions
Auditor        - Read-only, full audit log access
Manager        - Approve documents, view reports
Viewer         - Read-only basic data
CustomRole     - Flexible via permission matrix
```

### Permission Matrix (per module)
```
Module        | View | Add | Edit | Delete | Approve | Print | Export
--------------|------|-----|------|--------|---------|-------|--------
Tax Declar.   |  X   |  X  |  X   |   -    |    X    |   X   |   X
Invoices      |  X   |  X  |  X   |   -    |    X    |   X   |   X
Reports       |  X   |  -  |  -   |   -    |    -    |   X   |   X
Users         |  X   |  X  |  X   |   X    |    -    |   X   |   -
Audit Log     |  X   |  -  |  -   |   -    |    -    |   X   |   X
Settings      |  X   |  X  |  X   |   -    |    -    |   -   |   -
```

## 4. Database Schema (Core Tables)

### tenants
```sql
CREATE TABLE tenants (
  id          UUID PRIMARY KEY,
  name        VARCHAR(255) NOT NULL,
  tax_code    VARCHAR(20) UNIQUE,
  address     TEXT,
  is_active   BOOLEAN DEFAULT true,
  created_at  TIMESTAMPTZ DEFAULT now()
);
```

### users
```sql
CREATE TABLE users (
  id              UUID PRIMARY KEY,
  tenant_id       UUID REFERENCES tenants(id),
  username        VARCHAR(100) UNIQUE NOT NULL,
  password_hash   VARCHAR(255),
  email           VARCHAR(255),
  phone           VARCHAR(20),
  full_name       VARCHAR(255),
  vneid_id        VARCHAR(100) UNIQUE,  -- VNeID subject identifier
  is_active       BOOLEAN DEFAULT true,
  is_locked       BOOLEAN DEFAULT false,
  failed_attempts INT DEFAULT 0,
  locked_until    TIMESTAMPTZ,
  totp_secret     TEXT,  -- encrypted
  totp_enabled    BOOLEAN DEFAULT false,
  last_login      TIMESTAMPTZ,
  created_at      TIMESTAMPTZ DEFAULT now(),
  updated_at      TIMESTAMPTZ DEFAULT now()
);
```

### roles
```sql
CREATE TABLE roles (
  id          UUID PRIMARY KEY,
  tenant_id   UUID REFERENCES tenants(id),
  name        VARCHAR(100) NOT NULL,
  is_system   BOOLEAN DEFAULT false,  -- system predefined
  created_at  TIMESTAMPTZ DEFAULT now()
);
```

### permissions
```sql
CREATE TABLE permissions (
  id            UUID PRIMARY KEY,
  role_id       UUID REFERENCES roles(id) ON DELETE CASCADE,
  module        VARCHAR(100) NOT NULL,
  action        VARCHAR(50) NOT NULL,  -- view/add/edit/delete/approve/print/export
  scope         VARCHAR(50) DEFAULT 'all'  -- all/own/branch
);
```

### user_roles
```sql
CREATE TABLE user_roles (
  user_id       UUID REFERENCES users(id) ON DELETE CASCADE,
  role_id       UUID REFERENCES roles(id) ON DELETE CASCADE,
  assigned_by   UUID REFERENCES users(id),
  assigned_at   TIMESTAMPTZ DEFAULT now(),
  PRIMARY KEY (user_id, role_id)
);
```

### audit_log
```sql
CREATE TABLE audit_log (
  id          BIGSERIAL PRIMARY KEY,
  tenant_id   UUID REFERENCES tenants(id),
  user_id     UUID REFERENCES users(id),
  event_type  VARCHAR(50) NOT NULL,  -- login/logout/permission_change/sign_doc/failed_auth
  resource    VARCHAR(255),
  action      VARCHAR(100),
  ip_address  INET,
  user_agent  TEXT,
  metadata    JSONB,
  created_at  TIMESTAMPTZ DEFAULT now()
) PARTITION BY RANGE (created_at);
```

### refresh_tokens
```sql
CREATE TABLE refresh_tokens (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     UUID REFERENCES users(id) ON DELETE CASCADE,
  token_hash  VARCHAR(64) NOT NULL,
  expires_at  TIMESTAMPTZ NOT NULL,
  revoked     BOOLEAN DEFAULT false,
  created_at  TIMESTAMPTZ DEFAULT now()
);
```

## 5. API Endpoints

### Auth
```
POST   /api/v1/auth/login              - Password login
POST   /api/v1/auth/refresh            - Refresh token
POST   /api/v1/auth/logout             - Logout, revoke tokens
GET    /api/v1/auth/vneid/login        - VNeID OIDC redirect
POST   /api/v1/auth/vneid/callback     - VNeID callback
POST   /api/v1/auth/2fa/setup          - Enable TOTP
POST   /api/v1/auth/2fa/verify         - Verify TOTP setup
POST   /api/v1/auth/2fa/login          - 2FA step (post-password)
POST   /api/v1/auth/password/reset     - Request password reset
POST   /api/v1/auth/password/change    - Change password (authenticated)
```

### User Management (RBAC)
```
GET    /api/v1/users                   - List users (paginated)
POST   /api/v1/users                   - Create user
GET    /api/v1/users/:id               - Get user
PUT    /api/v1/users/:id               - Update user
DELETE /api/v1/users/:id               - Deactivate user
GET    /api/v1/users/:id/roles         - Get user roles
POST   /api/v1/users/:id/roles         - Assign roles
DELETE /api/v1/users/:id/roles/:role_id - Remove role
GET    /api/v1/roles                   - List roles
POST   /api/v1/roles                   - Create custom role
PUT    /api/v1/roles/:id               - Update role permissions
DELETE /api/v1/roles/:id               - Delete role
```

### Digital Signature
```
POST   /api/v1/sign/documents          - Sign document
POST   /api/v1/sign/verify             - Verify signature
GET    /api/v1/sign/certificates       - List registered certificates
POST   /api/v1/sign/certificates       - Register certificate
DELETE /api/v1/sign/certificates/:id   - Remove certificate
```

### Audit
```
GET    /api/v1/audit/log               - Query audit logs (paginated, filterable)
GET    /api/v1/audit/log/export        - Export audit logs
```

## 6. Security Controls

| Control | Implementation | Standard |
|---------|---------------|----------|
| Password hashing | bcrypt cost 12 | OWASP |
| Token signing | RSA-256 or Ed25519 | JWT RFC 7519 |
| Rate limiting | 10 req/min per IP on auth endpoints | OWASP |
| CORS | Whitelist origins | OWASP |
| CSP | Strict Content-Security-Policy | OWASP |
| HSTS | max-age=31536000 | OWASP |
| SQL injection | Parameterized queries (pgx) | OWASP |
| XSS | Output encoding, CSP | OWASP |
| CSRF | SameSite cookies + CSRF tokens | OWASP |
| Data encryption | AES-256-GCM for PII at rest | Decree 13/2023 |
| Session encryption | TLS 1.3 for transit | Decree 13/2023 |
| Audit immutability | Append-only, hash chain | VSA/VACPA |
| Password policy | Min 8 chars, complexity, 90d expiry | Decree 13/2023 |

## 7. Token Structure

### Access Token (JWT)
```json
{
  "sub": "user-uuid",
  "tenant_id": "tenant-uuid",
  "roles": ["accountant", "manager"],
  "permissions": ["tax_declaration:write", "invoice:read"],
  "iat": 1721800000,
  "exp": 1721800900,
  "jti": "unique-token-id"
}
```

### Refresh Token
- Opaque UUID stored in database (SHA-256 hashed)
- 7-day expiry
- Single-use (rotate on each refresh)
- Revocable by admin

## 8. Error Codes

| Code | HTTP | Description |
|------|------|-------------|
| AUTH_001 | 401 | Invalid credentials |
| AUTH_002 | 401 | Token expired |
| AUTH_003 | 401 | Token revoked |
| AUTH_004 | 403 | Insufficient permissions |
| AUTH_005 | 423 | Account locked |
| AUTH_006 | 401 | 2FA required |
| AUTH_007 | 401 | Invalid 2FA code |
| AUTH_008 | 429 | Rate limit exceeded |
| AUTH_009 | 400 | Invalid VNeID token |
| AUTH_010 | 500 | Internal auth error |