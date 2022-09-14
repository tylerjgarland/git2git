package github

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/thoas/go-funk"

	"github.com/tylerjgarland/git2git/repositories"
)

func GetRepositories(token string) ([]repositories.GitRepository, bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Default().Print(err)
			panic("error getting github repositories")
		}
	}()

	client := resty.New()

	userName := getGithubUsername(token)

	if userName == "" {
		log.Default().Print("invalid github token")
		panic("invalid github token")
	}

	repositoriesResponse, err := client.R().
		EnableTrace().
		SetResult(githubRepositoryCollection{}).
		SetJSONEscapeHTML(false).
		SetHeader("Authorization", fmt.Sprintf("token %s", token)).
		SetQueryParams(map[string]string{
			"q": "user:" + userName,
		}).
		Get("https://api.github.com/search/repositories")

	if err != nil {
		panic("error getting github repositories:" + err.Error())
	}

	result := repositoriesResponse.Result().(*githubRepositoryCollection)

	//No repositories found.
	if len(result.Items) == 0 {
		log.Default().Print("No repositories found.")
		return nil, false
	}
	//https: //oauth-key-goes-here@github.com/username/repo.git
	//https://stackoverflow.com/questions/42148841/github-clone-with-oauth-access-token
	//"https://github.com/tylerjgarland/vocal-voter-web.git"
	return funk.Map(result.Items, func(repository githubRepository) repositories.GitRepository {
		return repositories.GitRepository{
			Name:          repository.FullName,
			HTTPUrlToRepo: strings.Replace(repository.CloneUrl, "https://", "https://"+token+"@", 1),
			Archived:      repository.Archived,
		}
	}).([]repositories.GitRepository), true
}

func CreateRepository(token string, repositoryPtr *repositories.GitRepository) string {
	defer func() {
		if err := recover(); err != nil {
			log.Default().Print(err)
			panic("error creating github repository")
		}
	}()

	client := resty.New()

	userName := getGithubUsername(token)

	if userName == "" {
		panic("invalid github token")
	}

	repositoryExists, _ := client.R().
		EnableTrace().
		SetResult(githubRepository{}).
		SetJSONEscapeHTML(false).
		SetHeader("Authorization", fmt.Sprintf("token %s", token)).
		SetPathParam("user", userName).
		SetPathParam("repo", repositoryPtr.Name).
		Get("https://api.github.com/repos/{user}/{repo}")

	if repositoryExists.StatusCode() != 404 {
		return fmt.Sprintf("https://%s@github.com/%s/%s.git", token, userName, repositoryPtr.Name)
	}

	_, err := client.R().
		EnableTrace().
		SetResult(githubRepository{}).
		SetJSONEscapeHTML(false).
		SetHeader("Authorization", fmt.Sprintf("token %s", token)).
		SetBody(map[string]string{
			"name":    repositoryPtr.Name,
			"private": "true",
		}).
		Post("https://api.github.com/user/repos")

	if err != nil {
		return ""
	}

	return fmt.Sprintf("https://%s@github.com/%s/%s.git", token, userName, repositoryPtr.Name)
}

func getGithubUsername(token string) string {
	defer func() {
		if err := recover(); err != nil {
			log.Default().Print(err)
			panic("error getting github username")
		}
	}()
	client := resty.New()

	//Grab username for repository search.
	userResponse, err := client.R().
		EnableTrace().
		SetResult(githubUser{}).
		SetJSONEscapeHTML(false).
		SetHeader("Authorization", fmt.Sprintf("token %s", token)).
		SetQueryParams(map[string]string{
			"access_token": token,
		}).
		Get("https://api.github.com/user")

	if err != nil {
		panic("error getting github username: " + err.Error())
	}
	return userResponse.Result().(*githubUser).Login
}

type githubUser struct {
	Login string
}

type githubRepositoryCollection struct {
	Items []githubRepository
}

type githubRepository struct {
	FullName string `json:"name"`
	CloneUrl string `json:"clone_url"`
	Archived bool
	Private  bool
}
