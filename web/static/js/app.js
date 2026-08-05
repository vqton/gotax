"use strict";

/* ─── API base ─────────────────────────────────────────────── */
const API_BASE = "/api/v1";

/* ─── Sidebar navigation (MISA SME layout) ─────────────────── */

const MODULE_GROUPS = [
  {
    label: "Tổng quan",
    items: [
      { href: "/app/dashboard.html", icon: "📊", label: "Dashboard" },
    ]
  },
  {
    label: "Hệ thống",
    items: [
      { href: "/app/company.html", icon: "🏢", label: "Công ty" },
      { href: "/app/users.html", icon: "👥", label: "Người dùng" },
      { href: "/app/audit.html", icon: "📜", label: "Nhật ký" },
    ]
  },
  {
    label: "Danh mục",
    items: [
      { href: "/app/coa.html", icon: "📋", label: "Hệ thống tài khoản" },
      { href: "/app/customers.html", icon: "🧑", label: "Khách hàng" },
      { href: "/app/suppliers.html", icon: "🏭", label: "Nhà cung cấp" },
      { href: "/app/items.html", icon: "📦", label: "Hàng hóa" },
      { href: "/app/warehouses.html", icon: "🏬", label: "Kho" },
    ]
  },
  {
    label: "Kế toán",
    items: [
      { href: "/app/journal-entries.html", icon: "📝", label: "Chứng từ" },
      { href: "/app/periods.html", icon: "📅", label: "Kỳ kế toán" },
      { href: "/app/exchange-rates.html", icon: "💱", label: "Tỷ giá" },
      { href: "/app/opening-balances.html", icon: "⚖️", label: "Số dư đầu kỳ" },
    ]
  },
  {
    label: "Mua hàng",
    items: [
      { href: "/app/purchase-orders.html", icon: "🛒", label: "Đơn đặt hàng" },
      { href: "/app/grn.html", icon: "📥", label: "Nhập kho" },
      { href: "/app/purchase-invoices.html", icon: "🧾", label: "Hóa đơn mua" },
      { href: "/app/ap-aging.html", icon: "📊", label: "Công nợ phải trả" },
    ]
  },
  {
    label: "Bán hàng",
    items: [
      { href: "/app/sales-quotations.html", icon: "💼", label: "Báo giá" },
      { href: "/app/sales-orders.html", icon: "📋", label: "Đơn hàng" },
      { href: "/app/delivery-notes.html", icon: "🚚", label: "Xuất kho" },
      { href: "/app/sales-invoices.html", icon: "🧾", label: "Hóa đơn bán" },
      { href: "/app/ar-aging.html", icon: "📊", label: "Công nợ phải thu" },
    ]
  },
  {
    label: "Quỹ / Ngân hàng",
    items: [
      { href: "/app/cash-receipts.html", icon: "💵", label: "Phiếu thu" },
      { href: "/app/cash-payments.html", icon: "💸", label: "Phiếu chi" },
      { href: "/app/cash-transfers.html", icon: "🔄", label: "Chuyển quỹ" },
      { href: "/app/bank-statements.html", icon: "🏦", label: "Sổ ngân hàng" },
      { href: "/app/payment-orders.html", icon: "📄", label: "Uy nhiệm chi" },
    ]
  },
  {
    label: "Kho",
    items: [
      { href: "/app/stock-balances.html", icon: "📊", label: "Tồn kho" },
      { href: "/app/stock-transfers.html", icon: "🔄", label: "Chuyển kho" },
      { href: "/app/stock-adjustments.html", icon: "📝", label: "Điều chỉnh" },
    ]
  },
  {
    label: "Tài sản / CCDC",
    items: [
      { href: "/app/fixed-assets.html", icon: "🏗️", label: "Tài sản cố định" },
      { href: "/app/fa-categories.html", icon: "📁", label: "Nhóm tài sản" },
      { href: "/app/depreciation.html", icon: "📉", label: "Khấu hao" },
    ]
  },
  {
    label: "Tiền lương",
    items: [
      { href: "/payroll/periods.html", icon: "📅", label: "Kỳ lương" },
      { href: "/payroll/timekeeping.html", icon: "⏰", label: "Chấm công" },
      { href: "/payroll/payslips.html", icon: "💰", label: "Phiếu lương" },
      { href: "/payroll/declarations.html", icon: "📋", label: "Tuyên thuế" },
      { href: "/payroll/config.html", icon: "⚙️", label: "Cấu hình" },
    ]
  },
  {
    label: "Thuế",
    items: [
      { href: "/app/tax-declarations.html", icon: "📋", label: "Tờ khai" },
      { href: "/app/tax-calendar.html", icon: "📅", label: "Lịch thuế" },
      { href: "/app/vat-report.html", icon: "📊", label: "Báo cáo GTGT" },
    ]
  },
  {
    label: "Báo cáo",
    items: [
      { href: "/app/trial-balance.html", icon: "📊", label: "Bảng cân đối phát sinh" },
      { href: "/app/balance-sheet.html", icon: "📑", label: "Bảng cân đối kế toán" },
      { href: "/app/income-statement.html", icon: "📈", label: "Kết quả kinh doanh" },
      { href: "/app/cash-flow.html", icon: "💰", label: "Lưu chuyển tiền tệ" },
    ]
  },
];

/* ─── Sidebar HTML ─────────────────────────────────────────── */

function sidebarHTML(activePath) {
  const groups = MODULE_GROUPS.map(g => {
    const items = g.items.map(item => {
      const active = location.pathname === item.href || activePath === item.href;
      return `<a href="${item.href}" class="flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm transition
        ${active ? 'bg-blue-50 text-blue-700 font-semibold' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}">
        <span class="text-base">${item.icon}</span>
        <span>${item.label}</span>
      </a>`;
    }).join("");
    return `<div class="mb-1">
      <p class="px-3 py-1.5 text-xs font-semibold text-gray-400 uppercase tracking-wider">${g.label}</p>
      <div class="space-y-0.5">${items}</div>
    </div>`;
  }).join("");

  return `
  <div class="flex flex-col h-full">
    <div class="px-4 py-4 border-b border-gray-200">
      <a href="/app/dashboard.html" class="flex items-center gap-2">
        <span class="text-xl">📊</span>
        <span class="text-lg font-bold text-blue-700">GoTax</span>
      </a>
      <p class="text-xs text-gray-400 mt-1">Circular 99/2025</p>
    </div>
    <nav class="flex-1 overflow-y-auto px-3 py-4 space-y-1">
      ${groups}
    </nav>
    <div class="px-4 py-3 border-t border-gray-200">
      <div class="flex items-center gap-2">
        <div class="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center text-sm font-semibold text-blue-700" id="user-avatar">U</div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-gray-900 truncate" id="user-name">User</p>
          <p class="text-xs text-gray-500" id="user-role">—</p>
        </div>
        <button onclick="App.logout()" class="text-gray-400 hover:text-gray-600" title="Đăng xuất">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/></svg>
        </button>
      </div>
    </div>
  </div>`;
}

/* ─── Top bar HTML ─────────────────────────────────────────── */

function topbarHTML(title) {
  return `
  <div class="flex items-center justify-between h-14 px-6 bg-white border-b border-gray-200">
    <div class="flex items-center gap-3">
      <button onclick="document.getElementById('sidebar').classList.toggle('-translate-x-full')" class="lg:hidden text-gray-500 hover:text-gray-700">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/></svg>
      </button>
      <h1 class="text-lg font-semibold text-gray-900">${title}</h1>
    </div>
    <div class="flex items-center gap-4">
      <div class="relative">
        <input type="text" placeholder="Tìm kiếm..." class="w-64 pl-9 pr-3 py-1.5 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent">
        <svg class="absolute left-2.5 top-2 w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
      </div>
      <span class="text-sm text-gray-500" id="company-name"></span>
    </div>
  </div>`;
}

/* ─── App shell mount ──────────────────────────────────────── */

function mountAppShell(title, activePath) {
  const sidebar = document.getElementById("sidebar");
  const topbar = document.getElementById("topbar");
  if (sidebar) sidebar.innerHTML = sidebarHTML(activePath);
  if (topbar) topbar.innerHTML = topbarHTML(title);

  // Load user info
  try {
    const store = Alpine.store("auth");
    if (store?.user) {
      const nameEl = document.getElementById("user-name");
      const roleEl = document.getElementById("user-role");
      const avatarEl = document.getElementById("user-avatar");
      if (nameEl) nameEl.textContent = store.user.username || "User";
      if (roleEl) roleEl.textContent = store.user.role || "—";
      if (avatarEl) avatarEl.textContent = (store.user.username || "U")[0].toUpperCase();
    }
  } catch(e) {}
}

/* ─── API client (with JWT refresh) ────────────────────────── */

async function apiFetch(path, opts = {}) {
  const store = Alpine.store("auth");
  const headers = { "Content-Type": "application/json", ...opts.headers };
  if (store?.accessToken) headers["Authorization"] = `Bearer ${store.accessToken}`;
  const r = await fetch(`${API_BASE}${path}`, { ...opts, headers });
  if (r.status === 401) {
    const ok = await store.refresh();
    if (!ok) { window.location = "/login"; throw new Error("unauthorized"); }
    headers["Authorization"] = store.authHeader;
    return fetch(`${API_BASE}${path}`, { ...opts, headers });
  }
  return r;
}

async function apiGet(path) {
  const r = await apiFetch(path);
  if (!r.ok) throw await r.json().catch(() => ({ error: `HTTP ${r.status}` }));
  return r.json();
}

async function apiPost(path, body) {
  const r = await apiFetch(path, { method: "POST", body: JSON.stringify(body) });
  if (!r.ok) throw await r.json().catch(() => ({ error: `HTTP ${r.status}` }));
  return r.json();
}

async function apiPut(path, body) {
  const r = await apiFetch(path, { method: "PUT", body: JSON.stringify(body) });
  if (!r.ok) throw await r.json().catch(() => ({ error: `HTTP ${r.status}` }));
  return r.json();
}

async function apiDelete(path) {
  const r = await apiFetch(path, { method: "DELETE" });
  if (!r.ok) throw await r.json().catch(() => ({ error: `HTTP ${r.status}` }));
  return r.ok;
}

/* ─── Formatters ───────────────────────────────────────────── */

function fmtVND(v) {
  if (v == null || isNaN(v)) return "0";
  return new Intl.NumberFormat("vi-VN").format(Math.round(v));
}

function fmtDate(s) {
  if (!s) return "—";
  return new Date(s).toLocaleDateString("vi-VN");
}

function fmtDateTime(s) {
  if (!s) return "—";
  return new Date(s).toLocaleString("vi-VN");
}

/* ─── Status helpers ───────────────────────────────────────── */

const STATUS_COLORS = {
  DRAFT: "bg-gray-100 text-gray-700",
  ACTIVE: "bg-green-100 text-green-700",
  INACTIVE: "bg-red-100 text-red-700",
  OPEN: "bg-green-100 text-green-700",
  CLOSED: "bg-gray-200 text-gray-500",
  LOCKED: "bg-red-100 text-red-700",
  POSTED: "bg-blue-100 text-blue-700",
  CANCELLED: "bg-red-100 text-red-700",
  APPROVED: "bg-green-100 text-green-700",
  PENDING: "bg-amber-100 text-amber-700",
  REJECTED: "bg-red-100 text-red-700",
  PAID: "bg-green-100 text-green-700",
  UNPAID: "bg-amber-100 text-amber-700",
  PARTIAL: "bg-blue-100 text-blue-700",
  SUBMITTED: "bg-amber-100 text-amber-700",
  PROCESSING: "bg-blue-100 text-blue-700",
  COMPLETED: "bg-green-100 text-green-700",
  CONFIRMED: "bg-green-100 text-green-700",
  RECEIVED: "bg-blue-100 text-blue-700",
  SENT: "bg-blue-100 text-blue-700",
  ACCEPTED: "bg-green-100 text-green-700",
  EXPIRED: "bg-gray-200 text-gray-500",
  DELIVERED: "bg-blue-100 text-blue-700",
  INVOICED: "bg-amber-100 text-amber-700",
};

function statusBadge(status) {
  const cls = STATUS_COLORS[status] || "bg-gray-100 text-gray-700";
  return `<span class="inline-block px-2 py-0.5 rounded-full text-xs font-medium ${cls}">${status}</span>`;
}

/* ─── Empty state ──────────────────────────────────────────── */

function emptyState(icon, title, subtitle) {
  return `
  <div class="text-center py-12">
    <span class="text-4xl">${icon}</span>
    <h3 class="mt-3 text-sm font-medium text-gray-900">${title}</h3>
    <p class="mt-1 text-sm text-gray-500">${subtitle}</p>
  </div>`;
}

/* ─── Loading skeleton ─────────────────────────────────────── */

function loadingSkeleton(rows = 5) {
  const trs = Array.from({ length: rows }, () =>
    `<tr class="border-b border-gray-100">
      <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-24"></div></td>
      <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-32"></div></td>
      <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-20"></div></td>
      <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-16"></div></td>
    </tr>`
  ).join("");
  return `<tbody>${trs}</tbody>`;
}

/* ─── Alpine store ─────────────────────────────────────────── */

document.addEventListener("alpine:init", () => {
  Alpine.store("app", {
    loading: false,
    error: "",
    success: "",

    flashSuccess(msg) {
      this.success = msg;
      setTimeout(() => (this.success = ""), 4000);
    },

    flashError(msg) {
      this.error = msg;
      setTimeout(() => (this.error = ""), 5000);
    },
  });
});

/* ─── App utilities ────────────────────────────────────────── */

const App = {
  logout() {
    try { Alpine.store("auth").clear(); } catch(e) {}
    window.location = "/login";
  },

  companyId() {
    try { return Alpine.store("auth")?.user?.company_id || "CMP001"; } catch(e) { return "CMP001"; }
  },

  confirm(msg) {
    return window.confirm(msg);
  },
};
