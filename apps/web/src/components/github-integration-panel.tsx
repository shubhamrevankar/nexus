"use client";

import type { GitHubFileSyncResult, GitHubRepository, GitHubRepositoryFile } from "@nexus/types";
import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import { ingestGitHubRepository, listGitHubFiles, listGitHubRepositories, syncGitHubFiles } from "@/lib/api";

interface GitHubIntegrationPanelProps {
  token: string;
  workspaceId: string | null;
}

export function GitHubIntegrationPanel({ token, workspaceId }: GitHubIntegrationPanelProps) {
  const [repositories, setRepositories] = useState<GitHubRepository[]>([]);
  const [filesByRepository, setFilesByRepository] = useState<Record<string, GitHubRepositoryFile[]>>({});
  const [syncResultByRepository, setSyncResultByRepository] = useState<Record<string, GitHubFileSyncResult>>({});
  const [syncingRepositoryId, setSyncingRepositoryId] = useState<string | null>(null);
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

  async function handleFileSync(repository: GitHubRepository) {
    const tokenValue = window.prompt("Paste a GitHub token to index source files. Nexus will not store it.");
    if (!tokenValue || !workspaceId) {
      return;
    }

    setError(null);
    setSyncingRepositoryId(repository.id);

    try {
      const result = await syncGitHubFiles(token, {
        workspaceId,
        repositoryId: repository.id,
        token: tokenValue,
        maxFiles: 50
      });

      setSyncResultByRepository((current) => ({ ...current, [repository.id]: result }));
      setFilesByRepository((current) => ({ ...current, [repository.id]: result.files.filter((file) => file.indexed) }));
    } catch (caughtError) {
      setError(caughtError instanceof Error ? caughtError.message : "GitHub files could not be indexed");
    } finally {
      setSyncingRepositoryId(null);
    }
  }

  async function handleLoadFiles(repository: GitHubRepository) {
    if (!workspaceId) {
      return;
    }

    setError(null);
    try {
      const files = await listGitHubFiles(token, workspaceId, repository.id);
      setFilesByRepository((current) => ({ ...current, [repository.id]: files }));
    } catch (caughtError) {
      setError(caughtError instanceof Error ? caughtError.message : "GitHub files could not be loaded");
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
          repositories.map((repository) => {
            const syncResult = syncResultByRepository[repository.id];
            const indexedFiles = filesByRepository[repository.id] ?? [];

            return (
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
                <div className="mt-4 grid gap-2 sm:grid-cols-2">
                  <button
                    className="rounded-xl border border-white/10 px-3 py-2 text-xs font-medium text-zinc-300 transition hover:text-white disabled:opacity-60"
                    disabled={syncingRepositoryId === repository.id}
                    onClick={() => handleFileSync(repository)}
                    type="button"
                  >
                    {syncingRepositoryId === repository.id ? "Indexing..." : "Index source files"}
                  </button>
                  <button
                    className="rounded-xl border border-white/10 px-3 py-2 text-xs font-medium text-zinc-300 transition hover:text-white"
                    onClick={() => handleLoadFiles(repository)}
                    type="button"
                  >
                    Show indexed files
                  </button>
                </div>
                {syncResult ? (
                  <p className="mt-3 rounded-xl bg-emerald-400/10 p-3 text-xs text-emerald-200">
                    Indexed {syncResult.indexedCount} files, skipped {syncResult.skippedCount}.
                  </p>
                ) : null}
                {indexedFiles.length ? (
                  <ul className="mt-3 max-h-48 space-y-2 overflow-auto rounded-xl bg-black/30 p-3 text-xs text-zinc-300">
                    {indexedFiles.map((file) => (
                      <li className="flex justify-between gap-3" key={file.id}>
                        <span className="truncate">{file.path}</span>
                        <span className="shrink-0 text-zinc-500">{Math.ceil(file.size / 1024)} KB</span>
                      </li>
                    ))}
                  </ul>
                ) : null}
              </article>
            );
          })
        )}
      </div>
    </aside>
  );
}

const inputClass =
  "w-full rounded-2xl border border-white/10 bg-black/30 px-4 py-3 text-sm text-white outline-none ring-violet-400 transition placeholder:text-zinc-600 focus:ring-2";
