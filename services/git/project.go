package git

import (
	"fmt"
	"gitmonitor/services/utils"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

func (r *GitConfig) Fetch(askAuth func() transport.AuthMethod, defaultRemoteName string) error {
	const serviceName = "Fetch"

	// Fetch remote repository
	err := r.repo.Fetch(&git.FetchOptions{
		RemoteName: defaultRemoteName,
	})
	if err == transport.ErrAuthenticationRequired {
		r.auth = askAuth()
		err = r.repo.Fetch(&git.FetchOptions{
			RemoteName: defaultRemoteName,
			Auth:       r.auth,
		})
		if err != nil {
			utils.CheckErr(serviceName, err)
			return err
		}
	}
	return nil
}

func (r *GitConfig) Checkout(defaultBranchName string) error {
	const serviceName = "Fetch"
	w, err := r.repo.Worktree()
	utils.CheckErr(serviceName, err)

	// Checking out to origin/master
	return w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(defaultBranchName),
		Keep:   true,
	})
}

func (r *GitConfig) FetchAndCheckout(
	askAuth func() transport.AuthMethod,
	defaultBranchName string,
	defaultRemoteName string,
) error {
	const serviceName = "FetchAndCheckout"

	// Fetch remote repository
	err := r.repo.Fetch(&git.FetchOptions{
		RemoteName: defaultRemoteName,
	})
	if err == transport.ErrAuthenticationRequired {
		r.auth = askAuth()
		err = r.repo.Fetch(&git.FetchOptions{
			RemoteName: defaultRemoteName,
			Auth:       r.auth,
		})
		if err != nil {
			utils.CheckErr(serviceName, err)
			return err
		}
	}

	w, err := r.repo.Worktree()
	utils.CheckErr(serviceName, err)

	// Checking out to origin/master
	return w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(defaultBranchName),
		Keep:   true,
	})

}

func (r *GitConfig) GetOriginUrl() string {
	origin, err := r.repo.Remote("origin")
	if err != nil {
		return ""
	}
	return origin.Config().URLs[0]
}

func (r *GitConfig) GetPaths() ([]string, error) {
	refs, err := r.repo.Head()
	if err != nil {
		return []string{}, err
	}

	commit, err := r.repo.CommitObject(refs.Hash())
	if err != nil {
		return []string{}, err
	}
	tree, err := commit.Tree()
	if err != nil {
		return []string{}, err
	}

	var paths []string
	for _, entry := range tree.Entries {
		paths = append(paths, entry.Name)
	}

	return paths, nil
}

func (r *GitConfig) GetCommitObjects(
	remoteName string,
	defaultBranch string,
) ([]*object.Commit, error) {

	var cIter object.CommitIter
	var err error

	commitHash, err := r.repo.ResolveRevision(
		plumbing.Revision(fmt.Sprintf("refs/remotes/%s/%s", remoteName, defaultBranch)),
	)

	if err == nil {
		cIter, err = r.repo.Log(&git.LogOptions{
			From:  *commitHash,
			Order: git.LogOrderCommitterTime,
		})
	} else {
		cIter, err = r.repo.Log(&git.LogOptions{
			Order: git.LogOrderCommitterTime,
		})
	}

	if err != nil {
		return nil, err
	}
	var commits []*object.Commit
	cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})
	return commits, err
}
