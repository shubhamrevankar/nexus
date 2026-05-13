# GitHub Integration Slice

This slice introduces the first external-tool integration for Nexus.

## What It Adds

- A PostgreSQL table for GitHub repository metadata.
- A Go GitHub API client for fetching repository details.
- Authenticated API endpoints for connecting and listing repositories.
- Workspace-level authorization before any repository data is saved or read.
- A workspace UI panel for connecting a GitHub repository.

## API Endpoints

### Connect Repository

```http
POST /v1/integrations/github/repositories
Authorization: Bearer <nexus-session-token>
Content-Type: application/json
```

```json
{
  "workspaceId": "workspace-uuid",
  "owner": "openai",
  "repository": "codex",
  "token": "github-token"
}
```

The GitHub token is used for the request only. It is not stored by Nexus.

### List Connected Repositories

```http
GET /v1/integrations/github/repositories?workspaceId=<workspace-uuid>
Authorization: Bearer <nexus-session-token>
```

## Data Stored

Nexus stores repository metadata only:

- GitHub repository ID.
- Owner, name, and full name.
- Description and HTML URL.
- Default branch and language.
- Visibility, stars, forks, open issues.
- Last pushed time and last synced time.

## Current Scope

This is an ingestion foundation, not a full GitHub app integration.

Implemented now:

- Manual repository sync by owner/repository name.
- Temporary token use from the browser request.
- Workspace authorization checks.

Still needed later:

- GitHub App installation flow.
- Secure token storage or installation token exchange.
- Background sync jobs.
- Commit, pull request, issue, and workflow-run ingestion.
- Webhooks.
- Rate-limit handling and retry policies.
- Audit logs for integration changes.

