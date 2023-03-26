package git

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func Clone(vmName string, username string, repo string) {
	_, err := git.PlainClone(vmName+"/repository", false, &git.CloneOptions{
		URL: "https://github.com/" + username + "/" + repo + ".git",
	})

	if err != nil {
		fmt.Println(err)
	}

	repository, err := git.PlainOpen(vmName + "/repository")
	if err != nil {
		fmt.Println(err)
	}

	worktree, err := repository.Worktree()
	if err != nil {
		panic(err)
	}

	err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		fmt.Println(err)
	}
}

func Push(vmName string, username string, password string) {
	repository, err := git.PlainOpen("./repository")
	if err != nil {
		panic(err)
	}

	worktree, err := repository.Worktree()
	if err != nil {
		panic(err)
	}

	status, _ := worktree.Status()
	if status.IsClean() {
		fmt.Println("working tree clean !, there is no new commit")
	}

	worktree.Add(".")
	worktree.Commit("commit for binaries", &git.CommitOptions{
		Author: &object.Signature{
			Name: "FaaS Clients",
			When: time.Now(),
		},
	},
	)

	auth := &http.BasicAuth{
		Username: username,
		Password: password,
	}

	err = repository.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})

	fmt.Println("Success push binaries into repository")
}
