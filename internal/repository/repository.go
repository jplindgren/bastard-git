package repository

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jplindgren/bastard-git/internal/store"
)

var initialBranch = "main"

type Repository struct {
	WorkTree       string
	BGitFolder     string
	bGitTempFolder string
	paths          struct {
		bGitTempPath   string
		bGitPath       string
		headPath       string
		IndexPath      string
		objectPath     string
		refsPath       string
		bGitIgnorePath string
		logs           struct {
			path string
			head string
		}
	}
	Email string
	Store store.Store
}

func Init() (Repository, error) {
	currentRepository := GetRepository("")

	err := currentRepository.createGitInfra()
	return currentRepository, err
}

var bGitBinaryName = "bgit"
var bGitFolder = ".bgit"
var bGitLogsFolder = "logs"
var bGitTempFolder = ".bGitemp"
var bGitIgnoreFile = ".bgitignore"
var BGitRefsHeads = "refs/heads"

func GetRepository(user string) Repository {
	rooPath, _ := os.Getwd()

	testRepo := os.Getenv("BGIT_TEST_REPO") //can be used to tesst bgit commands

	currentRepository := Repository{}
	currentRepository.Email = user
	currentRepository.WorkTree = filepath.Join(rooPath, testRepo)
	currentRepository.bGitTempFolder = bGitTempFolder
	currentRepository.BGitFolder = bGitFolder

	bGitFolderPath := filepath.Join(currentRepository.WorkTree, bGitFolder)

	currentRepository.paths.bGitTempPath = filepath.Join(currentRepository.WorkTree, bGitTempFolder)
	currentRepository.paths.bGitPath = bGitFolderPath
	currentRepository.paths.headPath = filepath.Join(bGitFolderPath, "HEAD")
	currentRepository.paths.IndexPath = filepath.Join(bGitFolderPath, "index")
	currentRepository.paths.objectPath = filepath.Join(bGitFolderPath, "objects")
	currentRepository.paths.refsPath = filepath.Join(bGitFolderPath, BGitRefsHeads)
	currentRepository.paths.logs = struct {
		path string
		head string
	}{
		path: filepath.Join(bGitFolderPath, bGitLogsFolder),
		head: filepath.Join(bGitFolderPath, bGitLogsFolder, "HEAD"),
	}
	currentRepository.paths.bGitIgnorePath = filepath.Join(currentRepository.WorkTree, bGitIgnoreFile)
	currentRepository.Store = store.New(currentRepository.paths.objectPath)
	return currentRepository
}

func (r *Repository) IsValid() bool {
	_, err := os.Stat(r.paths.bGitPath)
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
	err := os.MkdirAll(r.paths.objectPath, fs.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(r.paths.refsPath, fs.ModePerm)
	if err != nil {
		return err
	}
	r.SetHead(initialBranch)

	err = os.MkdirAll(r.paths.logs.path, fs.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(r.paths.logs.path, BGitRefsHeads), fs.ModePerm)
	if err != nil {
		return err
	}

	fmt.Printf("Initialized empty Git repository in %s/.bgit/ \n", r.WorkTree)
	return nil
}

func (r *Repository) SetHead(branch string) string {
	branchRefHead := filepath.Join(BGitRefsHeads, branch)
	err := os.WriteFile(r.paths.headPath, []byte(branchRefHead), 0644)
	if err != nil {
		fmt.Printf("Error setting HEAD: %s \n", err)
		os.Exit(1)
	}
	return branchRefHead
}

func (r *Repository) getHeadRef() (string, string, error) {
	content, err := os.ReadFile(r.paths.headPath)
	if err != nil {
		return "", "", err
	}
	ref := string(content)
	return ref, filepath.Join(r.paths.bGitPath, ref), nil
}

func (r *Repository) updateRefHead(commit string, refHead string) error {
	headRefPath := filepath.Join(r.paths.bGitPath, string(refHead))
	if err := createPathFoldersIfNotExists(headRefPath); err != nil {
		return err
	}

	file, err := os.OpenFile(headRefPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(commit)
	return err
}

func (r *Repository) getIgnoredPaths() (map[string]bool, error) {
	ignoredPaths := make(map[string]bool)

	// Add .bgitignore file to the ignored paths
	ignoredPaths[r.paths.bGitIgnorePath] = true
	ignoredPaths[filepath.Join(r.WorkTree, bGitBinaryName)] = true
	ignoredPaths[filepath.Join(r.WorkTree, bGitBinaryName+".exe")] = true

	file, err := os.Open(r.paths.bGitIgnorePath)
	if err != nil {
		return ignoredPaths, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			continue
		}
		path := scanner.Text()
		if path == "" {
			continue
		}

		absPath := filepath.Join(r.WorkTree, path)
		ignoredPaths[absPath] = true
	}
	return ignoredPaths, nil
}

func checkIfFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func createPathFoldersIfNotExists(path string) error {
	if !checkIfFileExists(path) {
		dirPath := filepath.Dir(path)
		err := os.MkdirAll(dirPath, fs.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
