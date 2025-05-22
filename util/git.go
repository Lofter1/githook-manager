package util

import "github.com/go-git/go-git/v5/plumbing"

func GetBranchNameFromRef(ref string) string {
	refName := plumbing.ReferenceName(ref)
	return refName.Short()
}
