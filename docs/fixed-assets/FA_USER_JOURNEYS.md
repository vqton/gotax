# Fixed Asset Module — User Journeys

**Version:** 1.0
**Date:** July 2026
**Author:** BA Lead + Chief Accountant (20+ yrs each)

**Regulatory references:**
- Circular 45/2013/TT-BTC (TT 45) / TT 147/2024 — FA management
- VAS 03 — Tangible Fixed Assets, VAS 04 — Intangible FA
- Circular 99/2025/TT-BTC — Chart of Accounts
- Circular 200/2014/TT-BTC — Accounting regime
- Decree 123/2020/ND-CP — E-invoices
- IAS 16, IAS 36 — Property, Plant & Equipment, Impairment

**Use case references:** See `FA_USE_CASES.md` for detailed UC specs.

---

## UJ-1: FA Accountant (Kế toán TSCĐ) — Monthly Depreciation Run

**Persona:** Mr. Hoang, FA Accountant at a manufacturing company, 6 years experience. Manages 1,200+ FA items across 3 factories.

**Goal:** Process accurate monthly depreciation for all assets by Day 3 of each month. Register new FA within 2 days of receipt.

**Frustration:** Spreadsheet-based depreciation calc (15 separate Excel files, one per cost center). Manual GL journal entry. No approval workflow — chief accountant rubber-stamps everything at month-end. Mid-period transfers cause allocation nightmares.

### Journey — Before GoTax

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Day 1, 08:00 | Check email: HR confirms new laptop for Sales Dept arrived | Email from HR | Neutral |
| 08:15 | Collect paper delivery receipt from warehouse, find supplier invoice in email | Paper trail | Annoyed |
| 08:30 | Open Excel "FA Register 2026.xlsm", add row for new laptop | Manual entry | OK |
| 08:35 | Calculate: cost = 35M, useful life = 48 months → monthly dep = 729,167 VND | Mental math | Distracted |
| 08:40 | Email chief accountant for approval (no response expected until Day 5) | Email | Anxious |
| 09:00 | Open depreciation Excel #7 (Factory A cost center 6274) | File hunt | Frustrated |
| 09:15 | Copy-paste prior month column, adjust for 5 new assets, 1 disposal | Manual calc | Tired |
| 09:30 | Cell B127: forgot to update residual value of Machine TS-2022-0318 | Error caught | Alarmed |
| 10:00 | Mid-period transfer: forklift moved from Factory A to Factory B on 12th | Calculator out | Stressed |
| 10:15 | Prorata calculation: 11 days A, 19 days B. Correct? Not sure. | Second-guessing | Uncertain |
| 16:00 | 14 Excel files done, 1 more + GL entry remaining | Overwhelmed | Exhausted |
| Day 2 | Run Python script to generate GL journal from Excel → import to accounting system | Technical jank | Nervous |
| Day 3 | Chief accountant reviews: finds 2 errors (wrong account for Factory B, missed new asset) | Corrections | Embarrassed |
| Day 5 | Depreciation finally posted. Close FA register. Month-end done. | Deadline missed | Burned out |

**Mid-month event:** Factory A calls: "We transferred Machine XYZ to Factory C last week." Hoang sighs. Update allocation spreadsheet. Send email to chief accountant "FYI, dep allocation changed." No confirmation. Hope it's right.

**Quarterly pain:** Physical inventory = print 50-page PDF, walk factory with clipboard, hand-write discrepancies, scan papers back.

### Journey — After GoTax

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Day 1, 08:00 | New laptop arrives. Scan delivery QR → auto-creates FA draft from PO | Mobile scan | Fast |
| 08:05 | Open GoTax → FA → Pending Registration. System pre-filled: cost 35M, useful life 48m, method SL, account 6424 from category defaults | Auto-calc | Happy |
| 08:06 | Validate, attach e-invoice PDF, click Submit for Approval | One click | Efficient |
| 08:07 | Approval request sent to chief accountant | Instant | Done |
| 08:10 | Open GoTax → FA → Run Depreciation for July | Dashboard | Ready |
| 08:12 | Select period: Jul 2026. Click "Calculate" | 2 seconds | Fast |
| 08:13 | System shows preview: 1,245 assets, total dep 1.87B, grouped by expense account | Clear table | Confident |
| 08:14 | Spot-check: Machine TS-2022-0318 residual value correct ✓ | Quick verify | Trusting |
| 08:15 | Mid-period transfer: forklift moved 12 Jul. System auto-prorated (11d A, 19d B). Correct. ✓ | Auto-calc | Impressed |
| 08:16 | Click "Post to GL" → Journal JE-2026-07812 created. Accumulated dep updated. | Done | Satisfied |
| 08:20 | New FA approved notification → ACTIVE, dep starts next month | Auto-activate | Smooth |
| Day 1, 09:00 | Done with depreciation. 45 minutes total. Rest of month for other work. | Time saved | Happy |

**Emotional arc:** ![Frustrated at Excel mess] → [One-click registration] → [Auto-depreciation] → [Confident and early finish]

**Before:** 5 days, manual, errors, stress. **After:** 45 minutes, zero errors, audit trail.

**Pain points solved:**
- Excel hell: 15 spreadsheets eliminated → single click
- Mid-period transfer prorata: auto-calculated, no guessing
- No approval workflow: instant routing, auto-activation
- Manual GL entry: auto-posted to correct accounts
- Error risk: validation rules prevent mistakes

**Success metrics:**
- Month-end close: Day 1 (before: Day 5) — target < 2 days
- Error rate: 0% (before: 2-3 errors/month)
- New FA registration: < 1 day from receipt (before: 3-5 days)
- Time per asset processed: 2 min (before: 15 min)

---

## UJ-2: Chief Accountant (Kế toán trưởng) — FA Policy & Approval

**Persona:** Ms. Hoa, Chief Accountant, 22 years experience. Oversees 3 accountants, signs off on FA register, approves depreciation method changes.

**Goal:** Ensure FA accounting is accurate, compliant, and audit-ready. Approve only what needs approval. Spend time on analysis, not clerical reviews.

**Frustration:** Buried in paper approval requests. No dashboard showing FA health. Cannot see pending approvals without asking. Tax accountant asks for FA data separately. Hard to assess if useful life assumptions are still valid.

### Journey — Before GoTax

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Day 2, 09:00 | Inbox: FA accountant's email "Please approve July depreciation" | Email | Annoyed |
| 09:05 | Open attached Excel "Depreciation Jul 2026.xlsx" (15 tabs) | Heavy file | Overwhelmed |
| 09:10 | Scan totals: 1.87B dept A, 892M dept B. Compare prior month: dept A jumped 15% | Flag | Suspicious |
| 09:15 | Reply: "Why is dept A 15% higher?" Wait for response | Email ping-pong | Inefficient |
| 09:30 | FA accountant replies: "5 new machines installed. See tab 'New Assets'." | Explanation | Slightly better |
| 09:35 | Scroll to New Assets tab. Verify: 5 machines, costs OK, useful lives within norms | Manual check | Tedious |
| 09:40 | Reply: "Approved." No digital signature. No audit trail of approval. | Email approval | Weak |
| 10:00 | Paper inbox: 3 FA registration forms (printouts with physical signatures) | Paper stack | Old-school |
| 10:05 | Review each: costs 45M, 120M, 380M. Check useful lives against TT 45 schedule. Sign. | Signing | OK |
| 10:15 | Stack for assistant to scan back to FA accountant | Scanning delay | Inefficient |
| 14:00 | Tax accountant knocks: "Need FA register for CIT finalization" | Walk-in | Interrupted |
| 14:05 | Forward FA accountant's latest Excel (3 weeks old). "Is this current?" | Outdated data | Embarrassed |
| 14:30 | Approval request: "Change depreciation method for Press Machine from DB to SL" | Email request | Cautious |
| 14:35 | No data given: why change? What's the impact? Remaining life? | Missing context | Frustrated |
| 14:40 | Reply: "Send impact analysis." Wait. | Back-and-forth | Slow |

**Quarterly:** Pull FA register from Excel. Reconcile with GL manually. Find 2 discrepancies. Chasing.

### Journey — After GoTax

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Day 1, 08:00 | Dashboard opens: "FA Summary" widget | Visual overview | Informed |
| 08:01 | Total FA: 1,245. Carrying amount: 42.8B. This month dep: 1.87B. ✓ vs forecast | Auto | Confident |
| 08:02 | Pending approvals: 3 registrations, 1 method change | Badge count | Clear |
| 08:05 | Open registration #1: Factory C CNC machine, 380M, SL 120 months | Review screen | Fast |
| 08:06 | System shows: cost ±5% vs similar assets, useful life within category range | Auto-validation | Trusting |
| 08:07 | Click Approve → FA activated, GL posted, depreciation starts next month | One click | Efficient |
| 08:10 | Registration #2: useful life 180 months (category default is 120) — flag | Anomaly | Attentive |
| 08:11 | System: "Useful life 50% > category max. Reason: extended production plan." | Context shown | Informed |
| 08:12 | Accept. Approve with note. | Documented | Compliant |
| 08:15 | Method change: Press Machine DB → SL. System shows impact: monthly dep 12M→8.5M, total dep over remaining life 102M vs 144M | Side-by-side | Clear |
| 08:16 | Reason: "Asset utilization reduced." QA check: "Prospective only per VAS 03." | Validate | Compliant |
| 08:17 | Approve. System updates schedule prospectively. | Instant | Done |
| 08:20 | All approvals done in 20 minutes. | Zero backlog | Satisfied |

**Mid-quarter:** Opens GoTax → Reports → FA Aging. Sees 12 assets within 6 months of full depreciation. Emails CFO: "Plan replacement budget for these 12 assets." Proactive instead of reactive.

**Emotional arc:** ![Buried in paper] → [Dashboard clarity] → [One-click approvals with context] → [Proactive asset planning]

**Pain points solved:**
- No dashboard: FA summary + pending approvals at a glance
- Context-free approvals: system shows validation, impact, and regulatory flags
- Paper routing: digital approval with audit trail
- No FA data sharing: role-based access, real-time reports
- Tax reconciliation: separate tax book tracked alongside accounting book

**Success metrics:**
- Approval turnaround: < 1 hour (before: 2-3 days)
- Pending approvals: zero at any time (before: 5-10 pending)
- FA tax compliance: 100% (before: annual CIT adjustments needed)
- Time spent on FA approval: 20 min/day (before: 2 hours)

---

## UJ-3: CFO / Financial Manager — FA Investment Decisions

**Persona:** Mr. Duc, CFO, 15 years experience. Responsible for capital budgeting, investment decisions, financial reporting accuracy.

**Goal:** Know FA value, age, utilization. Plan replacement capex. Ensure FA reporting is accurate for board and investors.

**Frustration:** No FA analytics. Cannot answer "What's our FA age profile?" without IT team running queries. Capital budget requests come without FA utilization data. Impairment risk unknown.

### Journey — Before GoTax

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Week 1 | Board prep: "Show us our fixed asset position" | Request | Neutral |
| Week 1 | Email chief accountant: "Please send FA report" | Email | Dependent |
| Day 3 | Receive Excel with 20 columns, raw data, no formatting | Raw dump | Frustrated |
| Day 3 | Build pivot: FA by category, original cost, accumulated dep, carrying amount | Manual pivot | Slow |
| Day 3 | Calculate age: avg useful life used / total useful life per asset | Manual calc | Tedious |
| Day 4 | Notice Factory C machines: 85% depreciated → replacement needed soon | Insight | Concerned |
| Day 4 | No utilization data. Are they running full capacity? No way to know. | Missing data | Incomplete |
| Day 5 | Capital planning: Factory manager requests 2B for new press machine | Request received | Skeptical |
| Day 5 | No benchmark: current press age? utilization? maintenance cost history? | No data | Blind |
| Day 5 | Approve subjectively based on "seems reasonable" | Gut feel | Risky |
| Quarter-end | Balance sheet review: FA carrying amount 42.8B. Is impairment needed? | Question | Uncertain |
| Quarter-end | Email for impairment assessment → chief accountant says "no indicators" | Trust but verify | Weak |

**Year-end:** FA movement schedule prepared by FA accountant in Excel. CFO reviews, finds discrepancy in "Additions" column (1.2B vs GL 1.1B). Chasing.

### Journey — After GoTax

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Monday | Open GoTax → Reports → FA Dashboard | Self-service | Independent |
| 08:01 | FA Overview: 1,245 assets, 82.0B original cost, 42.8B carrying amount | Key metrics | Clear |
| 08:02 | Aging: 28% fully depreciated (still in use), 35% 50-99% used, 25% 25-50%, 12% <25% | Bar chart | Insightful |
| 08:03 | Filter: Factory C machines → 85% depreciated → avg remaining life 18 months | Drill-down | Focused |
| 08:05 | Click "Replacement Projection" → forecast for next 24 months: 2.8B needed | Auto-forecast | Prepared |
| 08:10 | Export FA aging report to PDF for board pack | One click | Professional |
| Week 2 | Factory manager requests 2B for press machine | Request | Neutral |
| Week 2 | Open GoTax → current press FA record: 10 years old, 95% depreciated, carrying amount 50M | Data at hand | Informed |
| Week 2 | Utilization log: press ran 75% capacity last year, 60% this year (declining) | Trend | Data-driven |
| Week 2 | Maintenance cost history: 150M in repairs last year alone vs 50M/year average | Cost signal | Clear |
| Week 2 | Approve replacement with full evidence pack for board | Confident decision | Satisfied |
| Quarter-end | System auto-checks impairment: all FAs recoverable amount > carrying amount | Auto-report | Compliant |
| Quarter-end | FA movement schedule generated: Opening + Additions - Disposals - Dep = Closing | Auto-matches GL | Verified |

**Half-year review:** CFO reviews FA utilization report: 23 assets idle > 6 months. Flags for chief accountant: "Dispose or reactivate." Proactive FA optimization.

**Emotional arc:** ![Waiting on email reports] → [Self-service dashboard] → [Drill-down analytics] → [Confident capex decisions]

**Pain points solved:**
- No FA analytics: real-time dashboard with aging, utilization, replacement forecast
- Blind capital requests: full asset history, utilization, maintenance cost per FA
- Manual movement schedule: auto-generated, GL-reconciled
- Impairment blind: auto-assessment with indicators
- No utilization tracking: idle asset identification → disposal or reallocation

**Success metrics:**
- FA ROI visibility: 100% of assets tracked with cost & utilization (before: 0%)
- Capital budget accuracy: ±10% of actual (before: ±30%)
- Impairment assessment: quarterly, automated (before: annual, ad-hoc)
- Board report preparation: 30 min (before: 3 days)
- Idle asset identification: real-time (before: never)

---

## UJ-4: Tax Accountant (Kế toán thuế) — FA Tax Compliance

**Persona:** Ms. Lan, Tax Accountant, 6 years experience. Handles VAT, CIT, e-invoices for a mid-size trading company.

**Goal:** Ensure all FA transactions have correct tax treatment. Maximize VAT deduction. Prepare CIT finalization with correct FA adjustments. Avoid GDT audit adjustments.

**Frustration:** Manual tracking of VAT input on FA purchases. FA tax depreciation often different from accounting dep — hard to reconcile. E-invoice for FA sale is a separate process. CIT finalization requires digging through FA records from two books.

### Journey — Before GoTax

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Day 5 | Monthly VAT declaration: 45 purchase invoices, need to identify FA-related VAT | Print all invoices | Tedious |
| 09:00 | Filter: any invoice with account 211x or over 30M | Manual scan | Slow |
| 09:15 | Find 8 FA purchase invoices, total VAT input 45M | Found | OK |
| 09:20 | Check: all 8 have valid e-invoices? 3 invoices > 20M, check bank payment evidence | Per-declaration check | Cautious |
| 09:30 | Invoice #4: FA purchase but not yet received (goods in transit). Can I claim VAT? | Uncertainty | Confused |
| 09:35 | Check Decree 123: VAT claimable when goods received or service completed, regardless of payment (except > 20M) | Research break | Distracted |
| 09:40 | Wait — this FA is CIP (2411), not yet 211. VAT still claimable? | More research | Frustrated |
| 09:50 | Yes, VAT on CIP is claimable (1332) per TT 219. Proceed. | Confidence | Relieved |
| 10:00 | Enter VAT declaration: FA-related VAT input = 45M. Cross-check with FA register. | Manual cross-check | Repetitive |

| Month 3 | CIT quarterly: need to compute CIT adjustment for FA depreciation | Heavy task | Dreaded |
| 09:00 | Export FA register from accounting system (report: "TSCD depreciation schedule") | Report | Neutral |
| 09:05 | Export tax FA book from tax system (separate software) | Second system | Annoyed |
| 09:10 | Compare line by line: accounting dep vs tax dep for each asset | Excel VLOOKUP | Error-prone |
| 09:30 | Find 12 differences: different methods (SL vs Accelerated), different useful lives | Diffs found | Concerned |
| 09:35 | Calculate total: accounting dep 1.87B, tax dep 2.45B → temporary difference = 580M | Manual calc | Tired |
| 09:40 | CIT adjustment: add back 580M to taxable income (since tax dep > accounting dep = future taxable) | Journal entry | Done |
| 09:45 | Disclosure: prepare FA tax footnote for CIT finalization | Copy tables | Repetitive |

| FA sale | E-invoice for building sold: value 5B, carrying amount 3.2B, gain 1.8B | Transaction | Critical |
| 10:00 | Check: is this building subject to VAT? Building sold after 2 years → VAT-exempt per Law 71 | Regulation check | Careful |
| 10:05 | Issue e-invoice: VAT-exempt sale, reason code... | GDT portal | Separate process |

### Journey — After GoTax

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Day 1 | Monthly VAT: Open GoTax → Reports → FA VAT Input | Report | Ready |
| 08:00 | System lists 8 FA purchase invoices for the month, VAT input 45M, all with valid e-invoices ✓ | Auto-list | Effortless |
| 08:01 | System flags: Invoice #4 is CIP (2411). "VAT claimable per TT 219." ✓ | Guidance | Confident |
| 08:02 | System: "3 invoices > 20M — bank payment verified" ✓ | Auto-check | Compliant |
| 08:03 | Export VAT input data for declaration → auto-maps to HTKK format | One click | Fast |
| Day 2 | CIT quarterly: Open GoTax → Reports → FA Tax Reconciliation | Dashboard | Ready |
| 08:00 | System shows: Accounting dep 1.87B vs Tax dep 2.45B — difference 580M | Side-by-side | Clear |
| 08:01 | 12 differences listed: asset by asset with reason (method, life, residual value) | Drill-down | Transparent |
| 08:02 | Auto-calculated CIT adjustment: add back 580M to taxable income | Already done | Trusting |
| 08:03 | Generate FA disclosure note for CIT finalization: template with figures | Auto-fill | Complete |
| FA sale | Process building disposal: 5B proceeds, carrying 3.2B | Disposal flow | Guided |
| 09:00 | System: "Building age > 2 years → VAT-exempt per Law 71. Use exemption code X." | Auto-flag | Compliant |
| 09:01 | Generate e-invoice from disposal → links to FA transaction | Integrated | Seamless |

**Year-end:** GoTax generates FA tax schedule: original cost + accumulated dep for tax book + carrying amount. Ready for CIT finalization attachment.

**Emotional arc:** ![Manual VAT filtering] → [Auto-FA-VAT report] → [Two-book reconciliation automated] → [Audit-ready tax disclosure]

**Pain points solved:**
- Manual VAT tracking: auto-identified FA purchase invoices with e-invoice verification
- Two-book reconciliation: FA accounting book + FA tax book side-by-side, differences auto-calculated
- CIT adjustment: auto-computed deferred tax impact
- E-invoice for FA sale: integrated disposal → e-invoice pipeline
- CIT disclosure: auto-generated FA note with two-book reconciliation

**Success metrics:**
- Tax filing accuracy: 100% (before: quarterly CIT adjustments needed)
- VAT deduction: 100% claimed correctly (before: 2-3 missed per quarter)
- Two-book reconciliation: instant (before: 1 day)
- Audit readiness: full FA tax trail available anytime (before: manual compilation)
- CIT finalization preparation: 2 hours (before: 2 days)

---

## UJ-5: Internal Auditor (Kiểm toán nội bộ) — FA Physical Verification

**Persona:** Mr. Quang, Internal Auditor, 8 years experience. Performs annual FA physical inventory for a group of 3 companies.

**Goal:** Complete annual FA physical inventory within 5 working days. Identify discrepancies. Report to management with actionable findings.

**Frustration:** Paper-based inventory lists. Print 100+ pages per company. Hand-write results. Hard to track real-time progress. Post-inventory data entry takes 3 days. Managers don't act on findings.

### Journey — Before GoTax

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Week 1, Mon | Plan annual inventory: print FA register from accounting system (480 pages for 3 companies) | Print station | Material waste |
| 08:00 | Bind 3 books: Company A (200 pages), B (150), C (130) | Binding | Physical |
| 08:30 | Walk to Factory A with clipboard, inventory list, pen | Field gear | Ready |
| 09:00 | Asset #1: CNC Machine TS-2021-0012 → found, condition good → checkmark | Physical check | OK |
| 09:05 | Asset #2: Forklift TS-2022-0031 → found, damaged fork → note "damaged" | Write note | Concerned |
| 09:30 | Asset #15: Computer TS-2025-0008 → not found → check 3 desks, no | Search | Stressed |
| 09:35 | Mark "MISSING" on paper. Later investigate. | Paper note | Uncertain |
| 12:00 | Lunch break. Count: 45/480 assets done. Progress: 9%. | Slow pace | Pressure |
| Day 3 | Factory A done (200 assets). Handwriting getting messy. | Fatigue | Tired |
| Day 3 | Return to office. Start data entry: type results from paper into Excel | Transcription | Error-prone |
| Day 4 | Discover: 1 asset marked "FOUND" but paper shows next to "DAMAGED" — unclear handwriting | Ambiguity | Frustrated |
| Day 5 | All entered. Generate discrepancy report: 5 missing, 3 damaged, 2 unregistered | Excel pivot | Done |
| Day 5 | Email report to chief accountant + CFO. No reply for 2 weeks. | Dead email | Disappointed |
| Day 5 | Discrepancies follow-up: "We'll investigate." No deadline given. | Vague response | Unsatisfied |
| Month 3 | Follow up: missing assets still not resolved. Report filed. No action. | Stale findings | Defeated |

### Journey — After GoTax

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Mon, 08:00 | Open GoTax → FA → Physical Inventory → Create Plan | Digital | Ready |
| 08:02 | Select company: Company A. Scope: All active assets (location = Factory A) | Filter | Efficient |
| 08:03 | System generates inventory list: 200 assets, QR-coded, ready for mobile | Mobile-friendly | Modern |
| 08:05 | Open GoTax mobile app on phone → Inventory Plan → Start | Field ready | Prepared |
| 09:00 | Walk Factory A. Asset #1: CNC Machine. Scan QR tag on machine → phone shows FA detail | Barcode scan | Fast |
| 09:01 | Confirm: location correct ✓, condition Good → swipe Found | 3 seconds | Instant |
| 09:05 | Asset #2: Forklift. Scan QR. Condition: Damaged → tap "Damaged", add photo of fork | Photo evidence | Accurate |
| 09:30 | Asset #15: Computer TS-2025-0008. No QR found. Search nearby desks. Not found. | Physical search | Concerned |
| 09:35 | In app: Mark "MISSING" → auto-create investigation case | One tap | Systematic |
| 12:00 | Progress: 65/200. System dashboard shows: 32% done, estimated finish: Wed PM | Real-time tracking | Motivated |
| Wed, 15:00 | Factory A done: 200/200. System: "Inventory complete for Factory A." | Green checkmark | Satisfied |
| Wed, 15:05 | Results: 192 Found ✓, 3 Damaged (photos attached), 5 Missing (investigation open) | Auto-summary | Clear |
| Wed, 15:06 | System generates inventory report: comparison with prior year, discrepancy % | Auto-report | Professional |
| Wed, 15:10 | Click "Submit to Chief Accountant" → notification sent with report attached | Digital | Complete |
| Thu, AM | Chief accountant reviews discrepancies: assigns 5 missing to FA accountant for investigation | Action assigned | Accountable |
| 14 days | Investigation results: 3 missing = disposed without record (adjust), 2 found in other location (transfer) | Resolution track | Closed |
| 30 days | All discrepancies resolved, management approved adjustments | Audit cleared | Satisfied |

**Emotional arc:** ![Paper clipboard fatigue] → [Mobile QR scanning] → [Real-time progress] → [Discrepancies tracked to resolution]

**Pain points solved:**
- Paper inventory lists: mobile app with QR code scanning
- Manual data entry: results recorded in field, zero transcription
- No real-time progress: live dashboard shows completion % across locations
- Unclear handwriting: structured input with photo evidence
- No follow-up: auto-assigned investigations with deadlines and reminders
- Slow discrepancy resolution: automated workflow with SLA tracking

**Success metrics:**
- Inventory completion: 100% within 3 days (before: 5-7 days)
- Data entry errors: 0 (before: 5-10 transcription errors)
- Discrepancy resolution: < 30 days (before: often unresolved)
- Audit trail: full digital trail with photos (before: paper, lost)
- Time saved: 60% reduction (3 days vs 7 days for 3 companies)

---

## UJ-6: External Auditor (Kiểm toán độc lập) — FA Schedule Substantiation

**Persona:** Mr. Phuc, External Auditor, 10 years experience at a Big4 firm. Leading the FA audit for a manufacturing client.

**Goal:** Obtain sufficient evidence that FA exists, is properly valued, depreciation is accurate, and no material misstatement.

**Frustration:** Client sends Excel exports that need reformatting. Depreciation testing requires re-performing calculations manually. Vouching for FA existence requires hours of document requests. Impairment assessment is usually a "check the box" with no real analysis.

### Journey — Before GoTax

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Day 1, 08:00 | FA Planning meeting with chief accountant | Meeting | Organized |
| 08:30 | Request: FA register, depreciation schedule, FA movement schedule, FA tax schedule | Standard list | Routine |
| 09:00 | Client sends: 4 separate Excel files (different formats, one export buggy) | Data dump | Annoyed |
| 09:30 | Reformat into audit workpaper template | Manual work | Tedious |
| 10:00 | Reconcile FA register vs GL: 211x balance = 82B, FA register total = 81.7B → diff 300M | Variance | Concerned |
| 10:05 | Ask client: "Why 300M difference?" Client: "Let me check." | Follow-up | Waiting |
| Day 2 | Client reply: "One asset 300M was disposed last week, GL posted but register not updated." | Correction | Relieved |
| 10:00 | Depreciation testing: select sample of 20 assets | Sampling | Methodical |
| 10:15 | For each: manually recalc monthly dep = (OriginalCost - ResidualValue) / UsefulLifeMonths | Calculator | Slow |
| 10:20 | Asset #5: declining balance method → formula more complex | Research | Time-consuming |
| 10:45 | 20 recalc done. All match client figures. ✓ But took 45 minutes. | Verified | Competent |
| 11:00 | Vouching: request invoices for 5 FA additions (> 500M each) | Evidence request | Standard |
| 11:05 | Client goes to filing cabinet. Returns 30 min later with 3 of 5 invoices. "Still looking for the others." | Paper hunt | Inefficient |
| 11:35 | 2 invoices found (scanned copies). 1 still missing. "Construction asset, we have contractor's completion cert." | Alternative evidence | Adjusting |
| 14:00 | FA existence: review annual inventory report (from client's internal audit) | Paper report | Passive |
| 14:30 | Impairment: client says "no indicators." No analysis provided. "We'll assess." | Thin evidence | Accepting |
| Day 3 | FA movement schedule: opening + additions - disposals - dep = closing ✓ | Re-perform | Verified |
| Day 3 | FA tax schedule: differences from accounting book noted for CIT footnote | Cross-check | Thorough |
| Day 4 | Audit conclusion: no material misstatement. One finding: FA register update lag. | Management letter | Done |

### Journey — After GoTax

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Day 1, 08:00 | FA audit planning: client grants auditor read-only access to GoTax | Portal access | Modern |
| 08:05 | Open GoTax → Reports → FA Audit Package → "Export for External Audit" | One-click | Fast |
| 08:06 | System generates single PDF/XLSX: FA register, depreciation schedule, movement schedule, tax schedule, GL reconciliation | Complete package | Impressed |
| 08:07 | GL reconciliation: 211x = 82.0B = FA register total 82.0B ✓ | Zero diff | Verified |
| 08:08 | Depreciation schedule: per-asset breakdown, method, monthly amount, YTD, accumulated | Detailed | Transparent |
| 08:10 | Depreciation testing: use GoTax "Audit Sample" tool — select 20 assets random sample | Built-in tool | Efficient |
| 08:15 | System shows: for each sample, recalc = posted, full working visible | Auto-verify | Trusting |
| 08:16 | Asset #5 (declining balance): click "Show calculation" → formula, inputs, monthly schedule visible | Drill-down | Verified in 2 min |
| 08:20 | Vouching: for 5 additions, click "Supporting Documents" → e-invoice, contract, delivery receipt linked per asset | Digital evidence | Paperless |
| 08:22 | Asset #3: e-invoice PDF opens in browser. Bank payment evidence linked. | Complete trail | Satisfied |
| 08:25 | FA existence: open GoTax "Physical Inventory Results" for current year | System record | Reliable |
| 08:26 | 98.5% verified by internal audit. 1.5% discrepancies tracked. | Verified | Sufficient |
| 08:30 | Impairment: open GoTax "Impairment Assessment" → system checks all assets | Auto-analysis | Robust |
| 08:31 | Result: no impairment indicators for any asset. Supporting: utilization report, market data references, maintenance records. | Documented | Strong |
| 08:35 | FA movement schedule: auto-generated, opening→additions→disposals→dep→closing. Cross-footed. | Built-in | Verified |
| Day 2 | Audit conclusion: no material misstatement. Control finding: none. | Clean opinion | Satisfied |

**Emotional arc:** ![Reformatting Excel dumps] → [Single audit package export] → [Auto-verification] → [Clean opinion, zero findings]

**Pain points solved:**
- Data export not auditor-friendly: audit package export with all schedules in one file
- Manual dep recalculation: built-in audit sample tool with auto-verification
- Vouching paper trail: digital evidence linked per asset (e-invoice, contract, payment)
- FA existence proof: digital inventory results with QR scan log
- Impairment assessment: auto-checks with documented support
- Manual GL reconciliation: auto-reconciled, zero-variance report

**Success metrics:**
- Audit fieldwork: 1 day (before: 4 days)
- Audit findings: 0 material misstatements (target)
- Sample testing: 20 assets in 10 min (before: 45 min)
- Evidence completeness: 100% digital, available on request (before: paper, partial)
- Client time with auditor: 1 hour kickoff only (before: 8+ hours answering requests)

---

## UJ-7: System Administrator (Quản trị hệ thống) — FA Module Setup

**Persona:** Ms. Thao, System Administrator / IT Manager, 7 years experience. Responsible for configuring GoTax modules for a newly onboarded company.

**Goal:** Configure FA module for first-time use: categories, accounts, approval workflows, opening balance import. Go live within 1 week.

**Frustration:** Complex setup with many interdependent configurations. No migration tool for legacy FA data. Hard to know if setup is correct before users start. Manual testing of depreciation calculation.

### Journey — Before GoTax (no FA module exists)

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Week 1 | Receive request: "Set up FA tracking for our new company" | New module | Motivated |
| Day 1 | No FA module exists — build from scratch or use external tool | Evaluate | Uncertain |
| Day 2 | Decide: use Excel for FA register, outsource depreciation calc to accounting firm | Workaround | Compromised |
| Day 2 | Set up Excel with columns: code, name, cost, dep method, useful life, monthly dep | Manual setup | Tedious |
| Day 3 | No integration with GL — depreciation entries must be manually journaled | Separate process | Inefficient |
| Day 3 | No approval workflow — all FA changes done offline | No control | Risky |
| Week 2 | Legacy FA data: 800 assets from old system. Export → reformat → paste into Excel | Data migration | Painful |
| Week 2 | No validation: 5 assets had cost below 30M (should not be FA). 3 duplicates found later. | Data quality | Errors |
| Week 3 | Accounting firm runs first depreciation. Takes 5 business days. | Outsourced | Costly |
| Week 4 | Chief accountant asks: "Can we get FA reports?" — need to build in Excel | Report building | Ongoing |

### Journey — After GoTax

| Time | Step | Touch Point | Emotion |
|------|------|-------------|---------|
| Day 1, 08:00 | Open GoTax → Admin → FA Module → Enable for company | Configuration | Fast |
| 08:05 | System prompts: "Configure FA categories first." | Guided setup | Clear |
| 08:10 | Open FA Category Management. System provides default categories from TT 45/TT 147: | Template | Helpful |
| | - Building (2111): SL, 25 yrs, account 6424 | | |
| | - Machinery (2112): SL, 10 yrs, account 6274 | | |
| | - Vehicle (2113): SL, 6 yrs, account 6414 | | |
| | - Office Equipment (2114): SL, 3-5 yrs, account 6424 | | |
| | - Computer (2115): SL, 3 yrs, account 6424 | | |
| | - Intangible (213): SL, up to 20 yrs, account 6424 | | |
| 08:15 | Customize: add "Solar Panel System" category, useful life 25 yrs, dep method SL | Add category | Flexible |
| 08:20 | System auto-creates account mapping based on categories → 2111/2141/6274, etc. | Auto-map | Efficient |
| 08:25 | Configure approval workflow: | | |
| 08:26 | FA registration > 500M → requires 2 approvers (chief accountant + CFO) | Treshold set | Tailored |
| 08:27 | Method change → chief accountant only. Disposal → chief accountant + CFO. | By operation | Configurable |
| 08:30 | Set default depreciation method: Straight-line (most common in Vietnam) | Policy | Standard |
| 08:35 | Configure number series: TS-YYYYMM-XXXX, starts at 0001 | Auto-number | Ready |
| 09:00 | Import opening balance: download template | Migration | Planned |
| 09:05 | Client sends legacy FA register: 800 rows in CSV | Data ready | Prepared |
| 09:10 | Map CSV columns to GoTax fields (drag-and-drop mapping) | Visual mapping | Intuitive |
| 09:15 | Upload: system validates 800 rows | Bulk validation | Processing |
| 09:20 | Results: 792 passed, 5 cost < 30M (warnings), 3 duplicates | Summary | Clear |
| 09:22 | Fix: skip 5 below-threshold assets (they'll be expensed). Merge 3 duplicates. | Quick fix | Resolved |
| 09:25 | Preview: total original cost 78.5B, total accumulated dep 35.7B, carrying amount 42.8B | Pre-import | Verified |
| 09:26 | Reconciliation: GL 211x opening balance = 78.5B ✓. GL 214x = 35.7B ✓. Match! | GL match | Confident |
| 09:27 | Click Import → 792 FA created with correct statuses | Imported | Satisfied |
| 09:30 | Status: 245 ACTIVE, 430 DEPRECIATING, 117 FULLY_DEPRECIATED | Distribution | Transparent |
| Day 2 | Test depreciation run: select first month → Calculate → 792 assets processed | Test | Verify |
| 09:00 | Check 5 test assets: manual recalc vs system → match ✓ | Accuracy test | Trusting |
| 09:15 | Run GL dry-run: no post, preview only → accounts correct ✓ | Dry-run | Safe |
| 09:30 | Enable FA module for users: assign roles (FA Accountant, Chief Accountant, etc.) | Role assignment | Complete |
| Day 3 | Train users: 2-hour session with FA accountant + chief accountant | Training | Done |
| Day 3 | Go live! First production depreciation run. | Live | Successful |

**Emotional arc:** ![No FA module, Excel chaos] → [Guided setup with defaults] → [One-click opening balance import] → [Go live in 3 days]

**Pain points solved:**
- No FA module: built-in module with comprehensive feature set
- Complex setup: guided wizard with TT 45/TT 147 default categories
- No migration tools: template + column mapping + validation + GL reconciliation
- No default policies: category defaults (method, useful life, accounts) auto-populated
- No approval setup: configurable workflow per operation type
- Manual testing: test run before go-live with dry-run mode
- No role control: granular permissions for FA operations

**Success metrics:**
- Time to configure: < 1 day (before: 1-2 weeks)
- Opening balance import: 800 assets in 30 min (before: manual 3-5 days)
- Validation accuracy: 100% of rules applied (before: errors found weeks later)
- Go-live readiness: Day 3 (before: 3-4 weeks)
- Training time: 2 hours (before: ad-hoc, weeks of follow-up questions)
- User adoption: 100% within first week

---

## Summary Matrix

| Journey | Actor | Key Pain Point | GoTax Solution | Time Saved |
|---------|-------|----------------|----------------|------------|
| UJ-1 | FA Accountant | 15 Excel files, manual GL, mid-period prorata | One-click depreciation, auto-GL, auto-prorata | 5 days → 45 min |
| UJ-2 | Chief Accountant | No dashboard, paper approvals, context-free requests | FA dashboard, digital approvals with impact analysis | 2 hr → 20 min |
| UJ-3 | CFO | No FA analytics, blind capex decisions | Aging/utilization dashboard, replacement forecast | 3 days → 30 min |
| UJ-4 | Tax Accountant | Manual VAT tracking, two-book reconciliation | Auto VAT report, side-by-side tax/accounting book | 2 days → 2 hr |
| UJ-5 | Internal Auditor | Paper inventory, no real-time tracking | Mobile QR scanning, real-time progress, auto-actions | 7 days → 3 days |
| UJ-6 | External Auditor | Data export reformatting, manual recalc | Audit package export, auto-verification | 4 days → 1 day |
| UJ-7 | System Admin | No FA module, no migration tools | Built-in module, guided setup, import tool | 3 weeks → 3 days |
