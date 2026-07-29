"use strict";

window._alerts = window._alerts || [];

const API_BASE = "";

/* ---------- Alert sink (shared) ---------- */
function mountAlerts(containerId) {
  const host = document.getElementById(containerId);
  if (!host || host._mounted) return;
  host._mounted = true;
  host.innerHTML = Alpine.raw(`
    <template x-for="a in alerts" :key="a.text + a.type">
      <div class="alert" :class="'alert--' + a.type" x-text="a.text" x-transition @click="alerts.shift()" role="alert"></div>
    </template>
  `);
}

/* ---------- Root store: JWT + session ---------- */
function authStore() {
  return {
    accessToken: null,
    refreshToken: null,
    user: null,
    expiresAt: 0,
    init() {
      try {
        const raw = localStorage.getItem("gotax_session");
        if (raw) {
          const s = JSON.parse(raw);
          this.accessToken = s.access_token;
          this.refreshToken = s.refresh_token;
          this.user = s.user;
          this.expiresAt = s.expires_at || 0;
        }
      } catch (_) { this.clear(); }
      this.scheduleRefresh();
      setTimeout(() => scheduleReauth());
    },
    set(payload) {
      this.accessToken = payload.access_token;
      this.refreshToken = payload.refresh_token;
      this.user = payload.user;
      this.expiresAt = (payload.expires_in || 900) * 1000 + Date.now();
      this.persist();
      this.scheduleRefresh();
    },
    persist() {
      localStorage.setItem("gotax_session", JSON.stringify({
        access_token: this.accessToken,
        refresh_token: this.refreshToken,
        user: this.user,
        expires_at: this.expiresAt,
      }));
    },
    clear() {
      this.accessToken = null;
      this.refreshToken = null;
      this.user = null;
      this.expiresAt = 0;
      localStorage.removeItem("gotax_session");
    },
    get isAuthenticated() {
      return !!this.accessToken && Date.now() < this.expiresAt;
    },
    get authHeader() {
      return this.accessToken ? `Bearer ${this.accessToken}` : "";
    },
    async refresh() {
      if (!this.refreshToken) return false;
      try {
        const r = await fetch(`${API_BASE}/api/v1/auth/refresh`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ refresh_token: this.refreshToken }),
        });
        if (!r.ok) throw new Error("refresh_failed");
        const data = await r.json();
        this.set(data);
        return true;
      } catch {
        this.clear();
        return false;
      }
    },
    scheduleRefresh() {
      if (this._timer) clearTimeout(this._timer);
      const left = this.expiresAt - Date.now() - 60_000;
      this._timer = setTimeout(() => this.refresh(), Math.max(left, 0));
    },
  };
}

/* ---------- Toast / notification (global) ---------- */
function notificationStore() {
  return {
    alerts: window._alerts,
    init() { mountAlerts("global-alert-host"); },
    push(type, text) {
      this.alerts.push({ type, text });
      setTimeout(() => this.alerts.shift(), 5000);
    },
    success(msg) { this.push("success", msg); },
    error(msg)   { this.push("error",   msg); },
    info(msg)    { this.push("info",    msg); },
    warning(msg) { this.push("warning", msg); },
  };
}

/* ---------- Field validation ---------- */
function validateEmail(v) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v);
}
function validateConfirm(pw, pw2) {
  return pw === pw2 && pw.length >= 8;
}

/* ---------- Reauth gate (idle detector) ---------- */
let inactiveTimer = null;
function scheduleReauth() {
  clearTimeout(inactiveTimer);
  inactiveTimer = setTimeout(() => {
    const store = Alpine.store("auth");
    if (store && store.isAuthenticated) {
      store.refresh().then(ok => { if (!ok) window.location = "/login"; });
    }
  }, 1000 * 60 * 30);
}
["mousemove", "keydown", "click", "scroll", "touchstart"].forEach(evt =>
  document.addEventListener(evt, scheduleReauth, { passive: true })
);

/* ---------- htmx: populate auth header ---------- */
document.addEventListener("htmx:configRequest", (e) => {
  const store = window.__alpine_store;
  if (store?.accessToken) {
    e.detail.headers["Authorization"] = `Bearer ${store.accessToken}`;
  }
});

/* ---------- Boot ---------- */
document.addEventListener("alpine:init", () => {
  Alpine.store("auth", authStore());
  Alpine.store("alert", notificationStore());
});
