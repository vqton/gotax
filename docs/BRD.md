# BRD: GoTax Authentication & Authorization Module

## Version: 1.0 | Date: 2026-07-24 | Status: DRAFT

---

## 1. Executive Summary

### Project
GoTax - Vietnam Tax & Accounting Management System

### Business Need
Current codebase has zero authentication. Tax/accounting data is legally protected (Decree 13/2023/ND-CP on personal data protection, Law on Tax Administration 108/2025/QH15). All Vietnamese government e-services now mandate VNeID (Decree 69/2024/ND-CP, Decision 29/2026/QD-TTg). A compliant AUTH module is prerequisite for production.

### Strategic Alignment
- National digital transformation (Project 06, 2022-2025, vision 2030)
- MOF e-tax modernisation roadmap
- Compliance: Law on E-Transactions 2023, Decree 23/2025 (digital signatures), Decree 13/2023 (PDP)

---

## 2. Business Objectives

| # | Objective | KPI |
|---|-----------|-----|
| BO1 | Enable secure user authentication | Zero auth bypass in pentest |
| BO2 | Comply with VNeID mandatory auth (Decision 29/2026/QD-TTg) | VNeID integration by launch |
| BO3 | Support RBAC for multi-tenant enterprises | Role granularity = MISA/FAST level |
| BO4 | Enable digital signature for tax documents | Support Decree 23/2025 compliant signing |
| BO5 | Audit trail for all tax operations | Immutable audit log |
| BO6 | Multi-factor authentication (2FA) | 2FA optional but available |

---

## 3. Scope

### In Scope
- User registration (local + VNeID SSO)
- Authentication: password, VNeID (OIDC), 2FA (TOTP/SMS)
- Role-based access control (RBAC) with granular permissions
- Digital signature integration (USB token + remote signing)
- Session management with JWT + refresh tokens
- Audit logging (all auth events)
- Multi-tenancy (company/org hierarchy)
- User/group management with admin console
- Password policy enforcement (Decree 13/2023 Art. 10)

### Out of Scope (Phase 2)
- Biometric authentication (face/fingerprint) - planned Phase 2
- IFRS reporting module
- Advanced fraud detection ML

---

## 4. Key Stakeholders

| Role | Interest |
|------|----------|
| Chief Accountant | Audit trail, digital signature, tax compliance |
| IT Administrator | User management, RBAC, SSO |
| Tax Authority (GDT) | Integration via VNeID, eTax API |
| External Auditor | Audit log, non-repudiation |
| CFO | Data security, compliance, cost |

---

## 5. Legal & Regulatory Compliance Matrix

| Law/Decree | Requirement | AUTH Module Impact |
|------------|-------------|-------------------|
| Law on E-Transactions 20/2023/QH15 | Digital signature legal validity | Must support qualified digital signatures |
| Decree 69/2024/ND-CP | VNeID for organizations by 01/07/2025 | Mandatory VNeID OIDC integration |
| Decree 23/2025/ND-CP | Digital signature standards | Remote signing, X.509 cert validation |
| Law on Tax Administration 108/2025/QH15 | Electronic tax transactions | AUTH for eTax API calls |
| Decree 13/2023/ND-CP | Personal data protection | Consent mgmt, data encryption, breach notification |
| Decision 29/2026/QD-TTg | SSO via VNeID only from 20/07/2026 | VNeID as primary auth |
| Decree 254/2026/ND-CP | E-invoice digital signature | Signing workflow for invoices |
| Circular 22/2020/TT-BTTTT | Digital signature software standards | Validate CKS software compliance |
| ISO/IEC 27001:2022 | ISMS (FAST benchmark) | Security controls |
| OWASP Top 10 | Web app security (BRAVO benchmark) | AUTH secure by design |
| VSA (Vietnam Standards on Auditing) | Audit evidence integrity | Audit log immutability |

---

## 6. Functional Requirements

### FR1: Authentication
| ID | Requirement | Priority |
|----|-------------|----------|
| FR1.1 | Password-based login with bcrypt hashing | P0 |
| FR1.2 | VNeID SSO via OpenID Connect | P0 |
| FR1.3 | TOTP-based 2FA (Google Authenticator compatible) | P1 |
| FR1.4 | SMS OTP 2FA | P1 |
| FR1.5 | Session timeout & idle lock | P0 |
| FR1.6 | Account lockout after N failed attempts | P0 |
| FR1.7 | Password reset with email/SMS OTP | P0 |
| FR1.8 | Concurrent session control | P1 |

### FR2: Authorization (RBAC)
| ID | Requirement | Priority |
|----|-------------|----------|
| FR2.1 | System Admin role (full access) | P0 |
| FR2.2 | Security Admin role (user/role mgmt) | P0 |
| FR2.3 | Data Admin role (per-company data) | P1 |
| FR2.4 | Role templates (Accountant, Auditor, Manager, Viewer) | P0 |
| FR2.5 | Custom role creation | P1 |
| FR2.6 | Permission per module (read/write/delete/approve/print) | P0 |
| FR2.7 | Data-level access control (per-org/per-branch) | P1 |
| FR2.8 | Group-based permission assignment | P0 |

### FR3: Digital Signature
| ID | Requirement | Priority |
|----|-------------|----------|
| FR3.1 | USB token digital signature (PKCS#11) | P0 |
| FR3.2 | Remote signing (API-based, e.g. MISA eSign, VNPT) | P0 |
| FR3.3 | Signature verification & cert validation (CRL/OCSP) | P0 |
| FR3.4 | Multi-step signing workflow | P1 |
| FR3.5 | Signature audit trail | P0 |

### FR4: Audit
| ID | Requirement | Priority |
|----|-------------|----------|
| FR4.1 | Log all login/logout events | P0 |
| FR4.2 | Log all permission changes | P0 |
| FR4.3 | Log all digital signature events | P0 |
| FR4.4 | Immutable audit storage | P0 |
| FR4.5 | Audit log export (CSV/PDF) | P1 |

### FR5: Multi-Tenancy
| ID | Requirement | Priority |
|----|-------------|----------|
| FR5.1 | Company/Org entity hierarchy | P0 |
| FR5.2 | Data isolation per tenant | P0 |
| FR5.3 | Cross-org admin delegation | P1 |

---

## 7. Non-Functional Requirements

| # | Requirement | Target |
|---|-------------|--------|
| NFR1 | Auth response time | < 500ms p95 |
| NFR2 | Concurrent users | 10,000 |
| NFR3 | Token expiry | Access: 15min, Refresh: 7d |
| NFR4 | Password hash | bcrypt cost 12 |
| NFR5 | Session encryption | AES-256-GCM |
| NFR6 | Audit retention | 5 years (Law on Accounting Art. 13) |
| NFR7 | Availability | 99.9% uptime |

---

## 8. Risks & Mitigation

| Risk | Impact | Probability | Mitigation |
|------|--------|------------|------------|
| VNeID API changes | High | Medium | Abstract OIDC layer |
| Digital signature vendor lock-in | Medium | Medium | Multiple CA support |
| Data breach via AUTH bypass | Critical | Low | OWASP security review, pentest |
| Regulatory change (new decree) | Medium | Medium | Config-driven policy engine |

---

## 9. Budget Estimate (Phase 1)

| Item | Cost (VND) |
|------|-----------|
| Development (4 months, 4 engineers) | ~800M |
| Security audit & pentest | ~150M |
| VNeID integration certification | ~50M |
| Infrastructure (cloud) | ~100M/year |
| **Total** | **~1.1B** |