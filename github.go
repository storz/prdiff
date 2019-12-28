package prdiff

import (
	"context"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

// GitHubAPI represents required API calls for GitHub.
type GitHubAPI interface {
	GetRepository(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error)
	GetRepositoryLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error)
	ListPullRequests(ctx context.Context, owner, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
	CompareRepositoryCommits(ctx context.Context, owner, repo, base, head string) (*github.CommitsComparison, *github.Response, error)
}

// GitHub manages GitHub API client from google/go-github.
type GitHub struct {
	client *github.Client
}

var _ GitHubAPI = (*GitHub)(nil)

// NewGitHub returns `GitHub` with OAuth access token of github.com.
func NewGitHub(ctx context.Context, accessToken string) *GitHub {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &GitHub{
		client: github.NewClient(tc),
	}
}

// GetRepository calls `Repositories.Get`.
func (gh *GitHub) GetRepository(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error) {
	return gh.client.Repositories.Get(ctx, owner, repo)
}

// GetRepositoryLatestRelease calls `Repositories.GetLatestRelease`.
func (gh *GitHub) GetRepositoryLatestRelease(ctx context.Context, owner, repo string) (*github.RepositoryRelease, *github.Response, error) {
	return gh.client.Repositories.GetLatestRelease(ctx, owner, repo)
}

// ListPullRequests calls `PullRequests.List`.
func (gh *GitHub) ListPullRequests(ctx context.Context, owner, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
	return gh.client.PullRequests.List(ctx, owner, repo, opt)
}

// CompareRepositoryCommits calls `Repositories.CompareCommits`.
func (gh *GitHub) CompareRepositoryCommits(ctx context.Context, owner, repo, base, head string) (*github.CommitsComparison, *github.Response, error) {
	return gh.client.Repositories.CompareCommits(ctx, owner, repo, base, head)
}
