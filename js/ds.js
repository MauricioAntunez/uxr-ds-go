// UXR Design System — ES2025+ Vanilla JS
// No HTMX. No frameworks. Progressive enhancement only.
const DS = (() => {
  'use strict';

  // --- CSRF injection ---
  function injectCSRF(token) {
    for (const form of document.forms) {
      if (!form.querySelector('input[name="_csrf"]')) {
        const input = Object.assign(document.createElement('input'), {
          type: 'hidden', name: '_csrf', value: token,
        });
        form.append(input);
      }
    }
  }

  // --- Tabs (delegated, aria-managed) ---
  function initTabs() {
    document.addEventListener('click', (e) => {
      const tab = e.target.closest('.tab');
      if (!tab) return;
      const tabList = tab.closest('.tab-list');
      const tabs = tab.closest('.tabs');
      if (!tabList || !tabs) return;
      for (const t of tabList.querySelectorAll('.tab')) {
        t.setAttribute('aria-selected', 'false');
      }
      tab.setAttribute('aria-selected', 'true');
      for (const panel of tabs.querySelectorAll('.tab-panel')) {
        panel.hidden = panel.id !== `panel-${tab.dataset.tab}`;
      }
    });
  }

  // --- Flash auto-dismiss ---
  function initFlashDismiss(delay = 5000) {
    for (const alert of document.querySelectorAll('.alert')) {
      setTimeout(() => {
        alert.style.opacity = '0';
        alert.style.transition = 'opacity 0.5s';
        setTimeout(() => alert.remove(), 500);
      }, delay);
    }
  }

  // --- Native <dialog> confirm pattern ---
  function initConfirmDialogs() {
    document.addEventListener('click', (e) => {
      const trigger = e.target.closest('[data-confirm-dialog]');
      if (!trigger) return;
      e.preventDefault();
      const dialog = document.getElementById(trigger.dataset.confirmDialog);
      if (!dialog) return;

      const { promise, resolve } = Promise.withResolvers();
      dialog.addEventListener('close', () => resolve(dialog.returnValue), { once: true });
      dialog.showModal();

      promise.then(value => {
        if (value === 'confirm') {
          trigger.closest('form')?.requestSubmit()
            ?? (window.location.href = trigger.href);
        }
      });
    });
  }

  // --- Expandable table rows ---
  function initExpandableRows() {
    document.addEventListener('click', (e) => {
      const expandBtn = e.target.closest('.expand-btn');
      if (!expandBtn) return;

      const row = expandBtn.closest('tr');
      const detailRow = row?.nextElementSibling;
      if (!detailRow?.classList.contains('detail-row')) return;

      const isExpanded = expandBtn.getAttribute('aria-expanded') === 'true';
      expandBtn.setAttribute('aria-expanded', !isExpanded);
      expandBtn.textContent = isExpanded ? '+' : '-';
      detailRow.hidden = isExpanded;
      row.classList.toggle('expanded', !isExpanded);
    });
  }

  // --- Polling helper with AbortSignal ---
  function createPoller(url, interval, callback) {
    let controller = new AbortController();
    let timer = null;

    async function poll() {
      try {
        const signal = AbortSignal.any([
          controller.signal,
          AbortSignal.timeout(10_000),
        ]);
        const resp = await fetch(url, { signal });
        if (!resp.ok) return;
        const data = await resp.json();
        callback(data);
      } catch (e) {
        if (e.name !== 'AbortError') console.error('Poll failed:', e);
      }
    }

    return {
      start() {
        if (timer) return;
        poll();
        timer = setInterval(poll, interval);
        document.addEventListener('visibilitychange', () => {
          if (document.hidden) {
            clearInterval(timer);
            timer = null;
          } else {
            poll();
            timer = setInterval(poll, interval);
          }
        });
      },
      stop() {
        clearInterval(timer);
        timer = null;
        controller.abort();
        controller = new AbortController();
      },
    };
  }

  // --- Intl formatters (locale-aware) ---
  const fmt = {
    number: (n) => new Intl.NumberFormat().format(n ?? 0),
    cost: (n) => new Intl.NumberFormat(undefined, {
      style: 'currency', currency: 'USD', minimumFractionDigits: 4,
    }).format(n ?? 0),
    relativeTime: (() => {
      const rtf = new Intl.RelativeTimeFormat(undefined, { numeric: 'auto' });
      return (date) => {
        const diff = (new Date(date) - Date.now()) / 1000;
        if (Math.abs(diff) < 60) return rtf.format(Math.round(diff), 'second');
        if (Math.abs(diff) < 3600) return rtf.format(Math.round(diff / 60), 'minute');
        if (Math.abs(diff) < 86400) return rtf.format(Math.round(diff / 3600), 'hour');
        return rtf.format(Math.round(diff / 86400), 'day');
      };
    })(),
  };

  // --- Set operations for permission checking (ES2025) ---
  function hasAnyPermission(userPerms, requiredPerms) {
    return new Set(userPerms).intersection(new Set(requiredPerms)).size > 0;
  }

  // --- Object.groupBy for table data grouping (ES2024) ---
  function groupBy(items, key) {
    return Object.groupBy(items, item => item[key]);
  }

  return {
    injectCSRF, initTabs, initFlashDismiss, initConfirmDialogs,
    initExpandableRows, createPoller, fmt, hasAnyPermission, groupBy,
  };
})();
