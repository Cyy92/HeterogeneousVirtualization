package git

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"gopkg.in/yaml.v2"
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

func Push() {
	// Extract exist working dir in userinfo.yaml
	var data map[string]interface{}
	yamlFile, err := ioutil.ReadFile("./config/userinfo.yaml")
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML data: %v", err)
	}

	gitInfo, ok := data["git_info"].([]interface{})
	if !ok {
		log.Fatal("Failed to get git_info from YAML data")
	}

	var wd = gitInfo[0].(map[interface{}]interface{})["exist_workingdir"]
	var username = gitInfo[0].(map[interface{}]interface{})["git_user"].(string)
	var password = gitInfo[0].(map[interface{}]interface{})["token"].(string)

	if wd == nil {
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
	} else {
		repository, err := git.PlainOpen(wd.(string))
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
}
