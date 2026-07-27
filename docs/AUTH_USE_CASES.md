# AUTH Module — Use Cases

**Version:** 1.0  
**Author:** BA Lead + Chief Accountant  

---

## UC-1: Login (Username + Password)

**Actor:** Any user (admin, chief_accountant, accountant, viewer)  
**Precondition:** User account exists, active, not locked

### Happy Path
1. User POSTs `/api/v1/auth/login` with `{username, password}`
2. System validates credentials format
3. System retrieves user by username
4. System verifies `IsActive == true`
5. System verifies account not locked
6. System compares password hash (bcrypt cost=12)
7. System generates access token (RS256, 15min) + refresh token (opaque, 7d)
8. System logs login success to audit trail
9. System returns `{access_token, refresh_token, expires_in, user}`
10. User stores tokens securely

### Alternative Path — First Login / Password Expired
1-6 same as happy path
7. System detects password age > 90 days → returns `password_expired: true`
8. System forces password change before granting access
9. User must call UC-4 before accessing any resource

### Alternative Path — 2FA Required
1-6 same as happy path
7. System detects 2FA enabled for user
8. System returns `requires_2fa: true, temp_token: <one-time-use-token>`
9. User calls `/api/v1/auth/2fa/verify` with `{temp_token, totp_code}`
10. System verifies TOTP → returns `{access_token, refresh_token}`

### Exception Path — Invalid Credentials
1-6 same, but step 5 fails or step 6 fails
7. System increments failed attempt counter
8. If attempts ≥ 5, lock account for 30min, notify user
9. System logs login failure to audit trail
10. System returns 401 `{error: "invalid credentials"}` (generic, no enumeration)

### Exception Path — Account Locked
1-3 same
4. System checks lock status → account locked
5. If lock expired, auto-release and continue happy path
6. If lock active, return 423 `{error: "account locked", retry_after: <minutes>}`

### Exception Path — Account Inactive
1-3 same
4. `IsActive == false` → return 401 `{error: "account inactive"}`

### Exception Path — Rate Limited
1. User sends login request
2. Rate limiter detects >5 attempts/15min for IP+username
3. Return 429 `{error: "too many attempts", retry_after: <seconds>}`

---

## UC-2: Token Refresh

**Actor:** Any authenticated user  
**Precondition:** User holds valid refresh token (not expired, not revoked)

### Happy Path
1. User POSTs `/api/v1/auth/refresh` with `{refresh_token}`
2. System validates refresh token (opaque, stored in DB/Redis)
3. System verifies token not revoked/expired
4. System rotates refresh token (old token revoked, new issued)
5. System returns `{access_token, refresh_token, expires_in}`

### Exception Path — Expired Token
2. System detects token expired → 401 `{error: "refresh token expired"}`
3. User must re-login (UC-1)

### Exception Path — Revoked Token
2. System detects token revoked → 401 `{error: "session expired"}`
3. User must re-login (UC-1)

### Exception Path — Replay Attack
4. Token was already rotated (reused) → revoke ALL tokens for user
5. Force re-login on all devices
6. Log security event to audit trail

---

## UC-3: Logout

**Actor:** Any authenticated user  
**Precondition:** User holds valid access token

### Happy Path
1. User POSTs `/api/v1/auth/logout` with `Authorization: Bearer <token>`
2. System blacklists access token
3. System revokes refresh token
4. System logs logout to audit trail
5. Return 200 `{message: "logged out"}`

---

## UC-4: Password Change

**Actor:** Any authenticated user  
**Precondition:** User knows current password

### Happy Path
1. User PUTs `/api/v1/auth/password` with `{current_password, new_password}`
2. System verifies current password
3. System validates new password against policy (min 12, complexity)
4. System validates new password not in history (10)
5. System hashes new password (bcrypt cost=12)
6. System revokes ALL refresh tokens → force re-login
7. System logs password change to audit trail
8. Return 200 `{message: "password changed"}`

### Exception Path — Wrong Current Password
2. System detects mismatch → 401 `{error: "current password incorrect"}`

### Exception Path — Weak Password
3. Policy validation fails → 422 `{error: "password does not meet requirements", details: [...]}`

### Exception Path — Recently Used
4. Password in history → 422 `{error: "password was used recently"}`

---

## UC-5: Admin Force Password Reset

**Actor:** admin  
**Precondition:** Admin authenticated

### Happy Path
1. Admin POSTs `/api/v1/users/:id/force-reset`
2. System generates temporary password
3. System forces password change on next login (`password_expired = true`)
4. System logs to audit trail
5. Return `{message: "password reset", temporary_password: <one-time>}`
6. Admin communicates temporary password out-of-band

---

## UC-6: Enable/Disable 2FA

**Actor:** Any user  
**Precondition:** User authenticated, not already enabled (for enable)

### Happy Path — Enable
1. User POSTs `/api/v1/auth/2fa/enable` → returns `{secret, qr_code_url, backup_codes}`
2. User scans QR code into authenticator app
3. User confirms by calling `/api/v1/auth/2fa/confirm` with `{totp_code}`
4. System verifies TOTP, enables 2FA
5. System logs to audit trail
6. Return `{message: "2FA enabled", backup_codes: [10 codes]}`

### Happy Path — Disable
1. User POSTs `/api/v1/auth/2fa/disable` with `{current_password, totp_code}`
2. System verifies password + TOTP
3. System disables 2FA
4. System logs to audit trail
5. Return `{message: "2FA disabled"}`

---

## UC-7: Session Management — View Active Sessions

**Actor:** Any user  
**Precondition:** User authenticated

### Happy Path
1. User GETs `/api/v1/auth/sessions`
2. System returns list of active sessions:
   ```json
   [{ "id": "sess-1", "device": "Chrome/Windows", "ip": "192.168.1.1",
      "last_activity": "2026-07-27T10:00:00Z", "created_at": "2026-07-27T08:00:00Z",
      "is_current": true }]
   ```

---

## UC-8: Session Management — Force Logout

**Actor:** Any user  
**Precondition:** User authenticated

### Happy Path — Single Session
1. User DELETEs `/api/v1/auth/sessions/:session_id`
2. System revokes session's refresh token
3. System logs to audit trail
4. Return 200

### Happy Path — All Sessions
1. User DELETEs `/api/v1/auth/sessions`
2. System revokes ALL refresh tokens for user
3. System logs to audit trail
4. Return 200

---

## UC-9: Admin Audit Log — Auth Events

**Actor:** admin  
**Precondition:** Admin authenticated

### Happy Path
1. Admin GETs `/api/v1/audit?entity_type=auth&limit=50`
2. System returns auth audit entries (login, logout, password change, 2FA, role change)
3. Admin can filter by action, date range, username

---

## UC-10: User Creation (Admin)

**Actor:** admin  
**Precondition:** Admin authenticated

### Happy Path
1. Admin POSTs `/api/v1/users` with `{username, password, full_name, email, role}`
2. System validates username uniqueness
3. System validates password policy
4. System validates role is valid
5. System hashes password (bcrypt cost=12)
6. System creates user with `is_active: true`, `password_changed_at: null` → force change on first login
7. System logs to audit trail
8. Return 201 `{message: "user created", data: user}`
