# GoTax AUTH Module - Use Cases

---

## UC-01: User Login (Password)

**Actor:** User (Accountant, Manager, etc.)
**Precondition:** User has valid account, password known

### Happy Path
```
1. User navigates to login page
2. System displays login form (username, password, tenant)
3. User enters credentials and submits
4. System validates credentials (bcrypt compare)
5. System checks account status (active, not locked)
6. System checks if 2FA required for user
   [If 2FA OFF] -> Issue JWT tokens, redirect to dashboard
   [If 2FA ON]  -> Issue temp token, redirect to 2FA page
7. System logs audit event (login success)
8. User accesses application
```

### Alternative Path A1: First Login - Force Password Change
```
3a. User enters temporary password
4a. System validates temporary password
5a. System detects "must_change_password" flag
6a. System redirects to password change form
7a. User sets new password (min 8 chars, complexity req)
8a. System hashes new password, clears flag
9a. Continue with normal login flow (step 6)
```

### Alternative Path A2: VNeID Login
```
3b. User clicks "Login with VNeID"
4b. System redirects to VNeID OIDC provider
5b. User authenticates on VNeID app (password + biometric)
6b. VNeID returns auth code to callback URL
7b. System exchanges code for VNeID tokens
8b. System validates token signature against JWKS endpoint
9b. System matches VNeID subject to local user account
   [If matched] -> Issue JWT, redirect to dashboard
   [If not matched] -> Initiate account linking flow
10b. Log audit event
```

### Exception Path E1: Invalid Credentials
```
4c. System validates credentials -> FAIL
5c. System increments failed_attempts counter
6c. If failed_attempts >= 5 within 15 min:
    -> Lock account for 30 min
    -> Send notification email to user
    -> Log audit event (account_locked)
7c. Return error AUTH_001 "Invalid credentials" with remaining attempts
```

### Exception Path E2: Account Locked
```
4d. System checks account status -> is_locked = true
5d. System checks locked_until timestamp
   [If locked_until > now] -> Return AUTH_005 "Account locked until HH:MM"
   [If locked_until < now] -> Unlock account, proceed with login
6d. Admin must unlock if locked by manual action
```

### Exception Path E3: Password Expired
```
4e. System checks password age > 90 days
5e. System redirects to password change form
6e. User must set new password before proceeding
7e. Cannot reuse last 5 passwords (history stored)
```

---

## UC-02: User Management (Admin)

**Actor:** SystemAdmin or SecurityAdmin
**Precondition:** Admin authenticated with appropriate role

### Happy Path - Create User
```
1. Admin navigates to User Management
2. Clicks "Add User"
3. Fills form: username, email, phone, full_name, role(s)
4. System generates temporary password
5. System sends activation email with temp password
6. User appears in list with status "Pending Activation"
7. Log audit event (user_created)
```

### Happy Path - Assign Role
```
1. Admin selects user
2. Clicks "Assign Role"
3. Selects role(s) from available list
4. System validates role permissions
5. System assigns role to user
6. Log audit event (role_assigned)
```

### Alternative Path - Group Assignment
```
1. Admin selects user group
2. Clicks "Assign Role to Group"
3. Selects role
4. All users in group inherit role
5. Log audit event (group_role_assigned)
```

### Exception Path - Duplicate User
```
3a. Admin enters username that exists
4a. System returns error: "Username already exists"
5a. Admin must use different username
```

---

## UC-03: Digital Signature - Sign Document

**Actor:** Authorized User
**Precondition:** User has registered digital certificate (USB token or remote)

### Happy Path - USB Token
```
1. User opens document requiring signature (e-invoice, tax declaration)
2. System prompts: "Please insert USB token and enter PIN"
3. User inserts USB token, enters PIN
4. System reads certificate from token via PKCS#11
5. System validates certificate (validity, CRL/OCSP check)
6. System creates digital signature on document hash
7. System attaches signature to document
8. System verifies signature integrity
9. System stores signed document + signature metadata
10. Log audit event (document_signed)
11. User receives confirmation
```

### Happy Path - Remote Signing
```
1. User opens document requiring signature
2. User selects remote signing method
3. User enters phone/email for OTP
4. System sends signing request to remote CA (e.g., MISA eSign, VNPT)
5. User receives push notification on mobile
6. User authenticates on mobile (biometric/PIN)
7. Remote CA creates signature, returns signed hash
8. System attaches signature to document
9. Log audit event (document_signed_remote)
```

### Alternative Path - Multi-Step Approval
```
1. Document requires 2 signatures: preparer + approver
2. Preparer signs first (step 1-9 above)
3. System marks document as "Partially Signed"
4. Approver receives notification
5. Approver reviews and signs
6. System marks document as "Fully Signed"
7. Log both signature events
```

### Exception Path - Certificate Expired
```
4a. System checks certificate validity -> EXPIRED
5a. System returns error: "Digital certificate expired on YYYY-MM-DD"
6a. User must renew certificate with CA provider
7a. Log audit event (signature_failed_cert_expired)
```

### Exception Path - Certificate Revoked
```
4b. System checks CRL/OCSP -> CERT_REVOKED
5b. System returns error: "Digital certificate has been revoked"
6b. User must obtain new certificate
7b. Log audit event (signature_failed_cert_revoked)
```

---

## UC-04: Multi-Tenant Access Control

**Actor:** DataAdmin (parent company)
**Precondition:** Admin authenticated, has cross-org permissions

### Happy Path
```
1. Admin switches tenant context from parent company
2. System checks cross-org delegation permission
3. System loads target company data
4. Admin operates within scope of delegated permissions
5. All actions logged with both admin identity and target tenant
```

### Exception Path - No Cross-Tenant Permission
```
2a. System checks cross-org delegation -> DENIED
3a. Return error AUTH_004 "Insufficient permissions for this organization"
```

---

## UC-05: Password Reset

**Actor:** Any user (unauthenticated)

### Happy Path
```
1. User clicks "Forgot Password"
2. User enters email or phone number
3. System validates account exists [no disclosure if not found]
4. System sends OTP to registered email/SMS
5. User enters OTP
6. System verifies OTP (valid for 5 min, single-use)
7. User enters new password
8. System hashes and stores new password
9. System revokes all existing refresh tokens
10. Log audit event (password_reset)
```

### Exception Path - Rate Limited
```
3a. User exceeded 3 reset requests in 1 hour
4a. Return rate limit error
5a. Log audit event (password_reset_rate_limited)
```

---

## UC-06: Audit Log Review

**Actor:** Auditor, SystemAdmin

### Happy Path
```
1. User navigates to Audit Log page
2. Applies filters: date range, event type, user, action
3. System queries partitioned audit_log table
4. Displays paginated results with: timestamp, user, action, IP, metadata
5. User can view detail of individual event
6. User can export filtered results to CSV/PDF
```

### Exception Path - Data Retention Expired
```
2a. User queries beyond 5-year retention period
3a. System returns: "Data beyond retention period has been purged"
```

---

## UC-07: Session Timeout

**Actor:** Any authenticated user

### Happy Path
```
1. User is idle for 30 minutes
2. Frontend detects inactivity timer
3. Attempts silent token refresh -> refresh token expired
4. System redirects to login page
5. User must re-authenticate
```

### Alternative Path - Token Refresh Before Expiry
```
1. Frontend detects access_token within 5 min of expiry
2. Frontend calls POST /auth/refresh with refresh_token
3. System validates refresh_token (hash match, not revoked, not expired)
4. System issues new access_token + rotated refresh_token
5. Previous refresh_token is revoked
6. User continues session uninterrupted
```

### Exception Path - Refresh Token Revoked
```
3a. System checks refresh_token -> revoked
4a. Return AUTH_003 "Session revoked, please login again"
5a. Log audit event (session_revoked)
```