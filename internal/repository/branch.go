package repository

import (
	"os"
	"path/filepath"
	"strings"
)

func (r *Repository) CreateBranch(branch string) error {
	//copy from last branch
	branchRefHead, err := os.ReadFile(r.Paths.headPath)
	if err != nil {
		return err
	}

	lastCommit, err := os.ReadFile(filepath.Join(r.Paths.bGitPath, string(branchRefHead)))
	if err != nil {
		return err
	}

	r.SetHEAD(branch)
	r.UpdateRefHead(string(lastCommit))

	//no need to update the index, since new branches points to the same commit
	return nil
}

func (r *Repository) BranchList() ([]string, error) {
	dirEntries, err := os.ReadDir(r.Paths.refsPath)
	if err != nil {
		return nil, err
	}

	var branches []string
	for _, entry := range dirEntries {
		branches = append(branches, entry.Name())
	}

	return branches, nil
}

func (r *Repository) GetBranchTip(branch string) (string, string, error) {
	bCommitHash, err := os.ReadFile(filepath.Join(r.Paths.refsPath, branch))
	if err != nil {
		if os.IsNotExist(err) { //no commits yet
			return "", "", nil
		}
		return "", "", err
	}

	commitHash := string(bCommitHash)
	commitContent, err := r.Store.Get(commitHash)
	if err != nil {
		return "", "", err
	}

	if commitContent == nil {
		return "", "", nil
	}

	lines := strings.Split(string(commitContent), "\n") //TODO: do not read all file
	data := strings.Split(lines[0], " ")
	treeHash := data[1]
	return commitHash, treeHash, nil
}

func (r *Repository) GetCurrentBranch() (string, error) {
	content, err := os.ReadFile(r.Paths.headPath)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(string(content), "refs/heads/"), nil
}
