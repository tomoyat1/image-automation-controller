package git

import (
	"context"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	"github.com/fluxcd/source-controller/pkg/git"
	gitstrat "github.com/fluxcd/source-controller/pkg/git/strategy"
	libgit2 "github.com/libgit2/git2go/v33"
)

type RepositoryManager struct {
	repos map[string]Repository
}

func (r *RepositoryManager) CloneInto(ctx context.Context,
	access RepoAccess, ref *sourcev1.GitRepositoryRef, path string) (*libgit2.Repository, error) {
	if repo, ok := r.repos[access.url]; ok {
		// Use the path of the known clone
		path = repo.Path
	}
	opts := git.CheckoutOptions{}
	if ref != nil {
		opts.Tag = ref.Tag
		opts.SemVer = ref.SemVer
		opts.Commit = ref.Commit
		opts.Branch = ref.Branch
	}
	checkoutStrat, err := gitstrat.CheckoutStrategyForImplementation(ctx, sourcev1.LibGit2Implementation, opts)
	if err == nil {
		_, err = checkoutStrat.Checkout(ctx, path, access.url, access.auth) // Checkout is idempotent
	}
	if err != nil {
		return nil, err
	}

	r.repos[access.url] = Repository{
		Path: path,
		URL:  access.url,
	}
	return libgit2.OpenRepository(path)
}

// Repository is a repository cloned and managed for use across reconciliations.
// TODO: decide what information we need to hold in a "Repository".
type Repository struct {
	Path string
	URL  string // TODO do we really need this field?
}

type RepoAccess struct {
	auth *git.AuthOptions
	url  string
}
