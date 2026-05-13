package tenancy

import (
	"context"
	"database/sql"
	"regexp"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (repository *Repository) CreateOrganizationWithWorkspace(ctx context.Context, userID string, organizationName string, workspaceName string) (WorkspaceSummary, error) {
	transaction, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return WorkspaceSummary{}, err
	}
	defer transaction.Rollback()

	var organization Organization
	err = transaction.QueryRowContext(ctx, `
INSERT INTO organizations (name)
VALUES ($1)
RETURNING id::text, name, created_at`, strings.TrimSpace(organizationName)).Scan(
		&organization.ID,
		&organization.Name,
		&organization.CreatedAt,
	)
	if err != nil {
		return WorkspaceSummary{}, err
	}
	organization.Role = "owner"

	if _, err := transaction.ExecContext(ctx, `
INSERT INTO organization_memberships (organization_id, user_id, role)
VALUES ($1, $2, 'owner')`, organization.ID, userID); err != nil {
		return WorkspaceSummary{}, err
	}

	workspace, err := createWorkspace(ctx, transaction, organization.ID, workspaceName)
	if err != nil {
		return WorkspaceSummary{}, err
	}

	if err := transaction.Commit(); err != nil {
		return WorkspaceSummary{}, err
	}

	return WorkspaceSummary{Organization: organization, Workspaces: []Workspace{workspace}}, nil
}

func (repository *Repository) ListForUser(ctx context.Context, userID string) ([]WorkspaceSummary, error) {
	rows, err := repository.db.QueryContext(ctx, `
SELECT organizations.id::text, organizations.name, organization_memberships.role, organizations.created_at
FROM organizations
JOIN organization_memberships ON organization_memberships.organization_id = organizations.id
WHERE organization_memberships.user_id = $1
ORDER BY organizations.created_at ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summaries := make([]WorkspaceSummary, 0)
	for rows.Next() {
		var organization Organization
		if err := rows.Scan(&organization.ID, &organization.Name, &organization.Role, &organization.CreatedAt); err != nil {
			return nil, err
		}

		workspaces, err := repository.listWorkspaces(ctx, organization.ID)
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, WorkspaceSummary{Organization: organization, Workspaces: workspaces})
	}

	return summaries, rows.Err()
}

func (repository *Repository) listWorkspaces(ctx context.Context, organizationID string) ([]Workspace, error) {
	rows, err := repository.db.QueryContext(ctx, `
SELECT id::text, organization_id::text, name, slug, created_at
FROM workspaces
WHERE organization_id = $1
ORDER BY created_at ASC`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workspaces := make([]Workspace, 0)
	for rows.Next() {
		var workspace Workspace
		if err := rows.Scan(&workspace.ID, &workspace.OrganizationID, &workspace.Name, &workspace.Slug, &workspace.CreatedAt); err != nil {
			return nil, err
		}
		workspaces = append(workspaces, workspace)
	}

	return workspaces, rows.Err()
}

type workspaceCreator interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func createWorkspace(ctx context.Context, creator workspaceCreator, organizationID string, name string) (Workspace, error) {
	workspaceName := strings.TrimSpace(name)
	var workspace Workspace
	err := creator.QueryRowContext(ctx, `
INSERT INTO workspaces (organization_id, name, slug)
VALUES ($1, $2, $3)
RETURNING id::text, organization_id::text, name, slug, created_at`, organizationID, workspaceName, slugify(workspaceName)).Scan(
		&workspace.ID,
		&workspace.OrganizationID,
		&workspace.Name,
		&workspace.Slug,
		&workspace.CreatedAt,
	)
	return workspace, err
}

var slugPattern = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(value string) string {
	slug := strings.ToLower(strings.TrimSpace(value))
	slug = slugPattern.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "workspace"
	}

	return slug
}
