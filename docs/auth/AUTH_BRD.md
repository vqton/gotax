# AUTH Module — Business Requirements Document

**Version:** 1.0  
**Status:** Draft  
**Date:** 2026-07-27  
**BA Lead:** 20+ yrs VN accounting software  
**Chief Accountant:** 20+ yrs VN corporate accounting  

---

## 1. Executive Summary

GoTax AUTH module provides authentication, authorization, session management, and security for the GoTax GL platform. Current implementation has fundamental gaps preventing production deployment (see AUTH_READINESS.md). This BRD defines the complete auth system required for Circular 99/2025/TT-BTC compliance and feature parity with MISA/Fast/Bravo ERP.

## 2. Business Objectives

| # | Objective | Success Metric |
|---|-----------|----------------|
| OBJ-1 | Secure user authentication | Zero auth-related incidents |
| OBJ-2 | Regulatory compliance | Pass Circular 99 Art 28 audit |
| OBJ-3 | Role-based access control | All API endpoints enforce authorization |
| OBJ-4 | Audit trail for all auth events | Every auth action logged, immutable |
| OBJ-5 | Self-service credential mgmt | Password reset without admin |
| OBJ-6 | Multi-factor authentication | Optional 2FA for high-risk roles |
| OBJ-7 | Session lifecycle management | Active sessions visible, revocable |

## 3. User Roles

| Role | Permissions | VN Equivalent |
|------|------------|---------------|
| `admin` | Full system access, user mgmt, audit view | Chủ DN / Giám đốc |
| `chief_accountant` | All accounting, approve/post, user view | Kế toán trưởng |
| `accountant` | Create/edit entries, view reports | Kế toán viên |
| `viewer` | Read-only reports, no transactions | Người xem báo cáo |
| `auditor` | Audit log, read-only financials (future) | Kiểm toán nội bộ |

## 4. Feature Requirements

### FR-1: Authentication
- Username + password login
- Rate-limited (5 attempts/15min per IP+username)
- Account lockout after 5 failures (30min auto-release)
- bcrypt cost=12
- Multi-session support (same user, multiple devices)

### FR-2: Token Management
- Access token: JWT, 15min expiry, RS256-signed
- Refresh token: opaque, 7d expiry, rotation on use
- Token blacklist (Redis) for immediate revocation
- No HS256 — migrate to RS256

### FR-3: Password Policy
| Rule | Value |
|------|-------|
| Min length | 12 chars |
| Upper case | ≥ 1 |
| Lower case | ≥ 1 |
| Digit | ≥ 1 |
| Special char | ≥ 1 |
| Expiry | 90 days |
| History | 10 passwords |
| Max attempts | 5 before lockout |

### FR-4: Multi-Factor Authentication (2FA)
- Optional per-user
- TOTP (RFC 6238) via authenticator app
- Backup codes (10 x single-use)
- Required for `admin` and `chief_accountant` roles

### FR-5: Password Reset
- Request: email → reset token (15min expiry, single-use)
- Self-service: no admin intervention
- Notify user of password change

### FR-6: Session Management
- View active sessions (device, IP, last activity)
- Force logout specific session
- Force logout all sessions (password change invalidates all)

### FR-7: Audit Trail
| Event | Data |
|-------|------|
| Login success | user, IP, timestamp, device |
| Login failure | username, IP, timestamp, attempt count |
| Logout | user, IP, timestamp |
| Password change | user, timestamp, by (user/admin) |
| 2FA enable/disable | user, timestamp |
| Role change | user, changed by, old role, new role |
| Lockout | user, IP, timestamp, duration |
| Token refresh | user, IP, timestamp |

### FR-8: API Security
- All endpoints behind auth middleware (except login + ping + swagger)
- CORS allowlist
- Security headers: CSP, HSTS, X-Frame-Options, X-Content-Type-Options
- Rate limiting on all API endpoints

## 5. Regulatory Compliance

| Regulation | Requirement | How AUTH Module Meets |
|------------|------------|----------------------|
| Circular 99/2025/TT-BTC Art 28 | Alert/prevent unauthorized interference | Rate limit + lockout + audit |
| Circular 99/2025/TT-BTC Art 28 | Confidentiality and security | bcrypt + RS256 + HTTPS |
| Circular 99/200 Art 28 | Change traceability | Audit trail for all auth events |
| Decree 23/2025/ND-CP | Digital signature trust | (v2: integrate CA) |
| Law on Info Security 2015 | Personal data protection | Password hash never stored plaintext |

## 6. Assumptions

1. Email service available for password reset
2. Redis available for token blacklist + rate limiter
3. TOTP via authenticator app (Google Authenticator, Authy)
4. HTTPS enforced by reverse proxy / load balancer
5. JWT_SECRET / RSA keys provided via environment or secrets manager

## 7. Out of Scope (v2)

- SSO / SAML / LDAP / OIDC
- WebAuthn / FIDO2 / biometric
- IP geo-fencing
- Device fingerprinting
- Adaptive MFA (risk-based)
- Session recording
