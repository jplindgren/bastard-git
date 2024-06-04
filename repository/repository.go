package repository

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jplindgren/bastard-git/object"
	"github.com/jplindgren/bastard-git/store"
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
	Email string //TODO: real git gets user/email from global/local config
	Store store.Store
}

func Init() error {
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
	currentRepository.Email = "email@gmail.com" //TODO: git get from global/local config
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

func (r *Repository) IsGitRepo() bool {
	_, err := os.Stat(r.Paths.bGitPath)
	return !os.IsNotExist(err)
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
