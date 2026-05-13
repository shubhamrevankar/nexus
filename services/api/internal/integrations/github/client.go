package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type apiRepository struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	FullName      string     `json:"full_name"`
	Description   *string    `json:"description"`
	HTMLURL       string     `json:"html_url"`
	DefaultBranch string     `json:"default_branch"`
	Language      *string    `json:"language"`
	Private       bool       `json:"private"`
	Stars         int        `json:"stargazers_count"`
	Forks         int        `json:"forks_count"`
	OpenIssues    int        `json:"open_issues_count"`
	PushedAt      *time.Time `json:"pushed_at"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (client *Client) FetchRepository(ctx context.Context, token string, owner string, repo string) (Repository, error) {
	requestURL, err := url.JoinPath(client.baseURL, "repos", strings.TrimSpace(owner), strings.TrimSpace(repo))
	if err != nil {
		return Repository{}, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return Repository{}, err
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))
	request.Header.Set("User-Agent", "nexus-local-development")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	response, err := client.httpClient.Do(request)
	if err != nil {
		return Repository{}, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return Repository{}, fmt.Errorf("github api returned status %d", response.StatusCode)
	}

	var payload apiRepository
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return Repository{}, err
	}

	return Repository{
		GitHubID:      payload.ID,
		Owner:         strings.TrimSpace(owner),
		Name:          payload.Name,
		FullName:      payload.FullName,
		Description:   stringValue(payload.Description),
		HTMLURL:       payload.HTMLURL,
		DefaultBranch: payload.DefaultBranch,
		Language:      stringValue(payload.Language),
		Private:       payload.Private,
		Stars:         payload.Stars,
		Forks:         payload.Forks,
		OpenIssues:    payload.OpenIssues,
		PushedAt:      payload.PushedAt,
	}, nil
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
