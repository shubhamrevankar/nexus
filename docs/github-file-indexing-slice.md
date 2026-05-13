# GitHub File Indexing Slice

This slice lets Nexus ingest selected text files from a connected GitHub repository.

## What It Adds

- A PostgreSQL table for repository files.
- GitHub tree fetching from the repository default branch.
- Raw file content fetching for safe text-sized files.
- Skip rules for large files, lockfiles, generated paths, media/binary files, and sensitive-looking paths.
- API endpoints for syncing files and listing indexed files.
- Workspace UI controls to index and inspect files.

## API Endpoints

### Sync Repository Files

```http
POST /v1/integrations/github/files/sync
Authorization: Bearer <nexus-session-token>
Content-Type: application/json
```

```json
{
  "workspaceId": "workspace-uuid",
  "repositoryId": "repository-uuid",
  "token": "github-token",
  "maxFiles": 50
}
```

The GitHub token is used for this request only and is not stored.

### List Indexed Files

```http
GET /v1/integrations/github/files?workspaceId=<workspace-uuid>&repositoryId=<repository-uuid>
Authorization: Bearer <nexus-session-token>
```

## Data Stored

Nexus stores selected file metadata and text content:

- Repository ID.
- File path.
- Git blob SHA.
- File size.
- Text content for indexed files.
- Indexed/skipped status.
- Skip reason.
- Synced timestamp.

## Safety Limits

Current indexing is intentionally conservative:

- Maximum 50 indexed files per sync request.
- Maximum 100 KB per indexed file.
- Skips dependency/generated paths such as `node_modules`, `vendor`, `dist`, and `build`.
- Skips media/binary extensions.
- Skips lockfiles.
- Skips sensitive-looking paths containing `.env`, `secret`, `credential`, or `private_key`.

## Current Scope

This is source ingestion, not semantic search yet.

Still needed later:

- Background ingestion jobs.
- Better language detection.
- Secret scanning before persistence.
- Chunking files for retrieval.
- Embeddings and vector indexing.
- Code-aware search and citations.

