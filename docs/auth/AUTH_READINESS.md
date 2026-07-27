# AUTH Module — Production Readiness Assessment

**BA Lead:** 20+ yrs VN accounting software  
**Chief Accountant:** 20+ yrs VN corporate accounting  
**Review Date:** 2026-07-27  
**Regulatory Baseline:** Circular 99/2025/TT-BTC, Decree 130/2018/ND-CP, Decree 23/2025/ND-CP, Nghị định 85/2021/NĐ-CP (amended 2024)

---

## Verdict: ❌ NOT PRODUCTION-READY

| Layer | Score | Status |
|-------|-------|--------|
| Core Auth (login/JWT) | ⚠️ 50% | Missing rate limit, lockout, refresh |
| Password Security | ⚠️ 40% | bcrypt cost=10, no policy, no expiry |
| Session Mgmt | ❌ 0% | No revocation, no blacklist, no concurrent control |
| Access Control | ⚠️ 60% | RBAC exists, no row-level, no ABAC |
| Audit Trail | ⚠️ 30% | LOGIN audit action defined but never called |
| MFA/2FA | ❌ 0% | Not implemented |
| Password Reset | ❌ 0% | Not implemented |
| Regulatory Compliance | ⚠️ 35% | Missing Circular 99 Article 28 security requirements |
| Infrastructure | ❌ 20% | No HTTPS enforcement, no CORS, no CSP |
| Token Mgmt | ⚠️ 40% | HS256, 24h expiry, no refresh |
| **Overall** | **❌ 28%** | **NOT PROD-READY** |

---

## Critical Gaps (Block PROD)

### 1. No Rate Limiting on Login
- Attacker brute-force unlimited
- Required by Circular 99 Art 28: "alert or prevent intentional interference"
- Fix: rate limit 5 attempts/15min per IP + username

### 2. No Account Lockout
- No lock after N failed attempts
- Fix: lock 30min after 5 failures, notify user

### 3. No Password Policy
- No min length, no complexity, no expiry
- Circular 99 Art 28: "ensure confidentiality and security"
- Fix: min 12 chars, upper+lower+digit+special, 90-day expiry

### 4. bcrypt DefaultCost (10)
- Cost 10 marginal for 2026 (OWASP recommends 12+)
- Fix: cost = 12

### 5. No Refresh Token
- 24h JWT expiry, no renewal mechanism
- Users must re-login every 24h — UX fails
- Fix: short-lived access token (15min) + long-lived refresh token (7d) with rotation

### 6. No Token Revocation
- No blacklist → stolen token usable until expiry
- Fix: JWT blacklist with Redis, or short-lived tokens + refresh

### 7. No Audit Trail for Auth Events
- `AuditActionLogin` defined in code but never called
- No logging of: failed logins, password changes, role changes
- Fix: audit all auth events

### 8. No Password Reset Flow
- Users cannot recover accounts
- Fix: reset token via email, 15min expiry, one-time use

### 9. No 2FA/MFA
- Vietnamese accounting software (MISA, Fast, Bravo) all support SmartOTP/eToken
- Circular 99 Art 28: "security and alert mechanism"
- Fix: optional TOTP-based 2FA

### 10. HS256 (Symmetric)
- Single secret shared across all services
- No key rotation support
- Fix: RS256 with key pair, support key rotation

### 11. Global `jwtSecret` Variable
- Thread-safe? Read-only after init, but fragile
- No hot-reload for key rotation
- Fix: inject via config struct

### 12. No CORS Configuration
- No CORS middleware
- Mobile/web clients will fail cross-origin
- Fix: configure CORS allowlist

### 13. No Security Headers
- No CSP, HSTS, X-Frame-Options, X-Content-Type-Options
- Fix: add helmet-style middleware

---

## Gap Matrix vs Competitors

| Feature | MISA SME | Fast Acct | Bravo ERP | **GoTax (Current)** | **GoTax (Target)** |
|---------|----------|-----------|-----------|-------------------|-------------------|
| Login rate limit | ✅ | ✅ | ✅ | ❌ | ✅ |
| Account lockout | ✅ | ✅ | ✅ | ❌ | ✅ |
| Password policy | ✅ | ✅ | ✅ | ❌ | ✅ |
| bcrypt/scrypt | bcrypt | bcrypt | bcrypt | bcrypt cost=10 | bcrypt cost=12 |
| 2FA/SmartOTP | ✅ | ⚠️ (module) | ✅ | ❌ | ✅ |
| Token refresh | ✅ | ✅ | ✅ | ❌ | ✅ |
| Session mgmt | ✅ | ✅ | ✅ | ❌ | ✅ |
| Audit log | ✅ | ✅ | ✅ | partial | ✅ |
| Password reset | ✅ | ✅ | ✅ | ❌ | ✅ |
| RBAC | ✅ | ✅ | ✅ | basic | expanded |
| Row-level security | ✅ | ⚠️ | ✅ | ❌ | ❌ (v2) |
| IP restriction | ✅ | ✅ | ✅ | ❌ | ❌ (v2) |
| Device mgmt | ✅ | ❌ | ✅ | ❌ | ❌ (v2) |
| SSO/LDAP | ❌ | ❌ | ✅ | ❌ | ❌ (v2) |

---

## Compliance vs VN Regulations

| Regulation | Requirement | Status |
|------------|-------------|--------|
| Circular 99/2025/TT-BTC Art 28 | Alert/prevent unauthorized interference | ❌ |
| Circular 99/2025/TT-BTC Art 28 | Confidentiality and security | ⚠️ partial |
| Decree 130/2018/ND-CP Art 9 | Secure digital signature conditions | N/A (signature later) |
| Decree 23/2025/ND-CP | Electronic signature trust service | N/A (v2) |
| Nghị định 85/2021/NĐ-CP (amended) | E-commerce platform security | N/A (B2B only) |
| Law on Information Security 2015 | Personal data protection, breach notification | ⚠️ partial |
| Circular 78/2021/TT-BTC | E-tax transaction security | ❌ |

---

## Migration Priority

### Phase 1 — BLOCKER (must before PROD)
1. Rate limiting on login
2. Account lockout
3. Password policy
4. bcrypt cost → 12
5. Audit trail for auth events

### Phase 2 — HIGH (week 1-2)
6. Refresh token + token rotation
7. Token revocation/blacklist
8. Password reset flow
9. CORS + security headers

### Phase 3 — MEDIUM (week 3-4)
10. 2FA/TOTP
11. RS256 with key rotation
12. Session management (list active, force logout)

### Phase 4 — FUTURE (v2)
13. IP restriction / geo-fencing
14. Device management
15. SSO / LDAP / SAML
16. WebAuthn / FIDO2
