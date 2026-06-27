import "./style.css";
import { calculate, type Variables } from "./api";

// ---------------------------------------------------------------------------
// Keypad definition. Each key has a label, the text it inserts, and a style
// class controlling its color/role in the corporate palette.
// ---------------------------------------------------------------------------
type KeyRole = "fn" | "op" | "num" | "util" | "equals";

interface Key {
  label: string;
  insert?: string; // text appended to the expression (defaults to label)
  role: KeyRole;
  action?: "clear" | "back" | "equals";
  span?: number; // number of grid columns this key occupies (default 1)
  title?: string; // tooltip
}

const KEYS: Key[] = [
  { label: "sin", insert: "sin(", role: "fn" },
  { label: "cos", insert: "cos(", role: "fn" },
  { label: "tan", insert: "tan(", role: "fn" },
  { label: "√", insert: "sqrt(", role: "fn" },
  { label: "^", role: "fn" },

  { label: "ln", insert: "ln(", role: "fn" },
  { label: "log", insert: "log(", role: "fn" },
  { label: "π", insert: "pi", role: "fn" },
  { label: "e", role: "fn" },
  { label: "%", role: "fn" },

  { label: "xⁿ", insert: "pow(", role: "fn" },
  { label: "max", insert: "max(", role: "fn" },
  { label: "min", insert: "min(", role: "fn" },
  { label: "abs", insert: "abs(", role: "fn" },
  { label: ",", role: "fn" },

  { label: "C", role: "util", action: "clear" },
  { label: "(", role: "util" },
  { label: ")", role: "util" },
  { label: "Ans", insert: "ans", role: "util" },
  { label: "⌫", role: "util", action: "back" },

  { label: "7", role: "num" },
  { label: "8", role: "num" },
  { label: "9", role: "num" },
  { label: "÷", insert: "/", role: "op" },
  { label: "×", insert: "*", role: "op" },

  { label: "4", role: "num" },
  { label: "5", role: "num" },
  { label: "6", role: "num" },
  { label: "−", insert: "-", role: "op" },
  { label: "+", role: "op" },

  { label: "1", role: "num" },
  { label: "2", role: "num" },
  { label: "3", role: "num" },
  { label: ".", role: "num" },
  { label: "=", role: "equals", action: "equals" },

  // Bottom row: a wide "0" plus assignment helper.
  { label: "0", role: "num", span: 4 },
  { label: "x =", insert: "x = ", role: "util", title: "Assign to variable x" },
];

interface HistoryEntry {
  expr: string;
  result: string;
}

// ---------------------------------------------------------------------------
// Application state.
// ---------------------------------------------------------------------------
let expression = "";
let resultText = "0";
let errorText = "";
let history: HistoryEntry[] = [];
let busy = false;
let variables: Variables = {};
// True immediately after "=", so the next digit/function starts a fresh
// expression instead of appending to the previous one.
let justEvaluated = false;

const app = document.querySelector<HTMLDivElement>("#app")!;
app.innerHTML = renderShell();

const exprEl = app.querySelector<HTMLDivElement>("#expr")!;
const resultEl = app.querySelector<HTMLDivElement>("#result")!;
const errorEl = app.querySelector<HTMLDivElement>("#error")!;
const historyEl = app.querySelector<HTMLUListElement>("#history-list")!;
const varsEl = app.querySelector<HTMLUListElement>("#vars-list")!;
const keypadEl = app.querySelector<HTMLDivElement>("#keypad")!;

keypadEl.innerHTML = KEYS.map((k, i) => {
  const style = k.span ? ` style="grid-column: span ${k.span};"` : "";
  const title = k.title ? ` title="${k.title}"` : "";
  return `<button class="key key--${k.role}" data-i="${i}" type="button"${style}${title}>${k.label}</button>`;
}).join("");

// ---------------------------------------------------------------------------
// Event wiring.
// ---------------------------------------------------------------------------
keypadEl.addEventListener("click", (e) => {
  const btn = (e.target as HTMLElement).closest<HTMLButtonElement>(".key");
  if (!btn) return;
  handleKey(KEYS[Number(btn.dataset.i)]);
});

app
  .querySelector<HTMLButtonElement>("#clear-history")!
  .addEventListener("click", () => {
    history = [];
    renderHistory();
  });

historyEl.addEventListener("click", (e) => {
  const li = (e.target as HTMLElement).closest<HTMLLIElement>("li[data-expr]");
  if (!li) return;
  if (justEvaluated) {
    expression = "";
    justEvaluated = false;
  }
  expression = li.dataset.expr!;
  errorText = "";
  render();
});

app
  .querySelector<HTMLButtonElement>("#clear-vars")!
  .addEventListener("click", () => {
    variables = {};
    renderVariables();
  });

varsEl.addEventListener("click", (e) => {
  const li = (e.target as HTMLElement).closest<HTMLLIElement>("li[data-var]");
  if (!li) return;
  applyFreshStart(li.dataset.var!);
  expression += li.dataset.var!;
  errorText = "";
  render();
});

window.addEventListener("keydown", onKeyboard);

render();
renderHistory();
renderVariables();

// ---------------------------------------------------------------------------
// Behaviour.
// ---------------------------------------------------------------------------
function handleKey(key: Key) {
  if (busy) return;
  errorText = "";

  switch (key.action) {
    case "clear":
      expression = "";
      resultText = "0";
      justEvaluated = false;
      break;
    case "back":
      expression = expression.slice(0, -1);
      break;
    case "equals":
      void evaluate();
      return;
    default: {
      const text = key.insert ?? key.label;
      applyFreshStart(text);
      expression += text;
    }
  }
  render();
}

// applyFreshStart handles input right after pressing "=". An operator continues
// from the previous result (1024 then "+" -> "ans+"); anything else begins a
// brand-new expression.
function applyFreshStart(next: string) {
  if (!justEvaluated) return;
  justEvaluated = false;
  if (/^[+\-*/%^]/.test(next)) {
    expression = "ans";
  } else {
    expression = "";
  }
}

async function evaluate() {
  const expr = expression.trim();
  if (!expr) return;

  busy = true;
  render();
  try {
    const { result, variables: updated } = await calculate(expr, variables);
    variables = updated;
    resultText = formatNumber(result);
    errorText = "";
    expression = resultText;
    justEvaluated = true;
    history = [{ expr, result: resultText }, ...history].slice(0, 12);
    renderHistory();
    renderVariables();
  } catch (err) {
    errorText = err instanceof Error ? err.message : "Calculation failed";
  } finally {
    busy = false;
    render();
  }
}

function onKeyboard(e: KeyboardEvent) {
  if (busy) return;
  const { key } = e;

  if (key === "Enter" || key === "=") {
    e.preventDefault();
    void evaluate();
    return;
  }
  if (key === "Escape") {
    expression = "";
    resultText = "0";
    errorText = "";
    justEvaluated = false;
    render();
    return;
  }
  if (key === "Backspace") {
    e.preventDefault();
    expression = expression.slice(0, -1);
    errorText = "";
    render();
    return;
  }
  if (/^[0-9+\-*/%^(),.=]$/.test(key)) {
    applyFreshStart(key);
    expression += key;
    errorText = "";
    render();
  }
}

// ---------------------------------------------------------------------------
// Rendering.
// ---------------------------------------------------------------------------
function render() {
  exprEl.textContent = expression || " ";
  resultEl.textContent = busy ? "…" : resultText;
  errorEl.textContent = errorText;
  errorEl.classList.toggle("is-visible", Boolean(errorText));
  resultEl.classList.toggle("is-error", Boolean(errorText));
}

function renderHistory() {
  if (history.length === 0) {
    historyEl.innerHTML = `<li class="history__empty">No calculations yet</li>`;
    return;
  }
  historyEl.innerHTML = history
    .map(
      (h) => `
      <li data-expr="${escapeHtml(h.expr)}" title="Click to reuse">
        <span class="history__expr">${escapeHtml(h.expr)}</span>
        <span class="history__result">= ${escapeHtml(h.result)}</span>
      </li>`,
    )
    .join("");
}

function renderVariables() {
  const entries = Object.entries(variables);
  if (entries.length === 0) {
    varsEl.innerHTML = `<li class="card__empty">Assign one with <code>x = 5</code></li>`;
    return;
  }
  // "ans" first, then alphabetical.
  entries.sort(([a], [b]) =>
    a === "ans" ? -1 : b === "ans" ? 1 : a.localeCompare(b),
  );
  varsEl.innerHTML = entries
    .map(
      ([name, value]) => `
      <li data-var="${escapeHtml(name)}" title="Click to insert">
        <span class="vars__name">${escapeHtml(name)}</span>
        <span class="vars__value">${escapeHtml(formatNumber(value))}</span>
      </li>`,
    )
    .join("");
}

function renderShell(): string {
  return `
    <main class="calculator" role="application" aria-label="Scientific calculator">
      <header class="brand">
        <div class="brand__mark" aria-hidden="true">∑</div>
        <div class="brand__text">
          <h1>Calc</h1>
          <p>Scientific Calculator</p>
        </div>
        <span class="brand__badge">Powered by Go</span>
      </header>

      <section class="panel">
        <div class="screen">
          <div id="expr" class="screen__expr"></div>
          <div id="result" class="screen__result">0</div>
          <div id="error" class="screen__error" role="alert"></div>
        </div>
        <div id="keypad" class="keypad"></div>
      </section>

      <aside class="sidebar">
        <section class="card vars">
          <div class="card__head">
            <h2>Variables</h2>
            <button id="clear-vars" type="button" class="card__action">Reset</button>
          </div>
          <ul id="vars-list" class="vars__list"></ul>
        </section>

        <section class="card history">
          <div class="card__head">
            <h2>History</h2>
            <button id="clear-history" type="button" class="card__action">Clear</button>
          </div>
          <ul id="history-list" class="history__list"></ul>
        </section>
      </aside>
    </main>
  `;
}

// ---------------------------------------------------------------------------
// Helpers.
// ---------------------------------------------------------------------------
function formatNumber(n: number): string {
  if (!Number.isFinite(n)) return String(n);
  // Trim floating-point noise while keeping reasonable precision.
  return String(Number(n.toPrecision(12)));
}

function escapeHtml(s: string): string {
  return s
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}
