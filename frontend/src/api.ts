// API client for the Go calculator backend.
//
// In development, Vite proxies "/api/*" to the Go server (see vite.config.ts),
// so we use a relative base URL. Override with VITE_API_BASE if needed.
const API_BASE = import.meta.env.VITE_API_BASE ?? "/api";

export type Variables = Record<string, number>;

export interface CalcSuccess {
  result: number;
  variables?: Variables;
}

export interface CalcError {
  error: string;
}

export interface CalcResult {
  result: number;
  variables: Variables;
}

/**
 * Evaluate an arithmetic expression on the backend, carrying a variable
 * environment so assignments (x = 5) and references (x * 2, ans) work across
 * calls. Returns the result plus the updated variable map.
 * Throws an Error (with the backend's message) on failure.
 */
export async function calculate(
  expression: string,
  variables: Variables = {},
): Promise<CalcResult> {
  let res: Response;
  try {
    res = await fetch(`${API_BASE}/calculate`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ expression, variables }),
    });
  } catch {
    throw new Error("Cannot reach the calculator service");
  }

  const data = (await res.json().catch(() => null)) as
    | CalcSuccess
    | CalcError
    | null;

  if (!res.ok || data === null || "error" in data) {
    const message =
      data && "error" in data ? data.error : `Request failed (${res.status})`;
    throw new Error(message);
  }

  return { result: data.result, variables: data.variables ?? {} };
}
