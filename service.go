package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/thoas/go-funk"
)

func Git2Git(gitlabToken string, githubToken string) {
	gitlabRepos, _ := GetGitlabRepositories(gitlabToken)
	githubRepos, _ := GetGithubRepositories(githubToken)

	gitReps, _ := funk.Difference(gitlabRepos, githubRepos)

	for _, repo := range gitReps.([]GitRepository) {
		os.Mkdir("./gitrepos/"+repo.Name, 0755)
		_, err := git.PlainClone("./gitrepos/"+repo.Name, false, &git.CloneOptions{
			URL:      repo.HTTPUrlToRepo,
			Progress: os.Stdout,
		})
		if err != nil {
			fmt.Println(err)
		}

	}

	fmt.Println(githubRepos, gitlabRepos)
}
