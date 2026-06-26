import "./style.css";
import { calculate } from "./api";

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
  { label: "⌫", role: "util", action: "back" },
  { label: "÷", insert: "/", role: "op" },

  { label: "7", role: "num" },
  { label: "8", role: "num" },
  { label: "9", role: "num" },
  { label: "×", insert: "*", role: "op" },
  { label: "−", insert: "-", role: "op" },

  { label: "4", role: "num" },
  { label: "5", role: "num" },
  { label: "6", role: "num" },
  { label: "+", role: "op" },
  { label: "=", role: "equals", action: "equals" },

  { label: "1", role: "num" },
  { label: "2", role: "num" },
  { label: "3", role: "num" },
  { label: "0", role: "num" },
  { label: ".", role: "num" },
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

const app = document.querySelector<HTMLDivElement>("#app")!;
app.innerHTML = renderShell();

const exprEl = app.querySelector<HTMLDivElement>("#expr")!;
const resultEl = app.querySelector<HTMLDivElement>("#result")!;
const errorEl = app.querySelector<HTMLDivElement>("#error")!;
const historyEl = app.querySelector<HTMLUListElement>("#history-list")!;
const keypadEl = app.querySelector<HTMLDivElement>("#keypad")!;

keypadEl.innerHTML = KEYS.map(
  (k, i) =>
    `<button class="key key--${k.role}" data-i="${i}" type="button">${k.label}</button>`,
).join("");

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
  expression = li.dataset.expr!;
  errorText = "";
  render();
});

window.addEventListener("keydown", onKeyboard);

render();
renderHistory();

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
      break;
    case "back":
      expression = expression.slice(0, -1);
      break;
    case "equals":
      void evaluate();
      return;
    default:
      expression += key.insert ?? key.label;
  }
  render();
}

async function evaluate() {
  const expr = expression.trim();
  if (!expr) return;

  busy = true;
  render();
  try {
    const value = await calculate(expr);
    resultText = formatNumber(value);
    errorText = "";
    history = [{ expr, result: resultText }, ...history].slice(0, 12);
    renderHistory();
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
  if (/^[0-9+\-*/%^(),.]$/.test(key)) {
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

      <aside class="history">
        <div class="history__head">
          <h2>History</h2>
          <button id="clear-history" type="button" class="history__clear">Clear</button>
        </div>
        <ul id="history-list" class="history__list"></ul>
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
