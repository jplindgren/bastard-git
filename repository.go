package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jplindgren/bastard-git/object"
	"github.com/jplindgren/bastard-git/store"
)

var initialBranch = "main"

type Repository struct {
	workTree       string
	bGitFolder     string
	bGitTempFolder string
	paths          struct {
		bGitTempPath string
		bGitPath     string
		headPath     string
		indexPath    string
		objectPath   string
		refsPath     string
	}
	email string //TODO: real git gets user/email from global/local config
	store store.Store
}

func InitRepository() error {
	currentRepository := GetRepository()

	err := currentRepository.createGitInfra()
	return err
}

var bGitFolder = ".bgit"
var bGitTempFolder = ".bGitemp"

func GetRepository() Repository {
	rooPath, _ := os.Getwd()

	testFolder := "srctest"

	currentRepository := Repository{}
	currentRepository.email = "email@gmail.com" //TODO: git get from global/local config
	currentRepository.workTree = filepath.Join(rooPath, testFolder)
	currentRepository.bGitTempFolder = bGitTempFolder
	currentRepository.bGitFolder = bGitFolder

	bGitFoderPath := filepath.Join(currentRepository.workTree, bGitFolder)

	currentRepository.paths.bGitTempPath = filepath.Join(currentRepository.workTree, bGitTempFolder)
	currentRepository.paths.bGitPath = bGitFoderPath
	currentRepository.paths.headPath = filepath.Join(bGitFoderPath, "HEAD")
	currentRepository.paths.indexPath = filepath.Join(bGitFoderPath, "index")
	currentRepository.paths.objectPath = filepath.Join(bGitFoderPath, "objects")
	currentRepository.paths.refsPath = filepath.Join(bGitFoderPath, "refs/heads")
	currentRepository.store = store.New(currentRepository.paths.objectPath)
	return currentRepository
}

func (r *Repository) IsGitRepo() bool {
	_, err := os.Stat(r.paths.bGitPath)
	return !os.IsNotExist(err)
}

func (r *Repository) createGitInfra() error {
	err := os.MkdirAll(r.paths.objectPath, fs.ModePerm)
	if err != nil {
		return err
	}
	r.setHEAD(initialBranch)

	fmt.Printf("Initialized empty Git repository in %s/.bgit/ \n", r.workTree)
	return nil
}

func (r *Repository) setHEAD(branch string) {
	err := os.WriteFile(r.paths.headPath, []byte(fmt.Sprintf("refs/heads/%s", branch)), 0644)
	if err != nil {
		fmt.Printf("Error setting HEAD: %s \n", err)
		os.Exit(1)
	}

	err = os.MkdirAll(r.paths.refsPath, fs.ModePerm)
	if err != nil {
		fmt.Printf("Error setting HEAD: %s \n", err)
		os.Exit(1)
	}
}

func (r *Repository) updateRefHead(commit string) {
	content, err := os.ReadFile(r.paths.headPath) //extract?
	if err != nil {
		fmt.Printf("Error reading HEAD %s \n", err)
		os.Exit(1)
	}

	path := string(content)
	completePath := filepath.Join(r.paths.bGitPath, path)
	err = os.WriteFile(completePath, []byte(commit), 0644)
	if err != nil {
		fmt.Printf("Error updating the ref %s \n", err)
		os.Exit(1)
	}
}

func (r *Repository) updateIndex(rootTree object.BGitObject) error {
	data := rootTree.FormatToIndex()
	return os.WriteFile(r.paths.indexPath, []byte(data), 0644)
}

func (r *Repository) SetIndex(content []byte) error {
	return os.WriteFile(r.paths.indexPath, content, 0644)
}

func (r *Repository) createBranch(branch string) error {
	//copy from last branch
	branchRefHead, err := os.ReadFile(r.paths.headPath)
	if err != nil {
		return err
	}

	lastCommit, err := os.ReadFile(filepath.Join(r.paths.bGitPath, string(branchRefHead)))
	if err != nil {
		return err
	}

	r.setHEAD(branch)
	r.updateRefHead(string(lastCommit))

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
		branches = append(branches, entry.Name())
	}

	return branches, nil
}

func (r *Repository) GetBranchTip(branch string) (string, string, error) {
	bCommitHash, err := os.ReadFile(filepath.Join(r.paths.refsPath, branch))
	if err != nil {
		return "", "", err
	}

	commitHash := string(bCommitHash)
	commitContent, err := r.store.Get(commitHash)
	if err != nil {
		return "", "", err
	}

	lines := strings.Split(string(commitContent), "\n") //TODO: do not read all file
	data := strings.Split(lines[0], " ")
	treeHash := data[1]
	return commitHash, treeHash, nil
}

func (r *Repository) GetCurrentBranch() (string, error) {
	content, err := os.ReadFile(r.paths.headPath)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(string(content), "refs/heads/"), nil
}
