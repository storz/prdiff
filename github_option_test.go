package prdiff

import (
	"reflect"
	"testing"

	"github.com/google/go-github/v29/github"
)

func TestNewClosedPullRequestListOptions(t *testing.T) {
	type args struct {
		opts []PullRequestListOptionFunc
	}
	tests := []struct {
		name string
		args args
		want *github.PullRequestListOptions
	}{
		{
			name: "no argument",
			args: args{
				opts: nil,
			},
			want: &github.PullRequestListOptions{
				State:     "closed",
				Sort:      "updated",
				Direction: "desc",
				ListOptions: github.ListOptions{
					Page:    0,
					PerPage: 100,
				},
			},
		},
		{
			name: "update state",
			args: args{
				opts: []PullRequestListOptionFunc{
					func(opts *github.PullRequestListOptions) {
						opts.State = "open"
					},
				},
			},
			want: &github.PullRequestListOptions{
				State:     "open",
				Sort:      "updated",
				Direction: "desc",
				ListOptions: github.ListOptions{
					Page:    0,
					PerPage: 100,
				},
			},
		},
		{
			name: "TestWithBase",
			args: args{
				opts: []PullRequestListOptionFunc{WithBase("a")},
			},
			want: &github.PullRequestListOptions{
				State:     "closed",
				Sort:      "updated",
				Direction: "desc",
				Base:      "a",
				ListOptions: github.ListOptions{
					Page:    0,
					PerPage: 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClosedPullRequestListOptions(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClosedPullRequestListOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
