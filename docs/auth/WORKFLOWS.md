# GoTax AUTH Module - Workflows, Dataflows & Processes

---

## 1. Authentication Dataflow

```
User                    Frontend                  Auth Service              Database
  |                         |                          |                       |
  |--Login Form---------->|                          |                       |
  |                         |--POST /auth/login------>|                       |
  |                         |                          |--Query user by-------|
  |                         |                          |  username            |
  |                         |                          |<--user record--------|
  |                         |                          |                       |
  |                         |                          |--bcrypt.Compare------|
  |                         |                          |  (password, hash)    |
  |                         |                          |                       |
  |                         |                          |--Check: active?------|
  |                         |                          |  locked? locked_until|
  |                         |                          |  2FA enabled?        |
  |                         |                          |                       |
  |                         |                          |--Generate JWT--------|
  |                         |                          |  (access_token)      |
  |                         |                          |--Store refresh-------|
  |                         |                          |  token (SHA-256)     |
  |                         |                          |--Write audit log-----|
  |                         |                          |                       |
  |                         |<--{tokens, user_info}----|                       |
  |<--Dashboard------------|                          |                       |
```

## 2. VNeID OIDC Flow

```
User            Frontend         Auth Service        VNeID OIDC          National ID DB
 |                 |                  |                  |                    |
 |--Click VNeID-->|                  |                  |                    |
 |                 |--Redirect to--->|                  |                    |
 |                 |  VNeID login    |--Auth request--->|                    |
 |                 |                  |                  |                    |
 |<--VNeID login page----------------|                  |                    |
 |--Auth on VNeID--------------------|                  |                    |
 |                 |                  |                  |--Verify citizen----|
 |                 |                  |                  |  (CCCD, biometric) |
 |                 |                  |                  |<--match OK---------|
 |                 |                  |<--Auth code------|                    |
 |                 |<--Redirect with--|                  |                    |
 |                 |  auth code       |                  |                    |
 |                 |--POST callback-->|                  |                    |
 |                 |  {code, state}   |--Exchange code-->|                    |
 |                 |                  |<--ID token-------|                    |
 |                 |                  |  + access token  |                    |
 |                 |                  |                  |                    |
 |                 |                  |--Validate sig----|                    |
 |                 |                  |  (JWKS endpoint) |                    |
 |                 |                  |--Match/creat-----|                    |
 |                 |                  |  local user      |                    |
 |                 |                  |--Issue JWT-------|                    |
 |                 |                  |--Write audit Log-|                    |
 |                 |<--{tokens}-------|                  |                    |
 |<--Dashboard-----|                  |                  |                    |
```

## 3. RBAC Evaluation Flow

```
Request: User X tries to perform Action Y on Resource Z

1. Extract JWT from Authorization header
2. Validate JWT signature + expiry
3. Extract user_id, tenant_id, roles from claims
4. Load role permissions from cache (Redis) or DB
5. Evaluate:
   a. Does user have role with resource Z + action Y?
   b. Is scope "all" or "own" or "branch" (match user context)?
   c. Is tenant_id match?
6. Decision:
   [ALLOW] -> Forward to handler
   [DENY]  -> Return 403 AUTH_004 with resource info
7. Log audit event (access_granted or access_denied)
```

## 4. Digital Signature Process

```
User                    GoTax System              PKCS#11/USB           Remote CA
 |                         |                          |                    |
 |--Sign document-------->|                          |                    |
 |                         |--Hash document (SHA-256) |                    |
 |                         |                          |                    |
 |--"Insert token + PIN"->|                          |                    |
 |                         |--PKCS#11 C_Sign--------->|                    |
 |                         |<--Signature + cert-------|                    |
 |                         |                          |                    |
 |                         |--Validate via CRL/OCSP-->|                    |
 |                         |<--Cert valid-------------|                    |
 |                         |                          |                    |
 |                         |--Attach sig to doc-------|                    |
 |                         |--Verify sig integrity----|                    |
 |                         |--Store sig+doc+metadata -|                    |
 |                         |--Write audit log---------|                    |
 |<--"Signed successfully"-|                          |                    |
```

## 5. RBAC Configuration Process

```
Start: Admin login
  |
  v
[1] Navigate to RBAC Settings
  |
  v
[2] View Roles -> System displays predefined + custom roles
  |
  +--[2a] Create Custom Role
  |     |-- Enter role name, description
  |     |-- Select permissions: module + action + scope
  |     |-- Save -> Log audit event
  |
  +--[2b] Edit Existing Role
  |     |-- Modify permission matrix
  |     |-- Save -> Propagate changes to all users (optional immediate)
  |     |-- Log audit event
  |
  +--[2c] Delete Role
        |-- Check no users assigned
        |-- Confirm deletion
        |-- Delete -> Log audit event
  |
  v
[3] Assign Roles to Users
  |-- Select user(s) or group
  |-- Select role(s)
  |-- Confirm -> Inherited permissions active immediately
  |-- Log audit event
  |
  v
End
```

## 6. User Lifecycle Process

```
Registration
  |
  +--[Local Registration]
  |     |-- Admin creates user
  |     |-- User receives temp password
  |     |-- First login -> force password change
  |     |-- Optionally enable 2FA
  |     v
  |
  +--[VNeID Registration]
  |     |-- User logs in via VNeID first time
  |     |-- System creates local account (linked to VNeID)
  |     |-- Admin assigns roles
  |     v
  |
Active State
  |-- Normal operations
  |-- Password expiry (90d)
  |-- Session management
  |
  +--[Suspension/Deactivation]
        |-- Admin deactivates user
        |-- All active sessions terminated
        |-- Refresh tokens revoked
        |-- User cannot login
        |-- [Reactivation] Admin re-activates -> normal flow
        v
Termination/Deletion
  |-- Data retained per legal requirements
  |-- Account flagged "deleted" (soft delete)
```

## 7. Session Management Rules

| Rule | Value |
|------|-------|
| Access token lifetime | 15 minutes |
| Refresh token lifetime | 7 days |
| Max concurrent sessions | Configurable (default: 5 per user) |
| Idle timeout | 30 minutes |
| Absolute session timeout | 24 hours (re-login required) |
| Refresh token rotation | Each refresh issues new token, revokes old |
| Token storage | httpOnly, Secure, SameSite cookies |
| Session revocation | Admin can revoke all sessions for a user |

## 8. Password Policy

| Rule | Value | Reference |
|------|-------|-----------|
| Min length | 8 characters | Decree 13/2023 |
| Complexity | Upper + lower + digit + special | OWASP |
| Max age | 90 days | Decree 13/2023 |
| History | Last 5 passwords cannot repeat | Best practice |
| Max failed attempts | 5 within 15 min | OWASP |
| Lockout duration | 30 minutes auto, manual unlock by admin | Best practice |

## 9. 2FA Policy

| Rule | Value |
|------|-------|
| Required for | All admin roles (SystemAdmin, SecurityAdmin) |
| Optional for | Standard users |
| Methods | TOTP (preferred), SMS OTP (fallback) |
| Recovery codes | 8 single-use codes issued on setup |
| Remember device | Optional 30-day trust, requires re-auth for sensitive ops |

## 10. Data Retention & Cleanup

| Data Type | Retention | Deletion |
|-----------|-----------|----------|
| Audit logs | 5 years (Law on Accounting) | Archived to cold storage then purged |
| User accounts | Indefinite (soft delete) | Hard delete after 5 years inactive |
| Sessions | Until expiry or revocation | Immediate on logout/revoke |
| Password hash history | As long as user exists | On account deletion |
| OTP codes | 5 minutes | After verification or expiry |
| Failed login attempts | 30 days | Automatic cleanup |