# GoTax Company Module — Use Cases

## Version: 1.0 | Standard: Circular 99/2025/TT-BTC, Law on Enterprise 59/2025/QH15

---

## UC-01: Register New Company (Doanh nghiep)

**Actor:** Admin (Super admin or tenant admin)
**Precondition:** Authenticated with `company:create` permission, tenant context exists

### Happy Path
1. User opens Company Registration form
2. User enters: legal name (VN + EN), legal form (LLC/JSC/Sole Proprietorship/etc.), MSDN, MST
3. System validates MST format: `^[0-9]{10}(-[0-9]{3})?$` (10 digits or 13 digits for branches)
4. System validates MSDN format per provincial DPI requirements
5. User enters: registered address, head office, phone, email, legal representative, chief accountant
6. User selects accounting regime: TT99 (default), TT133, or TT58 (eligibility checked)
7. User configures fiscal year start month (default January)
8. System creates company with status = ACTIVE
9. System auto-generates fiscal periods for first year
10. System loads default COA based on selected regime
11. System logs audit entry (CREATE, COMPANY, MST)
12. System returns 201 with full company profile

### Alternative Path 1a — MST already exists:
- System checks uniqueness across tenant
- Returns 409 CONFLICT, error: "tax code already registered"

### Alternative Path 1b — Invalid MST format:
- System validates regex pattern
- Returns 400, error: "invalid MST format — must be 10 or 13 digits"

### Alternative Path 1c — MSDN not found in DPI database:
- System performs optional MSDN lookup (if integration exists)
- Returns 400, warning: "MSDN not verified — proceed anyway? (Y/N)"
- User confirms, company created with MSDN_UNVERIFIED flag

### Alternative Path 1d — Regime ineligible:
- If TT58 (micro-enterprise): company must have < 10 employees or revenue < 3B
- System checks eligibility, returns 400 if ineligible: "company does not meet micro-enterprise criteria for TT58"

### Exception Path 1e — Database error:
- System rolls back transaction
- Returns 500, logs error, alerts admin

### Business Rules
- MST is primary identifier, must be unique per tenant
- Legal form determines applicable regulations (Enterprise Law, Investment Law)
- Fiscal year must be 12 months (short first year allowed for mid-year registration)
- Regime change only at fiscal year boundary (per Circular 99 Art. 7)
- Default currency = VND (mandatory), secondary currency optional

---

## UC-02: Create Branch (Chi nhanh)

**Actor:** Chief Accountant
**Precondition:** Authenticated with `company:branch:create` permission, parent company exists

### Happy Path
1. User selects parent company
2. User enters: branch name, MST chi nhanh (13 digits: MST goc + -XXX)
3. System validates MST branch format (parent MST + hyphen + 3 digits)
4. User enters: address, phone, branch manager
5. System creates branch under parent company hierarchy
6. Branch inherits parent's accounting regime, fiscal year config
7. System logs audit entry
8. System returns 201

### Alternative Path 2a — MST branch format invalid:
- Returns 400: "branch MST must be parent MST + hyphen + 3 digits"

### Alternative Path 2b — Parent company not found:
- Returns 404: "parent company not found"

### Exception Path 2c — Parent company SUSPENDED or DISSOLVED:
- Returns 400: "cannot create branch under inactive company"

### Business Rules
- Branch inherits parent company's regime and fiscal year
- Branch has its own MST but shares parent's legal entity status
- Branch can have its own bank accounts and e-invoice patterns
- Branch transactions visible in parent consolidation

---

## UC-03: Configure Fiscal Year

**Actor:** Chief Accountant
**Precondition:** Authenticated with `company:fiscalyear:write` permission

### Happy Path
1. User selects company
2. User sets fiscal year start month (1-12)
3. System validates: no existing journal entries in proposed year
4. User confirms
5. System auto-generates 12 monthly periods + 4 quarterly periods
6. System sets first period as OPEN, remaining as FUTURE
7. System logs audit entry
8. System returns 200

### Alternative Path 3a — Journal entries exist in current year:
- Returns 400: "cannot change fiscal year — journal entries exist. Change at year boundary only."

### Alternative Path 3b — Short year for mid-year registration:
- User selects "short first year"
- System generates periods from registration month to fiscal year end
- System logs "SHORT_FISCAL_YEAR" audit entry

### Exception Path 3c — Invalid month (not 1-12):
- Returns 400: "fiscal year start must be month 1-12"

### Business Rules
- Fiscal year change only at boundary of current fiscal year (Circular 99 Art. 16)
- Must be consistent for at least one accounting year
- Short year only allowed for newly registered companies

---

## UC-04: Close Period (Khoa so)

**Actor:** Chief Accountant
**Precondition:** Authenticated with `company:period:close` permission

### Happy Path
1. User selects company and period (monthly or quarterly)
2. System validates: all journal entries are posted and balanced
3. System validates: no un-submitted tax declarations for period
4. System validates: sub-ledgers (AR, AP, inventory) are reconciled
5. User confirms close
6. System marks period = CLOSED
7. System auto-opens next period (if exists)
8. System logs audit entry with timestamp + user
9. System returns 200

### Alternative Path 4a — Unbalanced journal entries:
- Returns 400, lists unbalanced entries: "journal entry JE-2026-0012 not balanced (debit != credit)"

### Alternative Path 4b — Period already closed:
- Returns 409: "period already closed"

### Alternative Path 4c — Permanent close (Khoa so vinh vien):
- For annual close: system marks PPERMANENTLY_CLOSED
- No reopen permitted (per accounting law Art. 13)
- System generates annual financial statement data

### Exception Path 4d — Database error during close:
- System rolls back, returns 500
- Period remains OPEN, admin alerted

### Business Rules
- Once PERMANENTLY_CLOSED, period cannot be reopened (Law on Accounting Art. 17)
- Monthly close: optional but recommended
- Annual close: mandatory for financial statement preparation
- Closed period can be reopened (not permanently) with audit trail

---

## UC-05: Register E-Invoice Pattern

**Actor:** Chief Accountant
**Precondition:** Authenticated with `company:einvoice:write` permission

### Happy Path
1. User selects company
2. User enters: invoice pattern code (Mau so HD), serial (Ky hieu), description
3. System validates pattern format per GDT regulations (Decree 123/2020)
4. User selects: invoice type (GTGT/export/retail/etc.), form (Kieu HD)
5. System registers pattern with company profile
6. System sets status = REGISTERED
7. System logs audit entry
8. System returns 201

### Alternative Path 5a — Pattern already registered:
- Returns 409: "this invoice pattern already registered for company"

### Alternative Path 5b — GDT integration enabled:
- System pushes registration to GDT portal
- If GDT responds OK: status = ACTIVE
- If GDT fails: status = REGISTERED (pending manual sync)

### Exception Path 5c — Invalid pattern format:
- Returns 400: "pattern code format invalid — must follow GDT specification"

### Business Rules
- E-invoice required for all VAT-taxable enterprises (Decree 123/2020)
- Each company can have multiple patterns (e.g., GTGT + export)
- Pattern must be registered before invoice issuance
- GDT notification required within 5 days of change

---

## UC-06: Register Digital Signature

**Actor:** Chief Accountant or Admin
**Precondition:** Authenticated with `company:signature:write` permission

### Happy Path
1. User selects company
2. User selects: USB token or Remote HSM
3. For USB token: user inserts token, system reads serial number from provider API
4. For Remote HSM: user selects provider (ViettelCA/VNPTCA/CloudCA/others), enters credentials
5. User enters: owner name, valid from/to date
6. System validates: certificate chain, expiry date
7. System registers signature with company
8. System sets status = ACTIVE
9. System sets 30-day expiry alert
10. System returns 201

### Alternative Path 6a — Certificate expired:
- Returns 400: "digital certificate expired on YYYY-MM-DD"

### Alternative Path 6b — Remote HSM authentication failure:
- Returns 401: "remote HSM provider authentication failed — check credentials"

### Exception Path 6c — USB token not detected:
- Returns 400: "USB token not found — check connection or driver"

### Business Rules
- Each company must have >= 1 active digital signature for tax filing (Decree 23/2025)
- USB token: physical possession required
- Remote HSM: supported for enterprises with high-volume filing
- Expiry check at every tax submission
- Multiple signatures allowed per company (e.g., different owners)

---

## UC-07: Register Bank Account

**Actor:** Accountant or Chief Accountant
**Precondition:** Authenticated with `company:bank:write` permission

### Happy Path
1. User selects company
2. User selects bank from Napas bank list (or enters manually)
3. User enters: branch name, account number, account holder, currency (VND/USD/EUR)
4. System validates: account number format per bank rules
5. User sets as default bank account (optional)
6. System registers bank account, status = ACTIVE
7. System logs audit entry
8. System returns 201

### Alternative Path 7a — Bank not in system list:
- User enters bank name manually
- System sets bank_verified = false
- Admin verification required later

### Alternative Path 7b — Bank account verification (trial deposit):
- System initiates trial deposit (micro-amount)
- User confirms receipt amount
- If match: bank_verified = true
- If no match after 24h: verification fails

### Business Rules
- Multiple bank accounts per company
- Default bank account for cash/book entries
- Bank account format varies per bank (Napas standard where possible)
- Account closure: set status = CLOSED, no deletion

---

## UC-08: Switch Company Context

**Actor:** Any authenticated user with access to multiple companies
**Precondition:** Authenticated, assigned to multiple companies

### Happy Path
1. User calls `POST /api/v1/companies/switch` with target company_id
2. System validates user has access to target company
3. System updates session context: current_company = target company
4. System returns company context: company info, fiscal year, regime, default settings
5. All subsequent API calls use current_company context

### Alternative Path 8a — User not assigned to target company:
- Returns 403: "user does not have access to this company"

### Alternative Path 8b — Target company inactive:
- Returns 400: "cannot switch to inactive company"

### Business Rules
- Session context expires with JWT token
- Company switch is idempotent
- Audit log: COMPANY_SWITCH event recorded
- Rate limit: max 10 switches per minute (prevent abuse)

---

## UC-09: Create Employee Record

**Actor:** HR Manager or Accountant
**Precondition:** Authenticated with `company:employee:create` permission

### Happy Path
1. User selects company, department
2. User enters: employee code, full name, title, email, phone
3. User enters optional: MST ca nhan, BHXH number, bank account
4. System validates: employee code unique within company
5. System creates employee, status = ACTIVE
6. System links employee to department
7. System optionally links to system user (if user account exists)
8. System returns 201

### Alternative Path 9a — Employee already exists (same MST ca nhan):
- Returns 409: "employee with this personal tax code already exists"

### Alternative Path 9b — Department not found:
- Returns 400: "department not found in this company"

### Alternative Path 9c — Link to system user:
- User provides user_id
- System validates user exists and not already linked to another employee
- System links employee ↔ user (1:1)

### Business Rules
- Employee code unique per company
- Employee can belong to one primary department
- Employee's MST ca nhan used for PIT declaration
- Employee linked to system user for login (optional)
- Employee data protected under Decree 13/2023 (PDP)

---

## UC-10: Generate IGAP (Internal Accounting Policy)

**Actor:** Chief Accountant
**Precondition:** Authenticated with `company:igap:generate` permission

### Happy Path
1. User selects company
2. User clicks "Generate IGAP"
3. System collects: company info, accounting regime, COA, fiscal year config, accounting methods
4. System generates IGAP document in PDF format
5. System includes: company profile, regime basis, COA full list, revenue recognition policy, expense policy, fixed asset policy, inventory policy, foreign currency policy
6. System logs audit entry
7. System returns PDF download

### Alternative Path 10a — IGAP template not configured:
- System uses default template
- User can customize after generation

### Business Rules
- IGAP required per Circular 99 Art. 4 (Accounting Policy Regulation)
- Must be approved by company director
- Must be reviewed annually
- Changes must be documented with version history

---

## UC-11: Configure Integration Profile (GDT)

**Actor:** Admin or Chief Accountant
**Precondition:** Authenticated with `company:integration:write` permission

### Happy Path
1. User selects company, integration type = GDT (HTKK/iHTKK)
2. User enters: GDT username, password (encrypted), tax office code
3. User optionally enters: proxy settings, timeout config
4. System validates: test connection to GDT sandbox
5. If OK: integration status = CONNECTED
6. If FAIL: user can save with status = DISCONNECTED (retry later)

### Alternative Path 11a — Connection test fails:
- Returns 400: "GDT connection test failed — check credentials or network"
- User can save anyway with PENDING status

### Business Rules
- Credentials stored encrypted at rest (AES-256-GCM)
- Integration secret never returned in API responses
- Regular connection health check (cron every 4 hours)
- GDT portal credentials managed per company
