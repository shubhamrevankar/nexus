"use client";

import type { AuthResponse, WorkspaceSummary } from "@nexus/types";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { GitHubIntegrationPanel } from "@/components/github-integration-panel";
import { listWorkspaces } from "@/lib/api";
import { clearSession, readSession } from "@/lib/session";

export function WorkspaceDashboard() {
  const router = useRouter();
  const [session, setSession] = useState<AuthResponse | null>(null);
  const [workspaceSets, setWorkspaceSets] = useState<WorkspaceSummary[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const storedSession = readSession();
    if (!storedSession) {
      router.push("/auth");
      return;
    }

    setSession(storedSession);
    listWorkspaces(storedSession.session.token)
      .then(setWorkspaceSets)
      .catch((caughtError) => setError(caughtError instanceof Error ? caughtError.message : "Failed to load workspaces"));
  }, [router]);

  function signOut() {
    clearSession();
    router.push("/auth");
  }

  if (!session) {
    return <main className="min-h-screen px-6 py-12 text-zinc-300">Loading workspace...</main>;
  }

  const activeWorkspaceId = workspaceSets[0]?.workspaces[0]?.id ?? session.workspaceSet.workspaces[0]?.id ?? null;

  return (
    <main className="min-h-screen px-6 py-10 text-white">
      <section className="mx-auto max-w-6xl">
        <header className="flex flex-col gap-4 border-b border-white/10 pb-8 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <p className="text-sm uppercase tracking-[0.3em] text-violet-300">Nexus workspace</p>
            <h1 className="mt-3 text-4xl font-semibold tracking-tight">Welcome, {session.session.user.name}</h1>
            <p className="mt-2 text-zinc-400">Your organization and workspace foundation is connected to the Go API.</p>
          </div>
          <button className="rounded-2xl border border-white/10 px-4 py-2 text-sm text-zinc-300 hover:text-white" onClick={signOut}>
            Sign out
          </button>
        </header>

        {error ? <p className="mt-6 rounded-2xl border border-red-400/30 bg-red-950/30 p-4 text-red-200">{error}</p> : null}

        <div className="mt-8 grid gap-6 lg:grid-cols-[1fr_360px]">
          <section className="rounded-3xl border border-white/10 bg-white/[0.04] p-6">
            <h2 className="text-xl font-semibold">Organizations</h2>
            <div className="mt-6 space-y-4">
              {workspaceSets.map((item) => (
                <article className="rounded-2xl border border-white/10 bg-black/20 p-5" key={item.organization.id}>
                  <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                    <div>
                      <h3 className="text-lg font-semibold">{item.organization.name}</h3>
                      <p className="text-sm text-zinc-500">Role: {item.organization.role}</p>
                    </div>
                    <span className="rounded-full bg-emerald-400/10 px-3 py-1 text-xs font-medium text-emerald-300">Active</span>
                  </div>
                  <div className="mt-4 grid gap-3 sm:grid-cols-2">
                    {item.workspaces.map((workspace) => (
                      <div className="rounded-2xl bg-white/[0.04] p-4" key={workspace.id}>
                        <p className="font-medium">{workspace.name}</p>
                        <p className="mt-1 text-sm text-zinc-500">/{workspace.slug}</p>
                      </div>
                    ))}
                  </div>
                </article>
              ))}
            </div>
          </section>

          <GitHubIntegrationPanel token={session.session.token} workspaceId={activeWorkspaceId} />
        </div>
      </section>
    </main>
  );
}
