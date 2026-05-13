import type { AuthResponse, GitHubFileSyncResult, GitHubRepository, GitHubRepositoryFile, WorkspaceSummary } from "@nexus/types";

const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export interface RegisterPayload {
  email: string;
  name: string;
  password: string;
  organizationName: string;
  workspaceName: string;
}

export interface LoginPayload {
  email: string;
  password: string;
}

export interface IngestGitHubRepositoryPayload {
  workspaceId: string;
  owner: string;
  repository: string;
  token: string;
}

export interface SyncGitHubFilesPayload {
  workspaceId: string;
  repositoryId: string;
  token: string;
  maxFiles: number;
}

export async function register(payload: RegisterPayload): Promise<AuthResponse> {
  return request<AuthResponse>("/v1/auth/register", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function login(payload: LoginPayload): Promise<AuthResponse> {
  return request<AuthResponse>("/v1/auth/login", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function listWorkspaces(token: string): Promise<WorkspaceSummary[]> {
  const response = await request<{ items: WorkspaceSummary[] }>("/v1/workspaces", {
    headers: {
      Authorization: `Bearer ${token}`
    }
  });

  return response.items;
}

export async function listGitHubRepositories(token: string, workspaceId: string): Promise<GitHubRepository[]> {
  const response = await request<{ items: GitHubRepository[] }>(
    `/v1/integrations/github/repositories?workspaceId=${encodeURIComponent(workspaceId)}`,
    {
      headers: {
        Authorization: `Bearer ${token}`
      }
    }
  );

  return response.items;
}

export async function ingestGitHubRepository(token: string, payload: IngestGitHubRepositoryPayload): Promise<GitHubRepository> {
  const response = await request<{ repository: GitHubRepository }>("/v1/integrations/github/repositories", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`
    },
    body: JSON.stringify(payload)
  });

  return response.repository;
}

export async function syncGitHubFiles(token: string, payload: SyncGitHubFilesPayload): Promise<GitHubFileSyncResult> {
  const response = await request<{ result: GitHubFileSyncResult }>("/v1/integrations/github/files/sync", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`
    },
    body: JSON.stringify(payload)
  });

  return response.result;
}

export async function listGitHubFiles(token: string, workspaceId: string, repositoryId: string): Promise<GitHubRepositoryFile[]> {
  const response = await request<{ items: GitHubRepositoryFile[] }>(
    `/v1/integrations/github/files?workspaceId=${encodeURIComponent(workspaceId)}&repositoryId=${encodeURIComponent(repositoryId)}`,
    {
      headers: {
        Authorization: `Bearer ${token}`
      }
    }
  );

  return response.items;
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(`${apiBaseUrl}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init.headers
    }
  });

  if (!response.ok) {
    const payload = (await response.json().catch(() => null)) as { error?: string } | null;
    throw new Error(payload?.error ?? `Request failed with status ${response.status}`);
  }

  return (await response.json()) as T;
}
