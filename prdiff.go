package prdiff

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/google/go-github/v28/github"
)

type prdiff struct {
	gh            *GitHub
	Owner         string
	Repo          string
	DefaultBranch string
}

func New(gh *GitHub, owner, repo, defaultBranch string) *prdiff {
	return &prdiff{
		gh:            gh,
		Owner:         owner,
		Repo:          repo,
		DefaultBranch: defaultBranch,
	}
}

// UseRemoteDefaultBranch set automatically default branch on github.com.
func (t *prdiff) UseRemoteDefaultBranch(ctx context.Context) error {
	r, _, err := t.gh.GetRepository(ctx, t.Owner, t.Repo)
	if err != nil {
		return fmt.Errorf("failed to get repository: %v\n", err)
	}
	t.DefaultBranch = r.GetDefaultBranch()
	return nil
}

// GetLatestRelease returns the latest published release.
func (t *prdiff) GetLatestRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	rr, _, err := t.gh.GetRepositoryLatestRelease(ctx, t.Owner, t.Repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository latest release: %v\n", err)
	}
	return rr, nil
}

// GetDiffPullRequest returns the list of pull request these were merged between base and head.
func (t *prdiff) GetDiffPullRequests(ctx context.Context, base, head string) ([]*github.PullRequest, error) {
	cc, _, err := t.gh.CompareRepositoryCommits(ctx, t.Owner, t.Repo, base, head)
	if err != nil {
		return nil, fmt.Errorf("failed to compare repository commits: %v\n", err)
	}
	if len(cc.Commits) == 0 {
		return nil, fmt.Errorf("No commits between %s and %s\n", base, head)
	}

	prs, err := t.getPullRequestsMergedBetween(
		ctx,
		cc.GetBaseCommit().GetCommit().GetCommitter().GetDate(),
		time.Now(),
		NewClosedPullRequestListOptions(WithBase(t.DefaultBranch)))
	if err != nil {
		return nil, fmt.Errorf("failed to get PullRequests: %v\n", err)
	}

	// Drop pull request if it equals base because `MergedAt` is after commit timestamp of merge.
	nprs := make([]*github.PullRequest, 0, len(prs))
	for _, pr := range prs {
		if pr.GetMergeCommitSHA() != base {
			nprs = append(nprs, pr)
		}
	}

	sort.Slice(nprs, func(i, j int) bool {
		return nprs[i].GetNumber() < nprs[j].GetNumber()
	})

	return nprs, nil
}

func (t *prdiff) getPullRequestsMergedBetween(ctx context.Context, t1 time.Time, t2 time.Time, opt *github.PullRequestListOptions) ([]*github.PullRequest, error) {
	if t1.After(t2) {
		t2, t1 = t1, t2
	}

	nprs := make([]*github.PullRequest, 0)
	for {
		prs, res, err := t.gh.ListPullRequests(ctx, t.Owner, t.Repo, opt)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("invalid status")
		}

		var added int
		for _, pr := range prs {
			if mergedAt := pr.GetMergedAt(); !mergedAt.IsZero() {
				if mergedAt.Equal(t1) || mergedAt.After(t1) && mergedAt.Before(t2) { // [t1, t2)
					nprs = append(nprs, pr)
					added++
				}
			}
		}

		if added == 0 && prs[len(prs)-1].GetUpdatedAt().Before(t1) {
			break
		}
		opt.Page = res.NextPage
	}

	return nprs, nil
}
