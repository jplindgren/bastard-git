package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jplindgren/bastard-git/object"
)

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Please provide a command")
		os.Exit(1)
	}

	command := args[0]
	funcToExecute, err := getCommandsToExecute(command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//fmt.Println(args[1:])

	err = funcToExecute(args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// # Internal

func getCommandsToExecute(command string) (func(args []string) error, error) {
	mapper := map[string]func(args []string) error{
		"init": commandInit,
		//"hash-object": commandHashObject,
		"status":   commandStatus,
		"cat-file": commandCatFile,
		"commit":   commandCommit,
		"checkout": commandCheckout,
		"branch":   commandBranch,
		"branches": commandBranchList,
		"test":     commandTest,
	}

	if val, ok := mapper[command]; ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("command %s not found", command)
	}
}

// # Commands

func commandInit(args []string) error {
	return InitRepository()
}

// func commandHashObject(args []string) error {
// 	repository := GetRepository()

// 	fileName := args[0]
// 	objType := "blob" //default
// 	if len(args) > 1 {
// 		objType = args[1]
// 	}

// 	_, err := repository.CreateObject(objType, fileName)
// 	return err
// }

func commandCatFile(args []string) error {
	repo := GetRepository()
	if !repo.IsGitRepo() {
		fmt.Println("Not a git repository")
		os.Exit(1)
	}

	hash := args[0]

	content, err := repo.store.Get(hash)
	os.Stdout.Write([]byte(content))
	return err
}

func commandStatus(args []string) error {
	repo := GetRepository()
	if !repo.IsGitRepo() {
		fmt.Println("Not a git repository")
		os.Exit(1)
	}

	fmt.Println("Getting status")

	reader := &FileSystemReader{}
	toBeCommited, err := reader.Diff(repo.workTree)
	if err != nil {
		return err
	}

	if len(toBeCommited) == 0 {
		fmt.Println("Nothing to commit")
	} else {
		fmt.Println("Changes to be commited:")
		for _, file := range toBeCommited {
			fmt.Printf("%s:   %s\n", file.status, file.name)
		}
	}

	return nil
}

func commandCommit(args []string) error {
	repo := GetRepository()
	if !repo.IsGitRepo() {
		fmt.Println("Not a git repository")
		os.Exit(1)
	}

	message := args[0]

	if message == "" {
		fmt.Println("Please provide a commit message")
		os.Exit(1)
	}

	fmt.Println("Creating commit with message: " + message)

	rootTree, err := GenerateObjectTree(repo.workTree)
	if err != nil {
		return err
	}

	// //TODO: commit should check if a previous commit exists and set it as parent
	commit := object.NewCommit(repo.email, message, rootTree.GetHash(), []string{})
	err = repo.store.Write(commit)
	if err != nil {
		return err
	}

	repo.updateIndex(rootTree)
	repo.updateRefHead(commit.ToString())
	return err
}

func commandCheckout(args []string) error {
	repo := GetRepository()
	if !repo.IsGitRepo() {
		fmt.Println("Not a git repository")
		os.Exit(1)
	}

	if len(args) == 0 {
		fmt.Println("Please provide an argument")
		os.Exit(1)
	}

	if len(args) == 2 && args[0] != "-b" {
		fmt.Println("Please provide -b as argument")
		os.Exit(1)
	}

	if len(args) == 1 {
		//get last tree
		_, treeHash, err := repo.GetBranchTip(args[0])
		if err != nil {
			return err
		}

		err = RecreateWorkingTree(treeHash)
		if err != nil {
			return err
		}

		repo.setHEAD(args[0])

		os.Stdout.WriteString("Switched to branch " + args[0])
	}

	if len(args) == 2 {
		err := repo.createBranch(args[1])
		if err != nil {
			return err
		}

		os.Stdout.WriteString("New branch created.\nSwitched to a new branch " + args[1])
	}
	return nil
}

func commandBranch(args []string) error {
	repo := GetRepository()
	if !repo.IsGitRepo() {
		fmt.Println("Not a git repository")
		os.Exit(1)
	}

	branch, err := repo.GetCurrentBranch()
	if err != nil {
		return err
	}

	os.Stdout.WriteString(branch)
	return nil
}

func commandBranchList(args []string) error {
	repo := GetRepository()

	branches, err := repo.BranchList()
	if err != nil {
		return err
	}

	bText := strings.Join(branches, "\n")
	os.Stdout.WriteString(bText)
	return nil
}

func commandTest(args []string) error {
	//err := DeleteWorkingTree()

	err := RecreateWorkingTree(args[0])
	return err
}

//05ab1b66b4456f6ecd3d40314b32b7763eaaba57
