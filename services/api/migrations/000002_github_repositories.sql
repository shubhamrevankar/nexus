CREATE TABLE IF NOT EXISTS github_repositories (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  workspace_id uuid NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
  github_id bigint NOT NULL,
  owner text NOT NULL,
  name text NOT NULL,
  full_name text NOT NULL,
  description text NOT NULL DEFAULT '',
  html_url text NOT NULL,
  default_branch text NOT NULL,
  language text NOT NULL DEFAULT '',
  private boolean NOT NULL,
  stars integer NOT NULL DEFAULT 0,
  forks integer NOT NULL DEFAULT 0,
  open_issues integer NOT NULL DEFAULT 0,
  pushed_at timestamptz,
  synced_at timestamptz NOT NULL DEFAULT now(),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (workspace_id, full_name)
);

CREATE INDEX IF NOT EXISTS github_repositories_workspace_id_idx ON github_repositories (workspace_id);

