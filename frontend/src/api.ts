// API client for the Go calculator backend.
//
// In development, Vite proxies "/api/*" to the Go server (see vite.config.ts),
// so we use a relative base URL. Override with VITE_API_BASE if needed.
const API_BASE = import.meta.env.VITE_API_BASE ?? "/api";

export interface CalcSuccess {
  result: number;
}

export interface CalcError {
  error: string;
}

/**
 * Evaluate an arithmetic expression on the backend.
 * Throws an Error (with the backend's message) on failure.
 */
export async function calculate(expression: string): Promise<number> {
  let res: Response;
  try {
    res = await fetch(`${API_BASE}/calculate`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ expression }),
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

  return data.result;
}
