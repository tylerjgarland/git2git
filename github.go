package main

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/thoas/go-funk"
)

var userName string

func GetGithubUsername(token string) string {
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
		SetResult(GithubUser{}).
		SetJSONEscapeHTML(false).
		SetHeader("Authorization", fmt.Sprintf("token %s", token)).
		SetQueryParams(map[string]string{
			"access_token": token,
		}).
		Get("https://api.github.com/user")

	if err != nil {
		panic("error getting github username: " + err.Error())
	}

	if userName == "" {
		userName = userResponse.Result().(*GithubUser).Login
	}

	return userName

}

func GetGithubRepositories(token string) ([]GitRepository, bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Default().Print(err)
			panic("error getting github repositories")
		}
	}()

	client := resty.New()

	userName := GetGithubUsername(token)

	if userName == "" {
		log.Default().Print("invalid github token")
		panic("invalid github token")
	}

	repositoriesResponse, err := client.R().
		EnableTrace().
		SetResult(GithubRepositoryCollection{}).
		SetJSONEscapeHTML(false).
		SetHeader("Authorization", fmt.Sprintf("token %s", token)).
		SetQueryParams(map[string]string{
			"q": "user:" + userName,
		}).
		Get("https://api.github.com/search/repositories")

	if err != nil {
		panic("error getting github repositories:" + err.Error())
	}

	result := repositoriesResponse.Result().(*GithubRepositoryCollection)

	//No repositories found.
	if len(result.Items) == 0 {
		log.Default().Print("No repositories found.")
		return nil, false
	}

	//https://stackoverflow.com/questions/42148841/github-clone-with-oauth-access-token
	return funk.Map(result.Items, func(repository GithubRepository) GitRepository {
		return GitRepository{
			Name:          repository.FullName,
			HTTPUrlToRepo: repository.CloneUrl,
			Archived:      repository.Archived,
		}
	}).([]GitRepository), true
}

func CreateGithubRepository(token string, repositoryPtr *GitRepository) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Default().Print(err)
			panic("error creating github repository")
		}
	}()

	client := resty.New()

	userName := GetGithubUsername(token)

	if userName == "" {
		panic("invalid github token")
	}

	_, err := client.R().
		EnableTrace().
		SetResult(GithubRepository{}).
		SetJSONEscapeHTML(false).
		SetHeader("Authorization", fmt.Sprintf("token %s", token)).
		SetBody(map[string]string{
			"name":    repositoryPtr.Name,
			"private": "true",
		}).
		Post("https://api.github.com/user/repos")

	if err != nil {
		return false
	}

	repositoryPtr.HTTPUrlToRepo = fmt.Sprintf("https://%s@github.com/%s/%s.git", token, userName, repositoryPtr.Name)

	return true
}

type GithubUser struct {
	Login string
}

type GithubRepositoryCollection struct {
	Items []GithubRepository
}

type GithubRepository struct {
	FullName string `json:"name"`
	CloneUrl string `json:"clone_url"`
	Archived bool
	Private  bool
}
