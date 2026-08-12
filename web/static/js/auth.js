// GoTax auth pages — plain JS + fetch. No Alpine, no framework init race.
(function () {
  "use strict";

  var GoTax = window.GoTax = window.GoTax || {};

  function errText(d, status) {
    return (d && d.error) || (d && d.message) || "Server trả lỗi " + status;
  }

  /* ─── Login ─── */

  GoTax.loginPage = function () {
    var form = document.getElementById("login-form");
    if (!form) return;
    var btn = document.getElementById("submit-btn");
    var errEl = document.getElementById("login-error");
    var btnText = btn ? document.getElementById("submit-label") : null;

    form.addEventListener("submit", function (ev) {
      ev.preventDefault();
      var username = document.getElementById("username").value.trim();
      var password = document.getElementById("password").value;
      if (!username || !password) {
        showErr("Nhập đầy đủ tên đăng nhập và mật khẩu.");
        return;
      }
      var body = { username: username, password: password };
      var mst = document.getElementById("mst");
      if (mst && mst.value.trim()) body.tenant_id = mst.value.trim();

      setBusy(true);
      fetch("/api/v1/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      })
        .then(function (r) {
          return r.json().catch(function () { return {}; }).then(function (d) { return { ok: r.ok, status: r.status, d: d }; });
        })
        .then(function (res) {
          setBusy(false);
          if (res.ok && res.d.access_token) {
            GoTax.Auth.saveLogin(res.d);
            window.location.href = "/app/dashboard.html";
            return;
          }
          if (res.status === 401 && res.d.requires_2fa) {
            localStorage.setItem("gotax_temp_token", res.d.temp_token || "");
            window.location.href = "/2fa";
            return;
          }
          showErr(errText(res.d, res.status));
        })
        .catch(function () {
          setBusy(false);
          showErr("Không kết nối được server. Kiểm tra server đang chạy.");
        });
    });

    function showErr(msg) {
      if (errEl) { errEl.textContent = msg; errEl.classList.remove("hidden"); }
      else { GoTax.Toast.show(msg, "error"); }
    }
    function setBusy(b) {
      if (btn) {
        btn.disabled = b;
        if (btnText) btnText.textContent = b ? "Đang đăng nhập…" : "ĐĂNG NHẬP";
      }
    }
  };

  /* ─── Forgot password ─── */

  GoTax.forgotPage = function () {
    var form = document.getElementById("forgot-form");
    if (!form) return;
    var email = document.getElementById("email");
    var msgEl = document.getElementById("form-msg");
    var btn = document.getElementById("submit-btn");
    var btnText = btn ? document.getElementById("submit-label") : null;

    form.addEventListener("submit", function (ev) {
      ev.preventDefault();
      var v = email ? email.value.trim() : "";
      if (!v) { setMsg("Nhập email.", "error"); return; }
      if (btn) btn.disabled = true;
      if (btnText) btnText.textContent = "Đang gửi…";
      fetch("/api/v1/auth/forgot-password", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: v }),
      })
        .then(function (r) { return r.json().catch(function () { return {}; }); })
        .then(function (d) {
          if (btn) btn.disabled = false;
          if (btnText) btnText.textContent = "GỬI YÊU CẦU";
          if (d.ok || d.message || !d.error) setMsg("Nếu email tồn tại, bạn sẽ nhận được liên kết đặt lại mật khẩu.", "success");
          else setMsg(errText(d, 400), "error");
        })
        .catch(function () {
          if (btn) btn.disabled = false;
          if (btnText) btnText.textContent = "GỬI YÊU CẦU";
          setMsg("Không kết nối được server.", "error");
        });
    });

    function setMsg(text, type) {
      if (!msgEl) { GoTax.Toast.show(text, type); return; }
      msgEl.textContent = text;
      msgEl.className = "alert " + (type === "success" ? "alert-success" : "alert-error");
      msgEl.classList.remove("hidden");
    }
  };

  /* ─── Reset password ─── */

  GoTax.resetPage = function () {
    var form = document.getElementById("reset-form");
    if (!form) return;
    var pw = document.getElementById("password");
    var pw2 = document.getElementById("password2");
    var errEl = document.getElementById("login-error");
    var btn = document.getElementById("submit-btn");
    var btnText = btn ? document.getElementById("submit-label") : null;
    var token = new URLSearchParams(window.location.search).get("token") || "";

    form.addEventListener("submit", function (ev) {
      ev.preventDefault();
      var p = pw ? pw.value : "";
      var p2 = pw2 ? pw2.value : "";
      if (p.length < 8) { showErr("Mật khẩu tối thiểu 8 ký tự."); return; }
      if (p !== p2) { showErr("Mật khẩu không khớp."); return; }
      if (btn) btn.disabled = true;
      if (btnText) btnText.textContent = "Đang lưu…";
      fetch("/api/v1/auth/reset-password", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ token: token, new_password: p }),
      })
        .then(function (r) { return r.json().catch(function () { return {}; }); })
        .then(function (d) {
          if (d.ok || d.message || !d.error) {
            window.location.href = "/login?reset=1";
          } else {
            if (btn) btn.disabled = false;
            if (btnText) btnText.textContent = "ĐẶT LẠI MẬT KHẨU";
            showErr(errText(d, 400));
          }
        })
        .catch(function () {
          if (btn) btn.disabled = false;
          if (btnText) btnText.textContent = "ĐẶT LẠI MẬT KHẨU";
          showErr("Không kết nối được server.");
        });
    });

    function showErr(msg) {
      if (errEl) { errEl.textContent = msg; errEl.classList.remove("hidden"); }
      else GoTax.Toast.show(msg, "error");
    }
  };

  /* ─── Password visibility toggles ─── */

  GoTax.togglePw = function (id) {
    var el = document.getElementById(id);
    if (!el) return;
    el.type = el.type === "password" ? "text" : "password";
  };

  document.addEventListener("DOMContentLoaded", function () {
    GoTax.loginPage();
    GoTax.forgotPage();
    GoTax.resetPage();
    // flash after password reset → /login?reset=1
    if (window.location.search.indexOf("reset=1") !== -1) {
      var errEl = document.getElementById("login-error");
      if (errEl) {
        errEl.textContent = "Mật khẩu đã được đặt lại. Đăng nhập bằng mật khẩu mới.";
        errEl.className = "alert alert-success";
      }
    }
  });
})();
