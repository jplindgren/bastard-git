package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestRepo(t *testing.T) *Repository {
	t.Helper()

	t.Setenv("BGIT_TEST_REPO", "test-repo")

	repo, err := Init()
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}

	return &repo
}

func teardownTestRepo(t *testing.T, repo *Repository) {
	t.Helper()

	// Clean up the temporary directory
	os.RemoveAll("test-repo")
}

func TestCreateBranch(t *testing.T) {
	repo := setupTestRepo(t)
	defer teardownTestRepo(t, repo)

	branchName := "feature/test-branch"
	err := repo.CreateBranch(branchName)
	if err != nil {
		t.Errorf("CreateBranch() error = %v, wantErr %v", err, false)
	}

	branches, err := repo.BranchList()
	if err != nil {
		t.Fatalf("Failed to list branches: %v", err)
	}

	branchCreated := false
	for _, branch := range branches {
		if branch == branchName {
			branchCreated = true
		}
	}

	if branchCreated {
		t.Fatalf(`Branches should not create head references until we have a commit, %v`, err)
	}

	// After create a commit, the branch is really created
	err = os.WriteFile(filepath.Join(repo.WorkTree, "test.txt"), []byte("test"), 0644)
	if err != nil {
		t.Fatalf(`Could not create new file, %v`, err)
	}

	err = repo.Commit("test commit")
	if err != nil {
		t.Fatalf(`Could not create commit, %v`, err)
	}

	branches, err = repo.BranchList()
	if err != nil {
		t.Fatalf("Failed to list branches: %v", err)
	}

	branchCreated = false
	for _, branch := range branches {
		if branch == branchName {
			branchCreated = true
		}
	}

	if !branchCreated {
		t.Fatalf(`Branch was not created, %v`, err)
	}
}
