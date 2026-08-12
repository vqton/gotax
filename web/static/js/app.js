// GoTax shared frontend — plain JS, no Alpine. Exposes window.GoTax.
(function () {
  "use strict";

  var SESSION_KEY = "gotax_session"; // { access_token, refresh_token, user, expires_at }

  var GoTax = (window.GoTax = window.GoTax || {});

  /* ─── Auth: JWT session in localStorage + cookie for page loads ─── */

  function readSession() {
    try {
      return JSON.parse(localStorage.getItem(SESSION_KEY) || "null");
    } catch (_) {
      return null;
    }
  }

  GoTax.Auth = {
    getSession: readSession,
    getAccessToken: function () {
      var s = readSession();
      return s ? s.access_token : null;
    },
    getRefreshToken: function () {
      var s = readSession();
      return s ? s.refresh_token : null;
    },
    getUser: function () {
      var s = readSession();
      return s && s.user ? s.user : null;
    },
    isAuthenticated: function () {
      var s = readSession();
      return !!(s && s.access_token && Date.now() < (s.expires_at || 0));
    },
    saveLogin: function (p) {
      var s = {
        access_token: p.access_token,
        refresh_token: p.refresh_token,
        user: p.user,
        expires_at: (p.expires_in || 900) * 1000 + Date.now(),
      };
      localStorage.setItem(SESSION_KEY, JSON.stringify(s));
      this.setCookie(s.access_token);
      this.scheduleRefresh();
    },
    clear: function () {
      localStorage.removeItem(SESSION_KEY);
      document.cookie = "gotax_token=; path=/; max-age=0; SameSite=Lax";
    },
    setCookie: function (token) {
      document.cookie = "gotax_token=" + encodeURIComponent(token) + "; path=/; max-age=900; SameSite=Lax";
    },
    refresh: function () {
      var rt = this.getRefreshToken();
      if (!rt) return Promise.resolve(false);
      return fetch("/api/v1/auth/refresh", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refresh_token: rt }),
      })
        .then(function (r) {
          if (!r.ok) throw new Error("refresh_failed");
          return r.json();
        })
        .then(function (d) {
          this.saveLogin(d);
          return true;
        }.bind(this))
        .catch(function () {
          this.clear();
          return false;
        }.bind(this));
    },
    scheduleRefresh: function () {
      var self = this;
      if (self._timer) clearTimeout(self._timer);
      var s = readSession();
      if (!s) return;
      var left = (s.expires_at || 0) - Date.now() - 60000;
      self._timer = setTimeout(function () { self.refresh(); }, Math.max(left, 0));
    },
    logout: function () {
      this.clear();
      window.location = "/login";
    },
  };

  /* ─── Toast ─── */

  GoTax.Toast = {
    show: function (text, type) {
      var host = document.getElementById("toast-host");
      if (!host) {
        host = document.createElement("div");
        host.id = "toast-host";
        host.className = "toast-host";
        document.body.appendChild(host);
      }
      var el = document.createElement("div");
      el.className = "toast toast-" + (type || "info");
      el.textContent = text;
      el.onclick = function () { el.remove(); };
      host.appendChild(el);
      setTimeout(function () { el.remove(); }, 5000);
    },
  };

  /* ─── htmx hooks ─── */

  document.addEventListener("htmx:configRequest", function (e) {
    var t = GoTax.Auth.getAccessToken();
    if (t) e.detail.headers["Authorization"] = "Bearer " + t;
  });

  document.addEventListener("htmx:beforeSwap", function (e) {
    if (e.detail.xhr.status === 401) {
      e.detail.shouldSwap = false;
      GoTax.Auth.clear();
      window.location = "/login";
    }
  });

  // HX-Trigger JSON from server: {"toast": {"type": "success", "text": "..."}}
  document.addEventListener("toast", function (e) {
    var t = e.detail;
    if (t && t.text) GoTax.Toast.show(t.text, t.type);
  });

  // Flowbite auto-inits once at load; re-init after htmx swaps bring in
  // new data-* elements (dropdowns, modals, collapses in fragments).
  document.addEventListener("htmx:afterSwap", function () {
    if (window.initFlowbite) window.initFlowbite();
  });

  /* ─── Formatters ─── */

  GoTax.fmtVND = function (v) {
    if (v == null || isNaN(v)) return "0";
    return new Intl.NumberFormat("vi-VN").format(Math.round(v));
  };
  GoTax.fmtDate = function (s) {
    if (!s) return "—";
    return new Date(s).toLocaleDateString("vi-VN");
  };
  GoTax.fmtDateTime = function (s) {
    if (!s) return "—";
    return new Date(s).toLocaleString("vi-VN");
  };

  /* ─── App helpers ─── */

  GoTax.App = {
    logout: function () { GoTax.Auth.logout(); },
    companyId: function () {
      var u = GoTax.Auth.getUser();
      return (u && u.company_id) || "CMP001";
    },
    confirm: function (msg) { return window.confirm(msg); },
  };
})();
