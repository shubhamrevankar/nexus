import type { AuthResponse } from "@nexus/types";

const sessionKey = "nexus.session";

export function readSession(): AuthResponse | null {
  if (typeof window === "undefined") {
    return null;
  }

  const value = window.localStorage.getItem(sessionKey);
  if (!value) {
    return null;
  }

  try {
    return JSON.parse(value) as AuthResponse;
  } catch {
    window.localStorage.removeItem(sessionKey);
    return null;
  }
}

export function saveSession(session: AuthResponse): void {
  window.localStorage.setItem(sessionKey, JSON.stringify(session));
}

export function clearSession(): void {
  window.localStorage.removeItem(sessionKey);
}

