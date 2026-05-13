package github

import (
	"context"
	"database/sql"
)

type RepositoryStore struct {
	db *sql.DB
}

func NewRepositoryStore(db *sql.DB) *RepositoryStore {
	return &RepositoryStore{db: db}
}

func (store *RepositoryStore) UpsertRepository(ctx context.Context, workspaceID string, repository Repository) (Repository, error) {
	var saved Repository
	err := store.db.QueryRowContext(ctx, `
INSERT INTO github_repositories (
  workspace_id,
  github_id,
  owner,
  name,
  full_name,
  description,
  html_url,
  default_branch,
  language,
  private,
  stars,
  forks,
  open_issues,
  pushed_at,
  synced_at,
  updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, now(), now())
ON CONFLICT (workspace_id, full_name)
DO UPDATE SET
  github_id = EXCLUDED.github_id,
  owner = EXCLUDED.owner,
  name = EXCLUDED.name,
  description = EXCLUDED.description,
  html_url = EXCLUDED.html_url,
  default_branch = EXCLUDED.default_branch,
  language = EXCLUDED.language,
  private = EXCLUDED.private,
  stars = EXCLUDED.stars,
  forks = EXCLUDED.forks,
  open_issues = EXCLUDED.open_issues,
  pushed_at = EXCLUDED.pushed_at,
  synced_at = now(),
  updated_at = now()
RETURNING id::text, workspace_id::text, github_id, owner, name, full_name, description, html_url, default_branch, language, private, stars, forks, open_issues, pushed_at, synced_at`,
		workspaceID,
		repository.GitHubID,
		repository.Owner,
		repository.Name,
		repository.FullName,
		repository.Description,
		repository.HTMLURL,
		repository.DefaultBranch,
		repository.Language,
		repository.Private,
		repository.Stars,
		repository.Forks,
		repository.OpenIssues,
		repository.PushedAt,
	).Scan(
		&saved.ID,
		&saved.WorkspaceID,
		&saved.GitHubID,
		&saved.Owner,
		&saved.Name,
		&saved.FullName,
		&saved.Description,
		&saved.HTMLURL,
		&saved.DefaultBranch,
		&saved.Language,
		&saved.Private,
		&saved.Stars,
		&saved.Forks,
		&saved.OpenIssues,
		&saved.PushedAt,
		&saved.SyncedAt,
	)
	return saved, err
}

func (store *RepositoryStore) ListRepositories(ctx context.Context, workspaceID string) ([]Repository, error) {
	rows, err := store.db.QueryContext(ctx, `
SELECT id::text, workspace_id::text, github_id, owner, name, full_name, description, html_url, default_branch, language, private, stars, forks, open_issues, pushed_at, synced_at
FROM github_repositories
WHERE workspace_id = $1
ORDER BY synced_at DESC`, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	repositories := make([]Repository, 0)
	for rows.Next() {
		var repository Repository
		if err := rows.Scan(
			&repository.ID,
			&repository.WorkspaceID,
			&repository.GitHubID,
			&repository.Owner,
			&repository.Name,
			&repository.FullName,
			&repository.Description,
			&repository.HTMLURL,
			&repository.DefaultBranch,
			&repository.Language,
			&repository.Private,
			&repository.Stars,
			&repository.Forks,
			&repository.OpenIssues,
			&repository.PushedAt,
			&repository.SyncedAt,
		); err != nil {
			return nil, err
		}
		repositories = append(repositories, repository)
	}

	return repositories, rows.Err()
}

func (store *RepositoryStore) GetRepository(ctx context.Context, workspaceID string, repositoryID string) (Repository, error) {
	var repository Repository
	err := store.db.QueryRowContext(ctx, `
SELECT id::text, workspace_id::text, github_id, owner, name, full_name, description, html_url, default_branch, language, private, stars, forks, open_issues, pushed_at, synced_at
FROM github_repositories
WHERE workspace_id = $1 AND id = $2`, workspaceID, repositoryID).Scan(
		&repository.ID,
		&repository.WorkspaceID,
		&repository.GitHubID,
		&repository.Owner,
		&repository.Name,
		&repository.FullName,
		&repository.Description,
		&repository.HTMLURL,
		&repository.DefaultBranch,
		&repository.Language,
		&repository.Private,
		&repository.Stars,
		&repository.Forks,
		&repository.OpenIssues,
		&repository.PushedAt,
		&repository.SyncedAt,
	)
	return repository, err
}

func (store *RepositoryStore) UpsertFile(ctx context.Context, repositoryID string, file RepositoryFile) (RepositoryFile, error) {
	var saved RepositoryFile
	err := store.db.QueryRowContext(ctx, `
INSERT INTO github_repository_files (repository_id, path, sha, size, content_text, indexed, skipped_reason, synced_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
ON CONFLICT (repository_id, path)
DO UPDATE SET
  sha = EXCLUDED.sha,
  size = EXCLUDED.size,
  content_text = EXCLUDED.content_text,
  indexed = EXCLUDED.indexed,
  skipped_reason = EXCLUDED.skipped_reason,
  synced_at = now(),
  updated_at = now()
RETURNING id::text, repository_id::text, path, sha, size, COALESCE(content_text, ''), indexed, skipped_reason, synced_at`,
		repositoryID,
		file.Path,
		file.SHA,
		file.Size,
		nullableContent(file.ContentText, file.Indexed),
		file.Indexed,
		file.SkippedReason,
	).Scan(
		&saved.ID,
		&saved.RepositoryID,
		&saved.Path,
		&saved.SHA,
		&saved.Size,
		&saved.ContentText,
		&saved.Indexed,
		&saved.SkippedReason,
		&saved.SyncedAt,
	)
	return saved, err
}

func (store *RepositoryStore) ListFiles(ctx context.Context, repositoryID string, indexedOnly bool, limit int) ([]RepositoryFile, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}

	query := `
SELECT id::text, repository_id::text, path, sha, size, COALESCE(content_text, ''), indexed, skipped_reason, synced_at
FROM github_repository_files
WHERE repository_id = $1`
	args := []any{repositoryID, limit}
	if indexedOnly {
		query += " AND indexed = true"
	}
	query += " ORDER BY path ASC LIMIT $2"

	rows, err := store.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]RepositoryFile, 0)
	for rows.Next() {
		var file RepositoryFile
		if err := rows.Scan(
			&file.ID,
			&file.RepositoryID,
			&file.Path,
			&file.SHA,
			&file.Size,
			&file.ContentText,
			&file.Indexed,
			&file.SkippedReason,
			&file.SyncedAt,
		); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, rows.Err()
}

func nullableContent(content string, indexed bool) sql.NullString {
	if !indexed {
		return sql.NullString{}
	}

	return sql.NullString{String: content, Valid: true}
}
