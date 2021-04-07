package services

import (
	"io"

	"github.com/go-git/go-git/v5"
)

func getBranchList(r *git.Repository) []string {
	iter, _ := r.Branches()
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
