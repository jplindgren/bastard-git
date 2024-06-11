package repository

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func setupTest(tb testing.TB) func(tb testing.TB) {
	tb.Setenv("BGIT_TEST_REPO", "test-repo")
	log.Println("setup suite")

	return func(tb testing.TB) {
		err := os.RemoveAll("test-repo")
		if err != nil {
			log.Fatalf(`Could not remove test repo, %v`, err)
		}
		log.Println("teardown test")
	}
}

func TestRepoInit(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	repo, err := Init()

	paths := []string{
		repo.paths.bGitPath,
		repo.paths.headPath,
		repo.paths.objectPath,
		repo.paths.refsPath,
		repo.paths.logs.path,
	}

	if err != nil {
		t.Fatalf(`Could not create repo, %v`, err)
	}

	for _, path := range paths {
		_, err := os.Stat(path)
		if err != nil {
			t.Fatalf(`Init repo did not created the path %s, %v`, path, err)
		}
	}
}

func TestCommit(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	repo, err := Init()
	if err != nil {
		t.Fatalf(`Could not create repo, %v`, err)
	}

	err = os.WriteFile(filepath.Join(repo.WorkTree, "test.txt"), []byte("test"), 0644)
	if err != nil {
		t.Fatalf(`Could not create new file, %v`, err)
	}

	err = repo.Commit("test commit")
	if err != nil {
		t.Fatalf(`Could not create commit, %v`, err)
	}

	cHash, tHash, err := repo.GetBranchTip("main")
	if err != nil {
		t.Fatalf(`Could not get branch tip, %v`, err)
	}

	commit, err := repo.Store.Get(cHash)
	if err != nil {
		t.Fatalf(`Could not get commit, %v`, err)
	}

	tree, err := repo.Store.Get(tHash)
	if err != nil {
		t.Fatalf(`Could not get commit tree, %v`, err)
	}

	t.Logf("commit: %v", string(commit))
	t.Logf("tree: %v", string(tree))
}
