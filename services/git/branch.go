package git

import (
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

func (r *GitConfig) GetBranchList() []string {
	iter, _ := r.repo.Branches()
	defer iter.Close()
	var branchesName []string
	for {
		ref, err := iter.Next()
		if err == io.EOF {
			break
		}
		branchesName = append(branchesName, ref.String())

	}
	return branchesName
}

func (r *GitConfig) GetRemoteBranches(askAuth func() transport.AuthMethod) ([]string, error) {
	var remoteBranches []string
	r.repo.Fetch(&git.FetchOptions{})

	// Get remote repository, by default is origin
	rem, err := r.repo.Remote("origin")
	if err != nil {
		return remoteBranches, err
	}

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		if err == transport.ErrAuthenticationRequired {
			// need authentication, ask user for auth method
			authMethod := askAuth()
			refs, err = rem.List(&git.ListOptions{
				Auth: authMethod,
			})
			if err != nil {
				return remoteBranches, err
			}
		} else {
			return remoteBranches, err
		}
	}

	for _, ref := range refs {
		if ref.Name().IsBranch() {
			remoteBranches = append(remoteBranches, ref.Name().Short())
		}
	}
	return remoteBranches, nil
}
