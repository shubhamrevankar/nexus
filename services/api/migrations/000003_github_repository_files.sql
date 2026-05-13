CREATE TABLE IF NOT EXISTS github_repository_files (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  repository_id uuid NOT NULL REFERENCES github_repositories(id) ON DELETE CASCADE,
  path text NOT NULL,
  sha text NOT NULL,
  size integer NOT NULL DEFAULT 0,
  content_text text,
  indexed boolean NOT NULL DEFAULT false,
  skipped_reason text NOT NULL DEFAULT '',
  synced_at timestamptz NOT NULL DEFAULT now(),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (repository_id, path)
);

CREATE INDEX IF NOT EXISTS github_repository_files_repository_id_idx ON github_repository_files (repository_id);
CREATE INDEX IF NOT EXISTS github_repository_files_indexed_idx ON github_repository_files (repository_id, indexed);

