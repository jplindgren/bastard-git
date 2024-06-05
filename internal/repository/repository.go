package repository

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jplindgren/bastard-git/internal/object"
	"github.com/jplindgren/bastard-git/internal/store"
)

var initialBranch = "main"

type Repository struct {
	WorkTree       string
	BGitFolder     string
	bGitTempFolder string
	Paths          struct {
		bGitTempPath string
		bGitPath     string
		headPath     string
		IndexPath    string
		objectPath   string
		refsPath     string
	}
	Email string
	Store store.Store
}

func Init() error {
	currentRepository := GetRepository("")

	err := currentRepository.createGitInfra()
	return err
}

var bGitFolder = ".bgit"
var bGitTempFolder = ".bGitemp"

func GetRepository(user string) Repository {
	rooPath, _ := os.Getwd()

	testFolder := "srctest"

	currentRepository := Repository{}
	currentRepository.Email = user
	currentRepository.WorkTree = filepath.Join(rooPath, testFolder)
	currentRepository.bGitTempFolder = bGitTempFolder
	currentRepository.BGitFolder = bGitFolder

	bGitFoderPath := filepath.Join(currentRepository.WorkTree, bGitFolder)

	currentRepository.Paths.bGitTempPath = filepath.Join(currentRepository.WorkTree, bGitTempFolder)
	currentRepository.Paths.bGitPath = bGitFoderPath
	currentRepository.Paths.headPath = filepath.Join(bGitFoderPath, "HEAD")
	currentRepository.Paths.IndexPath = filepath.Join(bGitFoderPath, "index")
	currentRepository.Paths.objectPath = filepath.Join(bGitFoderPath, "objects")
	currentRepository.Paths.refsPath = filepath.Join(bGitFoderPath, "refs/heads")
	currentRepository.Store = store.New(currentRepository.Paths.objectPath)
	return currentRepository
}

func (r *Repository) IsValid() bool {
	_, err := os.Stat(r.Paths.bGitPath)
	if os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Not a git repository")
		return false
	}

	if r.Email == "" {
		fmt.Fprintln(os.Stderr, "User not set. Please use 'export BGIT_USER=<email>'")
		return false
	}

	return true
}

func (r *Repository) createGitInfra() error {
	err := os.MkdirAll(r.Paths.objectPath, fs.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(r.Paths.refsPath, fs.ModePerm)
	if err != nil {
		fmt.Printf("Error setting HEAD: %s \n", err)
		os.Exit(1)
	}
	r.SetHEAD(initialBranch)

	fmt.Printf("Initialized empty Git repository in %s/.bgit/ \n", r.WorkTree)
	return nil
}

func (r *Repository) SetHEAD(branch string) {
	err := os.WriteFile(r.Paths.headPath, []byte(fmt.Sprintf("refs/heads/%s", branch)), 0644)
	if err != nil {
		fmt.Printf("Error setting HEAD: %s \n", err)
		os.Exit(1)
	}
}

func (r *Repository) UpdateRefHead(commit string) {
	content, err := os.ReadFile(r.Paths.headPath) //extract?
	if err != nil {
		fmt.Printf("Error reading HEAD %s \n", err)
		os.Exit(1)
	}

	path := string(content)
	completePath := filepath.Join(r.Paths.bGitPath, path)
	err = os.WriteFile(completePath, []byte(commit), 0644)
	if err != nil {
		fmt.Printf("Error updating the ref %s \n", err)
		os.Exit(1)
	}
}

func (r *Repository) UpdateIndex(rootTree object.BGitObject) error {
	data := rootTree.FormatToIndex()
	return os.WriteFile(r.Paths.IndexPath, []byte(data), 0644)
}

func (r *Repository) SetIndex(content []byte) error {
	return os.WriteFile(r.Paths.IndexPath, content, 0644)
}
