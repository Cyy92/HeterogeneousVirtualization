package git

import (
	"strings"

	"github.com/google/go-github/v50/github"
)

func NewClient(username string, password string) *github.Client {

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client := github.NewClient(tp.Client())
	return client
}
