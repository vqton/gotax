/* ─── HyperUI-derived component library ─────────────────────── */
/* Pure functions returning HTML strings. No build step. Alpine.js + Tailwind v4. */
/* NOTE: Functions with Alpine directives (filterBar, tabs, dropdown) must be */
/* called in page templates directly, NOT via x-html. x-html does not compile */
/* Alpine directives. Use these in static HTML or set innerHTML before Alpine init. */

/* ─── Shared constants ──────────────────────────────────────── */

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
  OVERPAID: "bg-purple-100 text-purple-700",
  SIGNED: "bg-blue-100 text-blue-700",
  VALIDATED: "bg-indigo-100 text-indigo-700",
  ISSUED: "bg-green-100 text-green-700",
  REPLACED: "bg-purple-100 text-purple-700",
};

const HYPERUI_ICON_COLORS = {
  blue:   { bg: "bg-blue-50",   text: "text-blue-600" },
  green:  { bg: "bg-green-50",  text: "text-green-600" },
  amber:  { bg: "bg-amber-50",  text: "text-amber-600" },
  red:    { bg: "bg-red-50",    text: "text-red-600" },
  purple: { bg: "bg-purple-50", text: "text-purple-600" },
  gray:   { bg: "bg-gray-50",   text: "text-gray-600" },
};

const HyperUI = {

  /* ─── Alerts / Toasts ──────────────────────────────────────── */
  alert(type, message, dismissible = false) {
    const styles = {
      success: "bg-green-50 text-green-800 border-green-200",
      error:   "bg-red-50 text-red-800 border-red-200",
      warning: "bg-amber-50 text-amber-800 border-amber-200",
      info:    "bg-blue-50 text-blue-800 border-blue-200",
    };
    const icons = {
      success: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>',
      error:   '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/>',
      warning: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4.5c-.77-.833-2.694-.833-3.464 0L3.34 16.5c-.77.833.192 2.5 1.732 2.5z"/>',
      info:    '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>',
    };
    const cls = styles[type] || styles.info;
    const icon = icons[type] || icons.info;
    const dismiss = dismissible
      ? `<button onclick="this.parentElement.remove()" class="ml-auto shrink-0 text-current opacity-50 hover:opacity-100">
           <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
         </button>`
      : '';
    return `
    <div class="flex items-start gap-3 rounded-lg border px-4 py-3 text-sm ${cls}" role="alert">
      <svg class="w-5 h-5 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">${icon}</svg>
      <span class="flex-1">${message}</span>
      ${dismiss}
    </div>`;
  },

  /* ─── Stat Cards ───────────────────────────────────────────── */
  statCard(label, value, opts = {}) {
    const { icon, color = "blue", subtext, href } = opts;
    const valClass = {
      blue:   "text-blue-700",
      green:  "text-green-700",
      amber:  "text-amber-700",
      red:    "text-red-700",
      purple: "text-purple-700",
      gray:   "text-gray-700",
    }[color] || "text-blue-700";
    const ic = HYPERUI_ICON_COLORS[color] || HYPERUI_ICON_COLORS.blue;
    const iconHtml = icon
      ? `<div class="w-10 h-10 rounded-lg ${ic.bg} flex items-center justify-center ${ic.text} text-lg">${icon}</div>`
      : '';
    const sub = subtext ? `<p class="text-xs text-gray-400 mt-1">${subtext}</p>` : '';
    const wrap = href ? `<a href="${href}" class="block hover:shadow-md transition-shadow">` : '<div>';
    const close = href ? '</a>' : '</div>';

    return `
    ${wrap}
      <div class="bg-white rounded-xl border border-gray-200 p-5">
        <div class="flex items-start justify-between">
          <div>
            <p class="text-sm text-gray-500 mb-1">${label}</p>
            <p class="text-2xl font-bold ${valClass}">${value}</p>
            ${sub}
          </div>
          ${iconHtml}
        </div>
      </div>
    ${close}`;
  },

  /* ─── Table ────────────────────────────────────────────────── */
  table(columns, rows, opts = {}) {
    const { emptyIcon = "📋", emptyTitle = "Không có dữ liệu", emptySubtitle = "", id } = opts;
    const thead = columns.map(c =>
      `<th class="text-left px-5 py-2.5 font-semibold text-gray-600 ${c.align === 'right' ? 'text-right' : ''}" ${c.width ? `style="width:${c.width}"` : ''}>${c.label}</th>`
    ).join('');

    const tbody = rows.length
      ? rows.map(row => {
          const cells = columns.map(c =>
            `<td class="px-5 py-3 ${c.align === 'right' ? 'text-right' : ''} ${c.class || ''}">${c.render ? c.render(row) : (row[c.key] ?? '—')}</td>`
          ).join('');
          return `<tr class="border-b border-gray-100 hover:bg-gray-50">${cells}</tr>`;
        }).join('')
      : `<tr><td colspan="${columns.length}" class="px-5 py-12 text-center">${HyperUI.emptyState(emptyIcon, emptyTitle, emptySubtitle)}</td></tr>`;

    return `
    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden" ${id ? `id="${id}"` : ''}>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>${thead}</tr>
          </thead>
          <tbody>${tbody}</tbody>
        </table>
      </div>
    </div>`;
  },

  /* ─── Table with header bar ────────────────────────────────── */
  tableWithHeader(title, badgeText, columns, rows, opts = {}) {
    const header = badgeText
      ? `<div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
           <h3 class="text-sm font-semibold text-gray-700">${title}</h3>
           <span class="text-xs text-gray-400">${badgeText}</span>
         </div>`
      : `<div class="px-5 py-4 border-b border-gray-100">
           <h3 class="text-sm font-semibold text-gray-700">${title}</h3>
         </div>`;

    const thead = columns.map(c =>
      `<th class="text-left px-5 py-2.5 font-semibold text-gray-600 ${c.align === 'right' ? 'text-right' : ''}">${c.label}</th>`
    ).join('');

    const tbody = rows.length
      ? rows.map(row => {
          const cells = columns.map(c =>
            `<td class="px-5 py-3 ${c.align === 'right' ? 'text-right' : ''} ${c.class || ''}">${c.render ? c.render(row) : (row[c.key] ?? '—')}</td>`
          ).join('');
          return `<tr class="border-b border-gray-100 hover:bg-gray-50">${cells}</tr>`;
        }).join('')
      : `<tr><td colspan="${columns.length}" class="px-5 py-12 text-center">${HyperUI.emptyState(opts.emptyIcon || "📋", opts.emptyTitle || "Không có dữ liệu", opts.emptySubtitle || "")}</td></tr>`;

    return `
    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      ${header}
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>${thead}</tr>
          </thead>
          <tbody>${tbody}</tbody>
        </table>
      </div>
    </div>`;
  },

  /* ─── Empty State ──────────────────────────────────────────── */
  emptyState(icon, title, subtitle) {
    return `
    <div class="text-center py-12">
      <span class="text-4xl">${icon}</span>
      <h3 class="mt-3 text-sm font-medium text-gray-900">${title}</h3>
      <p class="mt-1 text-sm text-gray-500">${subtitle}</p>
    </div>`;
  },

  /* ─── Loading Skeleton ─────────────────────────────────────── */
  loadingSkeleton(rows = 5, cols = 4) {
    const trs = Array.from({ length: rows }, () => {
      const tds = Array.from({ length: cols }, (_, i) => {
        const w = i === 0 ? 'w-24' : i === cols - 1 ? 'w-16' : 'w-32';
        return `<td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse ${w}"></div></td>`;
      }).join('');
      return `<tr class="border-b border-gray-100">${tds}</tr>`;
    }).join('');
    return `<tbody>${trs}</tbody>`;
  },

  /* ─── Badge ────────────────────────────────────────────────── */
  badge(text, color = "gray") {
    const cls = {
      gray:   "bg-gray-100 text-gray-700",
      green:  "bg-green-100 text-green-700",
      red:    "bg-red-100 text-red-700",
      amber:  "bg-amber-100 text-amber-700",
      blue:   "bg-blue-100 text-blue-700",
      purple: "bg-purple-100 text-purple-700",
      indigo: "bg-indigo-100 text-indigo-700",
    }[color] || "bg-gray-100 text-gray-700";
    return `<span class="inline-block px-2 py-0.5 rounded-full text-xs font-medium ${cls}">${text}</span>`;
  },

  /* ─── Modal ────────────────────────────────────────────────── */
  modal(showVar, title, content, opts = {}) {
    const { size = "md", footer = "" } = opts;
    const sizeMap = { sm: "max-w-md", md: "max-w-2xl", lg: "max-w-4xl", xl: "max-w-6xl" };
    const maxW = sizeMap[size] || sizeMap.md;
    const footerHtml = footer ? `<div class="flex justify-end gap-2 mt-6 pt-4 border-t border-gray-100">${footer}</div>` : '';

    return `
    <template x-if="${showVar}">
      <div class="fixed inset-0 bg-black/40 flex items-center justify-center z-50" @click.self="${showVar} = null">
        <div class="bg-white rounded-xl shadow-xl w-full ${maxW} max-h-[90vh] overflow-y-auto p-6" @click.stop>
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-bold text-gray-900">${title}</h3>
            <button @click="${showVar} = null" class="text-gray-400 hover:text-gray-600">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
            </button>
          </div>
          ${content}
          ${footerHtml}
        </div>
      </div>
    </template>`;
  },

  /* ─── Page Header ──────────────────────────────────────────── */
  pageHeader(title, actions = "") {
    return `
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-bold text-gray-900">${title}</h2>
      <div class="flex items-center gap-2">${actions}</div>
    </div>`;
  },

  /* ─── Filter Bar (static HTML — use in page template) ──────── */
  /* Returns HTML WITHOUT Alpine directives. Page must bind x-model/@click */
  /* in its own template. Example usage:
     <div id="filter-bar"></div>
     <script>document.getElementById('filter-bar').innerHTML = HyperUI.filterBarHTML([...])</script>
     Then in the page's x-data template, add x-model bindings.
     OR: just write the filter bar inline in the page HTML. */
  filterBarHTML(fields) {
    const inputs = fields.map(f => {
      if (f.type === "select") {
        const opts = f.options.map(o => `<option value="${o.value}">${o.label}</option>`).join('');
        return `<div>
          <label class="block text-xs font-medium text-gray-500 mb-1">${f.label}</label>
          <select class="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500">${opts}</select>
        </div>`;
      }
      return `<div>
        <label class="block text-xs font-medium text-gray-500 mb-1">${f.label}</label>
        <input type="${f.type || 'text'}" class="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500" ${f.placeholder ? `placeholder="${f.placeholder}"` : ''}>
      </div>`;
    }).join('');

    return `
    <div class="bg-white rounded-xl border border-gray-200 p-4 mb-6">
      <div class="flex flex-wrap items-end gap-4">
        ${inputs}
        <button class="px-3 py-2 text-sm text-blue-700 hover:bg-blue-50 rounded-lg transition">Lọc</button>
      </div>
    </div>`;
  },

  /* ─── Button ───────────────────────────────────────────────── */
  btn(text, opts = {}) {
    const { variant = "primary", size = "md", icon, onclick, disabled, type = "button" } = opts;
    const variants = {
      primary:  "bg-blue-700 text-white hover:bg-blue-800",
      secondary:"bg-gray-100 text-gray-700 hover:bg-gray-200",
      danger:   "bg-red-600 text-white hover:bg-red-700",
      ghost:    "text-gray-600 hover:text-gray-900 hover:bg-gray-50",
      success:  "bg-green-600 text-white hover:bg-green-700",
    };
    const sizes = {
      sm: "px-3 py-1.5 text-xs",
      md: "px-4 py-2 text-sm",
      lg: "px-5 py-2.5 text-sm",
    };
    const cls = `${variants[variant] || variants.primary} ${sizes[size] || sizes.md} rounded-lg font-medium transition disabled:opacity-50`;
    const iconHtml = icon ? `<span class="mr-1">${icon}</span>` : '';
    return `<button type="${type}" class="${cls}" ${onclick ? `onclick="${onclick}"` : ''} ${disabled ? 'disabled' : ''}>${iconHtml}${text}</button>`;
  },

  /* ─── Pagination ───────────────────────────────────────────── */
  pagination(current, total, perPage, onChange) {
    const totalPages = Math.ceil(total / perPage);
    if (totalPages <= 1) return '';
    const pages = [];
    for (let i = 1; i <= totalPages; i++) {
      if (i === 1 || i === totalPages || (i >= current - 1 && i <= current + 1)) {
        pages.push(i);
      } else if (pages[pages.length - 1] !== '...') {
        pages.push('...');
      }
    }
    const btns = pages.map(p => {
      if (p === '...') return `<span class="px-3 py-1 text-sm text-gray-400">...</span>`;
      const active = p === current ? 'bg-blue-700 text-white' : 'text-gray-700 hover:bg-gray-100';
      return `<button onclick="${onChange}(${p})" class="px-3 py-1 text-sm rounded-lg ${active}">${p}</button>`;
    }).join('');

    return `
    <div class="flex items-center justify-between px-5 py-3 border-t border-gray-100">
      <p class="text-xs text-gray-500">Hiển thị ${(current-1)*perPage+1}-${Math.min(current*perPage, total)} / ${total}</p>
      <div class="flex items-center gap-1">${btns}</div>
    </div>`;
  },

  /* ─── Tabs (static HTML — use in page template) ────────────── */
  /* Returns HTML WITHOUT Alpine directives. Page must add :class binding */
  /* in its own template. Example:
     <div id="tabs"></div>
     <script>document.getElementById('tabs').innerHTML = HyperUI.tabsHTML([...], 'activeTab')</script>
     Then in the x-data template, add :class to each button. */
  tabsHTML(items, activeVar) {
    const btns = items.map(item => {
      return `<button @click="${activeVar} = '${item.value}'"
        class="px-4 py-2 text-sm font-medium rounded-lg transition text-gray-600 hover:bg-gray-100"
        :class="{ 'bg-blue-700 text-white': ${activeVar} === '${item.value}', 'text-gray-600 hover:bg-gray-100': ${activeVar} !== '${item.value}' }">
        ${item.label}
      </button>`;
    }).join('');
    return `<div class="flex items-center gap-1 p-1 bg-gray-100 rounded-lg">${btns}</div>`;
  },

  /* ─── Dropdown (static HTML — use in page template) ────────── */
  /* Returns HTML WITH Alpine directives. Must be in page template, not x-html. */
  dropdownHTML(trigger, items, opts = {}) {
    const { align = "right", id = "dropdown-" + Math.random().toString(36).slice(2, 8) } = opts;
    const alignClass = align === "right" ? "right-0" : "left-0";
    const menuItems = items.map(item => {
      if (item.divider) return '<div class="border-t border-gray-100 my-1"></div>';
      const danger = item.danger ? 'text-red-600 hover:bg-red-50' : 'text-gray-700 hover:bg-gray-50';
      return `<a href="${item.href || '#'}" onclick="${item.onclick || ''}" class="block px-4 py-2 text-sm ${danger}">${item.icon ? `<span class="mr-2">${item.icon}</span>` : ''}${item.label}</a>`;
    }).join('');

    return `
    <div x-data="{ open: false }" class="relative" id="${id}">
      <div @click="open = !open">${trigger}</div>
      <div x-show="open" @click.away="open = false" x-transition
        class="absolute ${alignClass} mt-2 w-56 bg-white rounded-xl shadow-lg border border-gray-200 py-1 z-50">
        ${menuItems}
      </div>
    </div>`;
  },

  /* ─── Confirmation Dialog ──────────────────────────────────── */
  confirmDialog(showVar, message, onConfirm) {
    return `
    <template x-if="${showVar}">
      <div class="fixed inset-0 bg-black/40 flex items-center justify-center z-50" @click.self="${showVar} = false">
        <div class="bg-white rounded-xl shadow-xl w-full max-w-md p-6" @click.stop>
          <div class="flex items-start gap-3 mb-4">
            <div class="w-10 h-10 rounded-full bg-amber-100 flex items-center justify-center shrink-0">
              <svg class="w-5 h-5 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4.5c-.77-.833-2.694-.833-3.464 0L3.34 16.5c-.77.833.192 2.5 1.732 2.5z"/></svg>
            </div>
            <div>
              <h3 class="text-sm font-semibold text-gray-900">Xác nhận</h3>
              <p class="text-sm text-gray-600 mt-1">${message}</p>
            </div>
          </div>
          <div class="flex justify-end gap-2">
            <button @click="${showVar} = false" class="px-4 py-2 text-sm text-gray-600 hover:text-gray-900">Hủy</button>
            <button @click="${onConfirm}; ${showVar} = false" class="px-4 py-2 text-sm bg-red-600 text-white rounded-lg hover:bg-red-700">Xác nhận</button>
          </div>
        </div>
      </div>
    </template>`;
  },
};

/* ─── Backward-compatible wrappers ────────────────────────────── */

function statusBadge(status) {
  const cls = STATUS_COLORS[status] || "bg-gray-100 text-gray-700";
  return `<span class="inline-block px-2 py-0.5 rounded-full text-xs font-medium ${cls}">${status}</span>`;
}

function emptyState(icon, title, subtitle) {
  return HyperUI.emptyState(icon, title, subtitle);
}

function loadingSkeleton(rows = 5) {
  return HyperUI.loadingSkeleton(rows, 4);
}
