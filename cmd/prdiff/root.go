package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/storz/prdiff"
)

var (
	optToken, optOwner, optRepo, optTargetBranch string
)

const (
	diffSeparator = "..."
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prdiff [[base] head] [base...head]",
		Short: "prdiff enumerates GitHub Pull Requests merged into default branch since last release (or between two branches/commits/tags)",
		Args:  cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showDiff(parseArgs(args))
		},
	}
	cmd.PersistentFlags().StringVar(&optToken, "token", "", "GitHub access token (default: $GITHUB_TOKEN)")
	cmd.PersistentFlags().StringVarP(&optOwner, "owner", "o", "", "user/organization name (default: repository of current directory)")
	cmd.PersistentFlags().StringVarP(&optRepo, "repo", "r", "", "repository name (default: repository of current directory)")
	cmd.Flags().StringVarP(&optTargetBranch, "target", "t", "", "specified target branch what you want to get diff (default: default branch set on GitHub)")

	return cmd
}

func parseArgs(args []string) (string, string) {
	var base, head string
	switch len(args) {
	case 0:
		// nop
	case 1:
		if nargs := strings.Split(args[0], diffSeparator); len(nargs) > 1 {
			base = nargs[0]
			head = nargs[1]
		} else {
			head = args[0]
		}
	default:
		base = args[0]
		head = args[1]
	}
	return base, head
}

func detectRepository(owner, repo string) (string, string, error) {
	if owner != "" && repo != "" {
		return owner, repo, nil
	}

	var dir string
	rev, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		wd, err := os.Getwd()
		if err != nil {
			return owner, repo, fmt.Errorf("failed to get git repository and current working directory")
		}
		dir = wd
	} else {
		dir = string(rev)
	}

	if owner == "" {
		mayOwner, err := exec.Command("dirname", dir).Output()
		if err != nil {
			return "", "", err
		}
		paths := strings.Split(string(mayOwner), string(os.PathSeparator))
		owner = strings.TrimSpace(paths[len(paths)-1])
	}

	if repo == "" {
		mayRepo, err := exec.Command("basename", dir).Output()
		if err != nil {
			return "", "", err
		}
		repo = strings.TrimSpace(string(mayRepo))
	}

	if owner == "" || repo == "" {
		return "", "", fmt.Errorf("owner/repo must be specified")
	}
	return owner, repo, nil
}

func showDiff(base, head string) error {
	token := optToken
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	owner, repo, err := detectRepository(optOwner, optRepo)
	if err != nil {
		return err
	}
	fmt.Printf("repository: %s/%s\n", owner, repo)

	ctx := context.Background()
	gitHub := prdiff.NewGitHub(ctx, token)

	pd := prdiff.New(gitHub, owner, repo, optTargetBranch)
	if optTargetBranch == "" {
		pd.UseRemoteDefaultBranch(ctx)
	}

	if base == "" {
		latestRelease, err := pd.GetLatestRelease(ctx)
		if err != nil {
			return err
		}
		base = latestRelease.GetTagName()
	}
	if head == "" {
		head = pd.DefaultBranch
	}
	fmt.Printf("base: %s <- head: %s\n", base, head)

	prs, err := pd.GetDiffPullRequests(ctx, base, head)
	if err != nil {
		return err
	}

	rows := make([]string, len(prs))
	for i, pr := range prs {
		s := fmt.Sprintf("#%d %s by @%s", pr.GetNumber(), pr.GetTitle(), pr.GetUser().GetLogin())
		rows[i] = s
	}
	fmt.Printf("%s\n", strings.Join(rows, "\n"))

	return nil
}
