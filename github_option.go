package prdiff

import "github.com/google/go-github/v29/github"

type PullRequestListOptionFunc func(*github.PullRequestListOptions)

// NewClosedPullRequestListOptions generates *github.PullRequestListOptions for fetching closed utmost count PRs with update time desc.
func NewClosedPullRequestListOptions(opts ...PullRequestListOptionFunc) *github.PullRequestListOptions {
	prlo := &github.PullRequestListOptions{
		State:     "closed",
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 100,
		},
	}

	for _, opt := range opts {
		opt(prlo)
	}

	return prlo
}

// WithBase modifies Base field of *github.PullRequestListOptions.
func WithBase(base string) PullRequestListOptionFunc {
	return func(opts *github.PullRequestListOptions) {
		opts.Base = base
	}
}
