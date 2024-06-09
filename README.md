# bastard-git

Making a simple git in GO

## Authors

- [Joao Lindgren](https://github.com/jplindgren)
- Add your clickable username here. It should point to your GitHub account.

## Summary how it works

**Bastard git** is a CLI tool pet project which mimics a small set of features from GIT. Despite covering only a small part of GIT's functionalities, the basis of how the version controller works was maintained. The idea is to simulate as much as possible what we have in GIT. Head, Index, Objects, etc.

As GIT does, bastard git stores the entire content/at each time under the `./bgit` folder. When you commit file(s) **BGIT** will create a tree of objects of 3 types.
A commit, one or more trees, and one or more blobs. The ref/head is updated to point to the recently commit. As it happens with GIT, each commit saves the entire state of the repo not the diff, and differently from GIT we do not have especial treatments for big files, so take care with them.
**BGIT** should be able to reconstruct the working tree based on any commit at any time. In theory, for instance, clone the repo is just copy the `./bgit` folder and reconstruct the working tree from the index.

## Pieces

### Objects

Objects are saved under the `/objects` folder with their content compressed. Each object has a hash, and the first two parts of the hash are used to create a folder and the rest is the name of the saved file. `store.go` is responsible to save/retrieve data from there. The command `bgit cat-file <hash>` can be used to retrieve the object content decompressed.

- A commit has a parent commit(unless it is the first), a root tree, and author and a time. The commit hash is the message + author + time
- The tree represents directories. Each tree (dir) can contains other trees or blobs. The tree has a description of each tree/blob along with filemode and hashs.
- The blob represents a file at a specific point in time. The whole content is saved compressed and the hash is based on the content of the file.

### HEAD

File pointing to the current branch

### Refs/Heads

Each file represents one branch, and the branch is just a pointer to a commit (the content of the file is a commit hash)

### Index

File containing the state of the repo at the time. In real GIT the index is updated when the user adds a file to stage. **BGIT** only does that when the user commits. We can get the diff of files to add/modify/delete comparing the index with the working tree.

## Ignore files

If you want to ignore files, create a file called `.bgitignore` and add one relative path per line. Wildcards are not allowed atm.

## How to run

- set env variables

```bash
export BGIT_TEST_REPO="srctest"
export BGIT_USER="email@gmail.com"
```

- run make build

```bash
  make build
```

- create a folder with the same name as set in `BGIT_TEST_REPO` and init the repo

```bash
  mkdir $BGIT_TEST_REPO
  ./bgit init
```

- run status to check and add a new file to commit

```bash
  ./bgit status
  /* add file to folder defined in `BGIT_TEST_REPO` */
  ./bgit commit "your message"
```

## What BGIT does not have

- Concept of stage. In the real GIT, when the user does `git add <filename>`, he is generating the blob object and adding an entry to the index, and this file starts to be tracked by GIT. You can even switch branches and this file will still be tracked. Here for simplicity I did not implemented the concept of stage. Any file added to the repo will be commited and if you switch branches without commting them, they will BE LOST.
- .pack files. GIT have special treatments for big files, I did not dig deeper into it, but the idea is to optimize since regular files are always copied no matter how small were the changes. There is no concept of saving the diffs.
- Parameters, a LOT of them. Each GIT command has MANY parameters to control every aspect of the command behavior. Here, for simplicity we almost do not have parameters in the commands.
- Config file. GIT uses global and local configs, which are basically files saved either on the local repo, or in the user home folder. BGIT for simplicity uses only the email which is set using the env variable `BGIT_USER`
- Merge and Rebase
- Remotes
- Stash
- Generate Patchs
- Add more things

## To Be Implemented

- reflog and `bgit log` command
- reset
- clone?

## References and related links

https://medium.com/data-management-for-researchers/git-under-the-hood-part-1-object-storage-in-git-57c9adfb5e5f
https://stackoverflow.com/questions/22968856/what-is-the-file-format-of-a-git-commit-object-data-structure
https://www.youtube.com/watch?v=RxHJdapz2p0 (video)

https://stackoverflow.com/questions/15765366/how-does-git-track-file-changes-internally
https://stackoverflow.com/questions/4084921/what-does-the-git-index-contain-exactly

https://www.freecodecamp.org/news/git-internals-objects-branches-create-repo/

https://benhoyt.com/writings/go-1brc/
https://github.com/git/git/blob/v2.21.1/commit.c
https://www.freecodecamp.org/news/boost-programming-skills-read-git-code/
https://git-scm.com/book/en/v2/Git-Internals-Git-Objects
https://www.freecodecamp.org/news/git-internals-objects-branches-create-repo/
