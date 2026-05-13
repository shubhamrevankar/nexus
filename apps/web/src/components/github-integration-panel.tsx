"use client";

import type { GitHubRepository } from "@nexus/types";
import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import { ingestGitHubRepository, listGitHubRepositories } from "@/lib/api";

interface GitHubIntegrationPanelProps {
  token: string;
  workspaceId: string | null;
}

export function GitHubIntegrationPanel({ token, workspaceId }: GitHubIntegrationPanelProps) {
  const [repositories, setRepositories] = useState<GitHubRepository[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (!workspaceId) {
      return;
    }

    listGitHubRepositories(token, workspaceId)
      .then(setRepositories)
      .catch((caughtError) => setError(caughtError instanceof Error ? caughtError.message : "Failed to load GitHub repositories"));
  }, [token, workspaceId]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const formElement = event.currentTarget;

    if (!workspaceId) {
      setError("Create or select a workspace before connecting GitHub.");
      return;
    }

    setError(null);
    setIsSubmitting(true);

    const form = new FormData(event.currentTarget);

    try {
      const repository = await ingestGitHubRepository(token, {
        workspaceId,
        owner: String(form.get("owner") ?? ""),
        repository: String(form.get("repository") ?? ""),
        token: String(form.get("token") ?? "")
      });

      setRepositories((current) => [repository, ...current.filter((item) => item.id !== repository.id)]);
      formElement.reset();
    } catch (caughtError) {
      setError(caughtError instanceof Error ? caughtError.message : "GitHub repository could not be connected");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <aside className="rounded-3xl border border-white/10 bg-white/[0.04] p-6">
      <h2 className="text-xl font-semibold">GitHub integration</h2>
      <p className="mt-2 text-sm leading-6 text-zinc-400">
        Connect one repository with a temporary GitHub token. Nexus fetches repository metadata and does not store the token.
      </p>

      <form className="mt-5 space-y-3" onSubmit={handleSubmit}>
        <input className={inputClass} name="owner" placeholder="owner, e.g. openai" required />
        <input className={inputClass} name="repository" placeholder="repository, e.g. codex" required />
        <input className={inputClass} name="token" placeholder="GitHub fine-grained token" required type="password" />

        {error ? <p className="rounded-2xl border border-red-400/30 bg-red-950/30 p-3 text-sm text-red-200">{error}</p> : null}

        <button
          className="w-full rounded-2xl bg-violet-500 px-4 py-3 text-sm font-semibold text-white transition hover:bg-violet-400 disabled:cursor-not-allowed disabled:opacity-60"
          disabled={isSubmitting || !workspaceId}
          type="submit"
        >
          {isSubmitting ? "Syncing..." : "Connect repository"}
        </button>
      </form>

      <div className="mt-6 space-y-3">
        {repositories.length === 0 ? (
          <p className="rounded-2xl bg-black/20 p-4 text-sm text-zinc-500">No GitHub repositories connected yet.</p>
        ) : (
          repositories.map((repository) => (
            <article className="rounded-2xl border border-white/10 bg-black/20 p-4" key={repository.id}>
              <div className="flex items-start justify-between gap-3">
                <div>
                  <h3 className="font-semibold">{repository.fullName}</h3>
                  <p className="mt-1 line-clamp-2 text-sm text-zinc-500">{repository.description || "No description"}</p>
                </div>
                <span className="rounded-full bg-emerald-400/10 px-2 py-1 text-xs text-emerald-300">synced</span>
              </div>
              <dl className="mt-4 grid grid-cols-3 gap-2 text-xs text-zinc-400">
                <div>
                  <dt>Stars</dt>
                  <dd className="mt-1 text-white">{repository.stars}</dd>
                </div>
                <div>
                  <dt>Forks</dt>
                  <dd className="mt-1 text-white">{repository.forks}</dd>
                </div>
                <div>
                  <dt>Branch</dt>
                  <dd className="mt-1 text-white">{repository.defaultBranch}</dd>
                </div>
              </dl>
            </article>
          ))
        )}
      </div>
    </aside>
  );
}

const inputClass =
  "w-full rounded-2xl border border-white/10 bg-black/30 px-4 py-3 text-sm text-white outline-none ring-violet-400 transition placeholder:text-zinc-600 focus:ring-2";
