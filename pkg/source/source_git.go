package source

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitSource struct {
	RepoURL string
	Branch  string
}

func NewGitSource(url, branch string) *GitSource {
	return &GitSource{
		RepoURL: url,
		Branch:  branch,
	}
}

func (g *GitSource) Fetch(ctx context.Context, destDir string) error {
	_, err := git.PlainCloneContext(ctx, destDir, false, &git.CloneOptions{
		URL:           g.RepoURL,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + g.Branch),
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}
	return nil
}
