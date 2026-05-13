package tenancy

import "time"

type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type Workspace struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organizationId"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	CreatedAt      time.Time `json:"createdAt"`
}

type WorkspaceSummary struct {
	Organization Organization `json:"organization"`
	Workspaces   []Workspace  `json:"workspaces"`
}
