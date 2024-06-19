package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (r *Repository) CreateBranch(newBranch string) error {
	//copy from last branch
	lastCommit, currentBranchRef, err := r.lookupLastCommit()
	if err != nil {
		return err
	}

	newBranchRef := r.SetHead(newBranch)

	//in case of no commits yet, we do not update refHead and logs
	if lastCommit == "" {
		return nil
	}

	r.updateRefHead(lastCommit, newBranchRef)
	r.logHead(CHECKOUT, lastCommit, lastCommit, time.Now(), fmt.Sprintf("moving from %s to %s", currentBranchRef, newBranchRef))

	//no need to update the index, since new branches points to the same commit
	return nil
}

func (r *Repository) BranchList() ([]string, error) {
	dirEntries, err := os.ReadDir(r.paths.refsPath)
	if err != nil {
		return nil, err
	}

	var branches []string
	for _, entry := range dirEntries {
		if entry.IsDir() {
			innerEntries, err := os.ReadDir(filepath.Join(r.paths.refsPath, entry.Name()))
			if err != nil {
				return nil, err
			}
			branches = append(branches, filepath.Join(entry.Name(), innerEntries[0].Name()))
		} else {
			branches = append(branches, entry.Name())
		}
	}

	return branches, nil
}

// returns the last commit and tree hashes
func (r *Repository) GetBranchTip(branchRef string) (string, string, error) {
	// In case the client sends only the branch name, we append the refs/heads prefix
	if !strings.Contains(branchRef, BGitRefsHeads) {
		branchRef = filepath.Join(BGitRefsHeads, branchRef)
	}

	bCommitHash, err := os.ReadFile(filepath.Join(r.paths.bGitPath, branchRef))
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

// returns the current branch name and the full reference path
func (r *Repository) GetCurrentBranch() (string, string, error) {
	content, err := os.ReadFile(r.paths.headPath)
	if err != nil {
		return "", "", err
	}

	sContent := string(content)
	return strings.TrimPrefix(sContent, BGitRefsHeads), sContent, nil
}
