# GoTax AUTH Module - Legal & Standards Compliance Reference

---

## 1. Primary Legislation (Effective 2025-2026)

| Law/Decree | Effective | Key Requirements for AUTH |
|------------|-----------|--------------------------|
| Law on E-Transactions 20/2023/QH15 | 01/07/2024 | E-signature legal validity, digital signature standards, trust services |
| Decree 69/2024/ND-CP | 25/06/2024 | VNeID mandatory for organizations from 01/07/2025, eID levels 1 & 2 |
| Decree 23/2025/ND-CP | 10/04/2025 | Digital signature & trust services: PKI, CRL/OCSP, remote signing |
| Law on Tax Admin 108/2025/QH15 | 01/07/2026 | Electronic tax transactions, e-invoice signatures, data linkage |
| Decree 254/2026/ND-CP | 01/07/2026 | E-invoice rules, digital signature on invoices, tax authority codes |
| Circular 91/2026/TT-BTC | 01/07/2026 | Implementation guide for Decree 254 |
| Decision 29/2026/QD-TTg | 20/07/2026 | SSO via VNeID ONLY for National Public Service Portal |
| Decree 13/2023/ND-CP | 01/07/2023 | Personal data protection: consent, encryption, breach notification |
| Law on Network Security 86/2015/QH13 | 01/07/2016 | Network security requirements for information systems |
| Law on Cybersecurity 24/2018/QH14 | 01/01/2019 | Data localization, cybersecurity assessment |
| Circular 22/2020/TT-BTTTT | 01/01/2021 | Technical requirements for digital signature software |
| Circular 19/2021/TT-BTC | 18/03/2021 | E-transactions in tax (amended by 46/2024/TT-BTC) |

## 2. Secondary Standards

| Standard | Issuer | Relevance |
|----------|--------|-----------|
| VSA (Vietnam Standards on Auditing) | MOF/VACPA | Audit evidence integrity, retention |
| VSQC 1 (Quality Control) | MOF/VACPA | ISQC 1 equivalent for audit firms |
| VAS (Vietnam Accounting Standards) | MOF | Financial record authentication |
| ISO/IEC 27001:2022 | ISO | ISMS (FAST certified, BRAVO pursuing) |
| OWASP Top 10 (2021) | OWASP | Web app security baseline (BRAVO benchmark) |
| PKCS#11 | OASIS | Cryptoki for USB token digital signatures |
| RFC 7519 (JWT) | IETF | Token-based authentication |
| RFC 6238 (TOTP) | IETF | Time-based one-time passwords |
| RFC 2560 (OCSP) | IETF | Online Certificate Status Protocol |
| RFC 5280 (X.509) | IETF | Certificate path validation |

## 3. Outdated/Replaced Documents (DO NOT USE)

| Document | Replaced By | Status |
|----------|-------------|--------|
| Decree 130/2018/ND-CP (digital signatures) | Decree 23/2025/ND-CP | Superseded 10/04/2025 |
| Law on E-Transactions 51/2005/QH11 | Law 20/2023/QH15 | Superseded 01/07/2024 |
| Decree 59/2022/ND-CP (eID - old) | Decree 69/2024/ND-CP | Superseded 25/06/2024 |
| Law on Tax Admin 38/2019/QH14 | Law 108/2025/QH15 | Superseded 01/07/2026 |
| Decree 123/2020/ND-CP (e-invoice) | Decree 254/2026/ND-CP (amended) | Superseded 01/07/2026 |
| Decree 119/2018/ND-CP (e-invoice pilot) | Decree 123/2020/ND-CP | Superseded |

## 4. Key Compliance Deadlines

| Date | Requirement | Impact on AUTH |
|------|-------------|----------------|
| 01/07/2025 | VNeID mandatory for eTax login | AUTH must support VNeID OIDC |
| 01/01/2026 | Tax declaration for e-commerce HHKD | AUTH for platform tax data |
| 01/07/2026 | Law on Tax Admin 108/2025 effective | Full AUTH compliance required |
| 01/07/2026 | Decree 254/2026 effective | Digital sig on all e-invoices |
| 20/07/2026 | VNeID only for DVC Quoc Gia | VNeID primary auth method |
| Ongoing | Password change every 90 days | AUTH policy enforcement |
| Ongoing | Audit retention 5 years | Immutable audit log storage |

## 5. Data Classification for AUTH Module

| Data Type | Classification | Protection Required |
|-----------|---------------|-------------------|
| Password hash | Sensitive | bcrypt cost 12, never logged |
| TOTP secret | Sensitive | AES-256-GCM encrypted at rest |
| VNeID ID token | Sensitive | In-memory only, not persisted |
| Audit logs | Internal | Append-only, encrypted backup |
| User PII (name, email, phone) | Sensitive | Encrypted at rest, access controlled |
| Session tokens | Sensitive | httpOnly cookie, encrypted |
| Role/permission data | Internal | Standard access control |

## 6. Audit Log Events (Mandatory per VSA & Decree 13/2023)

```
Event Type                    Retention        Description
─────────────────────────────────────────────────────────────
LOGIN_SUCCESS                 5 years          User login successful
LOGIN_FAILED                  5 years          Failed login attempt
LOGOUT                        5 years          User logout
TOKEN_REFRESH                 1 year           Token rotation
PASSWORD_CHANGE               5 years          Password changed
PASSWORD_RESET                5 years          Password reset via forgot flow
2FA_ENABLED                   5 years          TOTP enabled
2FA_DISABLED                  5 years          TOTP disabled
USER_CREATED                  5 years          New user account
USER_DEACTIVATED              5 years          User deactivated
USER_REACTIVATED              5 years          User reactivated
ROLE_CREATED                  5 years          New role created
ROLE_MODIFIED                 5 years          Role permissions changed
ROLE_DELETED                  5 years          Role deleted
ROLE_ASSIGNED                 5 years          Role assigned to user
ROLE_REVOKED                  5 years          Role removed from user
DOCUMENT_SIGNED               5 years          Document digitally signed
DOCUMENT_SIGNED_REMOTE        5 years          Document signed via remote CA
SIGNATURE_FAILED              5 years          Signature attempt failed
CERTIFICATE_REGISTERED        5 years          Digital cert registered
CERTIFICATE_REMOVED           5 years          Digital cert removed
ACCESS_DENIED                 1 year           Permission denied attempt
SESSION_REVOKED               1 year           Admin revoked session
ACCOUNT_LOCKED                5 years          Account auto/manually locked
EXPORT_AUDIT_LOG              1 year           Audit log exported
```