"use strict";

const API = "/api/v1/payroll";

/* ─── Nav layout (single source of truth) ──────────────────── */

const NAV_ITEMS = [
  { href: "/payroll/periods.html",     label: "Kỳ lương",   icon: "📅" },
  { href: "/payroll/timekeeping.html", label: "Chấm công",  icon: "⏰" },
  { href: "/payroll/leave.html",       label: "Nghỉ phép",  icon: "🏖" },
  { href: "/payroll/payslips.html",    label: "Phiếu lương", icon: "💰" },
  { href: "/payroll/declarations.html",label: "Tuyên thuế",  icon: "📋" },
  { href: "/payroll/config.html",      label: "Cấu hình",   icon: "⚙️" },
];

function pwNavHTML(activePath) {
  const links = NAV_ITEMS.map(item => {
    const active = location.pathname === item.href || location.pathname === activePath;
    return `<a href="${item.href}" class="px-3 py-1.5 rounded-lg text-sm transition ${active ? 'bg-blue-50 text-blue-700 font-semibold' : 'text-gray-600 hover:text-blue-700 hover:bg-gray-50'}">${item.icon} ${item.label}</a>`;
  }).join("\n        ");
  return `
  <div class="flex items-center gap-2">
    <a href="/payroll/periods.html" class="text-lg font-bold text-blue-700 mr-2">📊 GoTax</a>
    <span class="text-gray-300 mr-2">|</span>
    ${links}
  </div>
  <a href="/login" onclick="Alpine.store('auth').clear()" class="text-sm text-gray-500 hover:text-gray-700">Đăng xuất</a>`;
}

function mountNav(containerId, activePath) {
  const el = document.getElementById(containerId);
  if (el) el.innerHTML = pwNavHTML(activePath);
}

/* ─── Flash messages ───────────────────────────────────────── */

function pwFlashHTML() {
  return `
  <template x-if="$store.pw.success">
    <div class="mb-4 rounded-md bg-green-50 border border-green-200 px-4 py-3 text-sm text-green-800" x-text="$store.pw.success"></div>
  </template>
  <template x-if="$store.pw.error">
    <div class="mb-4 rounded-md bg-red-50 border border-red-200 px-4 py-3 text-sm text-red-800" x-text="$store.pw.error"></div>
  </template>`;
}

/* ─── Fetch wrapper ────────────────────────────────────────── */

async function pwFetch(path, opts = {}) {
  const store = Alpine.store("auth");
  const headers = { "Content-Type": "application/json", ...opts.headers };
  if (store?.accessToken) headers["Authorization"] = `Bearer ${store.accessToken}`;
  const r = await fetch(`${API}${path}`, { ...opts, headers });
  if (r.status === 401) {
    const ok = await store.refresh();
    if (!ok) { window.location = "/login"; throw new Error("unauthorized"); }
    headers["Authorization"] = store.authHeader;
    return fetch(`${API}${path}`, { ...opts, headers });
  }
  return r;
}

async function pwGet(path) {
  const r = await pwFetch(path);
  if (!r.ok) throw await r.json().catch(() => ({ error: `HTTP ${r.status}` }));
  return r.json();
}

async function pwPost(path, body) {
  const r = await pwFetch(path, { method: "POST", body: JSON.stringify(body) });
  if (!r.ok) throw await r.json().catch(() => ({ error: `HTTP ${r.status}` }));
  return r.json();
}

async function pwPut(path, body) {
  const r = await pwFetch(path, { method: "PUT", body: JSON.stringify(body) });
  if (!r.ok) throw await r.json().catch(() => ({ error: `HTTP ${r.status}` }));
  return r.json();
}

async function pwDelete(path) {
  const r = await pwFetch(path, { method: "DELETE" });
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
  PROCESSING: "bg-blue-100 text-blue-700",
  SUBMITTED: "bg-amber-100 text-amber-700",
  APPROVED: "bg-green-100 text-green-700",
  PAID: "bg-purple-100 text-purple-700",
  CLOSED: "bg-gray-200 text-gray-500",
  PENDING: "bg-amber-100 text-amber-700",
  REJECTED: "bg-red-100 text-red-700",
  CALCULATED: "bg-blue-100 text-blue-700",
};

function statusClass(status) {
  return STATUS_COLORS[status] || "bg-gray-100 text-gray-700";
}

/* ─── Alpine store: payroll (shared state) ─────────────────── */

document.addEventListener("alpine:init", () => {
  Alpine.store("pw", {
    loading: false,
    error: "",
    success: "",

    async load(fn) {
      this.loading = true;
      this.error = "";
      try {
        return await fn();
      } catch (e) {
        this.error = e.error || e.message || "Lỗi không xác định";
        throw e;
      } finally {
        this.loading = false;
      }
    },

    flashSuccess(msg) {
      this.success = msg;
      setTimeout(() => (this.success = ""), 4000);
    },
  });
});
