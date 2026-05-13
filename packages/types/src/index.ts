export type ServiceStatus = "ok" | "degraded" | "unavailable";

export interface HealthStatus {
  service: string;
  status: ServiceStatus;
  version?: string;
}

export interface User {
  id: string;
  email: string;
  name: string;
  createdAt: string;
}

export interface Organization {
  id: string;
  name: string;
  role: "owner" | "admin" | "member";
  createdAt: string;
}

export interface Workspace {
  id: string;
  organizationId: string;
  name: string;
  slug: string;
  createdAt: string;
}

export interface WorkspaceSummary {
  organization: Organization;
  workspaces: Workspace[];
}

export interface AuthSession {
  token: string;
  expiresAt: string;
  user: User;
}

export interface AuthResponse {
  session: AuthSession;
  workspaceSet: WorkspaceSummary;
}

export interface GitHubRepository {
  id: string;
  workspaceId: string;
  githubId: number;
  owner: string;
  name: string;
  fullName: string;
  description: string;
  htmlUrl: string;
  defaultBranch: string;
  language: string;
  private: boolean;
  stars: number;
  forks: number;
  openIssues: number;
  pushedAt?: string | null;
  syncedAt: string;
}
