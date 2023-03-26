package git

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/google/go-github/v50/github"
)

var (
	repository  string
	description string
	private     bool
	autoInit    bool
)

func CreateNewRepo(client *github.Client, repoName string) {
	flag.StringVar(&repository, "name", repoName, "Repository name")
	flag.StringVar(&description, "description", "Git Repo for binaries", "User repo description")
	flag.BoolVar(&private, "private", false, "Public repo")
	flag.BoolVar(&autoInit, "auto-init", true, "Auto Initialization")

	ctx := context.Background()

	rr := &github.Repository{Name: &repository, Private: &private, Description: &description, AutoInit: &autoInit}
	repo, _, err := client.Repositories.Create(ctx, "", rr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully created new repo: %v\n", repo.GetName())
}
