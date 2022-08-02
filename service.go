package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/thoas/go-funk"
)

func Gitlab2Github(gitlabToken string, githubToken string) {
	os.RemoveAll("./gitrepos")
	os.Mkdir("./gitrepos", 0755)

	var wg sync.WaitGroup
	var gitlabReposChan, githubReposChan chan []GitRepository = make(chan []GitRepository, 1), make(chan []GitRepository, 1)

	go GetRepositories(gitlabToken, wg, gitlabReposChan, GetGitlabRepositories)
	go GetRepositories(githubToken, wg, githubReposChan, GetGithubRepositories)

	wg.Wait()

	gitlabRepos := <-gitlabReposChan
	githubRepos := <-githubReposChan

	gitReps, _ := funk.Difference(
		funk.Map(gitlabRepos, func(item GitRepository) string { return item.Name }),
		funk.Map(githubRepos, func(item GitRepository) string { return item.Name }),
	)

	downloadRepos := funk.Filter(gitlabRepos, func(item GitRepository) bool { return funk.Contains(gitReps, item.Name) }).([]GitRepository)

	//Limit to 3 concurrent git clones.
	concurrencyLimit := make(chan struct{}, 3)

	wg.Add(len(downloadRepos))

	for _, repo := range downloadRepos {
		func() {
			concurrencyLimit <- struct{}{}

			defer func() { <-concurrencyLimit }()

			go CloneRepository(repo, githubToken, &wg)
		}()
	}

	wg.Wait()
	fmt.Println("Sync complete")
}

func GetRepositories(token string, waitGroup sync.WaitGroup, reposChannel chan []GitRepository, function func(token string) ([]GitRepository, bool)) {
	defer waitGroup.Done()

	waitGroup.Add(1)
	repos, ok := function(token)

	if !ok {
		panic("error getting repositories")
	}
	reposChannel <- repos
}

func CloneRepository(repo GitRepository, gitHubToken string, wgPtr *sync.WaitGroup) bool {
	repositoryDownloadDir := fmt.Sprintf("./gitrepos/%s", repo.Name)

	os.Mkdir(repositoryDownloadDir, 0755)

	defer func() {
		wgPtr.Done()
		os.Remove(repositoryDownloadDir)
	}()

	gitRepo, err := git.PlainClone(repositoryDownloadDir, false, &git.CloneOptions{
		URL: repo.HTTPUrlToRepo,
	})

	if err != nil {
		log.Default().Print(err)
		fmt.Printf("Failed to clone repository: %s : %s", err.Error(), repo.Name)
		fmt.Println()
		return false
	}

	fmt.Printf("Downloaded: %s", repo.Name)
	fmt.Println()

	ok := CreateGithubRepository(gitHubToken, &repo)

	if !ok {
		fmt.Printf("Failed to create repository in GitHub: %s", repo.Name)
		return false
	}

	fmt.Println("Created repository in GitHub")

	remote, _ := gitRepo.Remote("origin")

	remote.Config().URLs = []string{repo.HTTPUrlToRepo}

	err = remote.Push(&git.PushOptions{
		RemoteName: "origin",
	})

	if err != nil {
		fmt.Printf("Failed to push repository to GitHub: %s", repo.Name)
		return false
	}

	fmt.Println("Pushed repository to GitHub")

	wgPtr.Done()

	fmt.Printf("Synced: %s", repo.Name)
	fmt.Println()

	return true
}
