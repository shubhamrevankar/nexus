"use client";

import type { AuthResponse } from "@nexus/types";
import { useRouter } from "next/navigation";
import type { FormEvent } from "react";
import { useState } from "react";
import { login, register } from "@/lib/api";
import { saveSession } from "@/lib/session";

type Mode = "register" | "login";

export function AuthPanel() {
  const router = useRouter();
  const [mode, setMode] = useState<Mode>("register");
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError(null);
    setIsSubmitting(true);

    const form = new FormData(event.currentTarget);

    try {
      let response: AuthResponse;
      if (mode === "register") {
        response = await register({
          email: String(form.get("email") ?? ""),
          name: String(form.get("name") ?? ""),
          password: String(form.get("password") ?? ""),
          organizationName: String(form.get("organizationName") ?? ""),
          workspaceName: String(form.get("workspaceName") ?? "")
        });
      } else {
        response = await login({
          email: String(form.get("email") ?? ""),
          password: String(form.get("password") ?? "")
        });
      }

      saveSession(response);
      router.push("/workspace");
    } catch (caughtError) {
      setError(caughtError instanceof Error ? caughtError.message : "Something went wrong");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <section className="rounded-3xl border border-white/10 bg-white/[0.04] p-6 shadow-2xl shadow-violet-950/30 backdrop-blur">
      <div className="mb-6 flex rounded-2xl bg-black/30 p-1 text-sm">
        <button
          className={tabClass(mode === "register")}
          type="button"
          onClick={() => setMode("register")}
        >
          Create workspace
        </button>
        <button className={tabClass(mode === "login")} type="button" onClick={() => setMode("login")}>
          Sign in
        </button>
      </div>

      <form className="space-y-4" onSubmit={handleSubmit}>
        {mode === "register" ? (
          <label className="block text-sm text-zinc-300">
            Name
            <input className={inputClass} name="name" placeholder="Ada Lovelace" required />
          </label>
        ) : null}

        <label className="block text-sm text-zinc-300">
          Email
          <input className={inputClass} name="email" placeholder="you@company.com" required type="email" />
        </label>

        <label className="block text-sm text-zinc-300">
          Password
          <input className={inputClass} minLength={8} name="password" required type="password" />
        </label>

        {mode === "register" ? (
          <div className="grid gap-4 sm:grid-cols-2">
            <label className="block text-sm text-zinc-300">
              Organization
              <input className={inputClass} name="organizationName" placeholder="Acme Labs" required />
            </label>
            <label className="block text-sm text-zinc-300">
              Workspace
              <input className={inputClass} name="workspaceName" placeholder="Engineering" required />
            </label>
          </div>
        ) : null}

        {error ? <p className="rounded-2xl border border-red-400/30 bg-red-950/30 p-3 text-sm text-red-200">{error}</p> : null}

        <button
          className="w-full rounded-2xl bg-violet-500 px-4 py-3 text-sm font-semibold text-white transition hover:bg-violet-400 disabled:cursor-not-allowed disabled:opacity-60"
          disabled={isSubmitting}
          type="submit"
        >
          {isSubmitting ? "Working..." : mode === "register" ? "Create Nexus workspace" : "Sign in"}
        </button>
      </form>
    </section>
  );
}

const inputClass =
  "mt-2 w-full rounded-2xl border border-white/10 bg-black/30 px-4 py-3 text-white outline-none ring-violet-400 transition placeholder:text-zinc-600 focus:ring-2";

function tabClass(active: boolean) {
  return `flex-1 rounded-xl px-3 py-2 font-medium transition ${
    active ? "bg-violet-500 text-white" : "text-zinc-400 hover:text-white"
  }`;
}
