package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/thoas/go-funk"
)

func Git2Git(gitlabToken string, githubToken string) {
	gitlabRepos, _ := GetGitlabRepositories(gitlabToken)
	githubRepos, _ := GetGithubRepositories(githubToken)

	gitReps, _ := funk.Difference(gitlabRepos, githubRepos)
	var wg sync.WaitGroup
	guard := make(chan struct{}, 1)

	for _, repo := range gitReps.([]GitRepository)[0:1] {
		wg.Add(1)
		guard <- struct{}{}
		defer wg.Done()
		go DownloadRepository(repo, githubToken)
		<-guard

	}
	wg.Wait()
	fmt.Println(githubRepos, gitlabRepos)
}

func DownloadRepository(repo GitRepository, gitHubToken string) {

	os.Mkdir("./gitrepos/"+repo.Name, 0755)

	defer func() {
		os.Remove("./gitrepos/" + repo.Name)
	}()

	gitRepo, err := git.PlainClone("./gitrepos/"+repo.Name, false, &git.CloneOptions{
		URL: repo.HTTPUrlToRepo,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Downloaded: " + repo.Name)

	MakeGithubRepository(gitHubToken, &repo)

	remote, err := gitRepo.Remote("origin")

	remote.Config().URLs = []string{repo.HTTPUrlToRepo}

	err = remote.Push(&git.PushOptions{
		RemoteName: "origin",
	})

	fmt.Println(err)
}
