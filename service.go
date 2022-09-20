package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"

	"github.com/tylerjgarland/git2git/repositories"
)

func SyncRepositories(originToken string, destinationToken string, origin func(token string) ([]repositories.GitRepository, bool), target func(token string) ([]repositories.GitRepository, bool), createRepositoryAsync func(token string, repoPtr *repositories.GitRepository) string) {
	defer func() {
		if err := recover(); err != nil {
			log.Default().Print(err)
			log.Default().Print("Failed to sync repositories. Exiting.")
		}
	}()

	os.RemoveAll("./gitrepos")
	os.Mkdir("./gitrepos", 0755)

	var wg sync.WaitGroup
	var originReposChan, targetReposChan chan []repositories.GitRepository = make(chan []repositories.GitRepository, 1), make(chan []repositories.GitRepository, 1)

	go GetRepositories(originToken, wg, originReposChan, origin)
	go GetRepositories(destinationToken, wg, targetReposChan, target)

	wg.Wait()

	var copyFromRepos []repositories.GitRepository
	// var copyToRepos []repositories.GitRepository

	copyFromRepos = <-originReposChan
	// copyToRepos = <-targetReposChan

	if len(copyFromRepos) == 0 {
		log.Default().Print("No repositories to copy.")
		return
	}

	var reposToDownload []repositories.GitRepository

	// workingGitRepos, _ := funk.Difference(
	// 	funk.Map(copyFromRepos, func(item repositories.GitRepository) string { return item.Name }),
	// 	funk.Map(copyToRepos, func(item repositories.GitRepository) string { return item.Name }),
	// )

	// var stringReposToDownload []string

	// reposToDownload = funk.Filter(copyFromRepos, func(item repositories.GitRepository) bool {
	// 	return item.Name == "machinations"
	// }).([]repositories.GitRepository)
	// reposToDownload = funk.Filter(copyFromRepos, func(item repositories.GitRepository) bool { return funk.Contains(workingGitRepos, item.Name) }).([]repositories.GitRepository)
	reposToDownload = copyFromRepos
	//Limit to 3 concurrent git clones.
	concurrencyLimit := make(chan struct{}, 1)

	wg.Add(len(reposToDownload))

	for _, repo := range reposToDownload {
		func() {
			concurrencyLimit <- struct{}{}

			defer func() { <-concurrencyLimit }()

			go syncRepository(repo, destinationToken, &wg, createRepositoryAsync)

		}()
	}

	wg.Wait()
	fmt.Println("Sync complete")
}

func GetRepositories(token string, waitGroup sync.WaitGroup, reposChannel chan []repositories.GitRepository, getRepositories func(token string) ([]repositories.GitRepository, bool)) {
	defer waitGroup.Done()

	waitGroup.Add(1)
	repos, _ := getRepositories(token)

	reposChannel <- repos
}

func syncRepository(repo repositories.GitRepository, gitHubToken string, wgPtr *sync.WaitGroup, createRepositoryAsync func(token string, repoPtr *repositories.GitRepository) string) bool {
	repositoryDownloadDir := fmt.Sprintf("./gitrepos/%s", repo.Name)

	os.Mkdir(repositoryDownloadDir, 0755)

	defer func() {
		wgPtr.Done()
		os.Remove(repositoryDownloadDir)
	}()

	gitRepo, err := git.PlainClone(repositoryDownloadDir, false, &git.CloneOptions{
		URL: repo.HTTPUrlToRepo,
	})

	opts := &git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	}

	if err != nil {
		log.Default().Print(err)
		fmt.Printf("Failed to clone repository: %s : %s", err.Error(), repo.Name)
		fmt.Println()
		return false
	}

	remote, _ := gitRepo.Remote("origin")

	if err := remote.Fetch(opts); err != nil {
		fmt.Printf("Repo failed to create: %s", repo.Name)
		fmt.Println()
		return false
	}

	fmt.Printf("Downloaded: %s", repo.Name)
	fmt.Println()

	pushURL := createRepositoryAsync(gitHubToken, &repo)

	if pushURL == "" {
		fmt.Printf("Repo failed to create: %s", repo.Name)
		fmt.Println()
		return false
	}

	remote, err = gitRepo.CreateRemote(&config.RemoteConfig{
		Name: "Destination",
		URLs: []string{pushURL},
	})

	if err != nil {
		fmt.Printf("Failed to create new remote: %s", repo.Name)
		fmt.Println()
		fmt.Println(err)
		return false
	}

	refs := make([]config.RefSpec, 0)
	refs = append(refs, config.RefSpec(fmt.Sprintf("%s:%s", "*", "*")))

	err = remote.Push(&git.PushOptions{
		RemoteName: "Destination",
		RefSpecs:   refs,
	})

	if err != nil {
		errorString := err.Error()
		if errorString == "authorization failed" {
			fmt.Printf("Not allowed to access repository. Check permissions.: %s", repo.Name)
			fmt.Println()
			return false
		} else if !strings.Contains(errorString, "deny updating a hidden ref") {
			fmt.Printf("Failed to push repository: %s", repo.Name)
			fmt.Println()
			fmt.Println(err)
			return false
		}
	}

	fmt.Printf("Synced: %s", repo.Name)
	fmt.Println()

	return true
}
