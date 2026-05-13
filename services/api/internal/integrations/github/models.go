package github

import "time"

type Repository struct {
	ID            string     `json:"id"`
	WorkspaceID   string     `json:"workspaceId"`
	GitHubID      int64      `json:"githubId"`
	Owner         string     `json:"owner"`
	Name          string     `json:"name"`
	FullName      string     `json:"fullName"`
	Description   string     `json:"description"`
	HTMLURL       string     `json:"htmlUrl"`
	DefaultBranch string     `json:"defaultBranch"`
	Language      string     `json:"language"`
	Private       bool       `json:"private"`
	Stars         int        `json:"stars"`
	Forks         int        `json:"forks"`
	OpenIssues    int        `json:"openIssues"`
	PushedAt      *time.Time `json:"pushedAt"`
	SyncedAt      time.Time  `json:"syncedAt"`
}
