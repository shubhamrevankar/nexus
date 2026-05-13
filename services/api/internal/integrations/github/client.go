package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

type apiTreeResponse struct {
	Tree []apiTreeEntry `json:"tree"`
}

type apiTreeEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	SHA  string `json:"sha"`
	Size int    `json:"size"`
}

type TreeFile struct {
	Path string
	SHA  string
	Size int
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

func (client *Client) FetchTreeFiles(ctx context.Context, token string, owner string, repo string, branch string) ([]TreeFile, error) {
	requestURL, err := url.JoinPath(client.baseURL, "repos", strings.TrimSpace(owner), strings.TrimSpace(repo), "git", "trees", strings.TrimSpace(branch))
	if err != nil {
		return nil, err
	}
	requestURL = requestURL + "?recursive=1"

	request, err := client.newGitHubRequest(ctx, http.MethodGet, requestURL, token, "application/vnd.github+json")
	if err != nil {
		return nil, err
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("github tree api returned status %d", response.StatusCode)
	}

	var payload apiTreeResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}

	files := make([]TreeFile, 0, len(payload.Tree))
	for _, entry := range payload.Tree {
		if entry.Type != "blob" {
			continue
		}
		files = append(files, TreeFile{Path: entry.Path, SHA: entry.SHA, Size: entry.Size})
	}

	return files, nil
}

func (client *Client) FetchRawFile(ctx context.Context, token string, owner string, repo string, path string, branch string) (string, error) {
	requestURL, err := url.JoinPath(client.baseURL, "repos", strings.TrimSpace(owner), strings.TrimSpace(repo), "contents", path)
	if err != nil {
		return "", err
	}
	requestURL = requestURL + "?ref=" + url.QueryEscape(strings.TrimSpace(branch))

	request, err := client.newGitHubRequest(ctx, http.MethodGet, requestURL, token, "application/vnd.github.raw")
	if err != nil {
		return "", err
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf("github file api returned status %d", response.StatusCode)
	}

	content, err := io.ReadAll(io.LimitReader(response.Body, 256*1024))
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (client *Client) newGitHubRequest(ctx context.Context, method string, requestURL string, token string, accept string) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, method, requestURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", accept)
	request.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))
	request.Header.Set("User-Agent", "nexus-local-development")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	return request, nil
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
