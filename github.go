package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/thoas/go-funk"
)

func GetGithubRepositories(token string) ([]GitRepository, bool) {

	client := resty.New()

	userResponse, err := client.R().
		EnableTrace().
		SetResult(GithubUser{}).
		SetJSONEscapeHTML(false).
		SetHeader("Authorization", fmt.Sprintf("token %s", token)).
		SetQueryParams(map[string]string{
			"access_token": token,
		}).
		Get("https://api.github.com/user")

	userName := userResponse.Result().(*GithubUser).Login

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
		panic(err)
	}

	result := repositoriesResponse.Result().(*GithubRepositoryCollection)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic occurred:", err)
		}
	}()

	//https://stackoverflow.com/questions/42148841/github-clone-with-oauth-access-token

	return funk.Map(result.Items, func(repository GithubRepository) GitRepository {
		return GitRepository{
			Name:          repository.FullName,
			HTTPUrlToRepo: repository.CloneUrl,
			Archived:      repository.Archived,
		}
	}).([]GitRepository), true
}

func MakeGithubRepository(token string, repositoryPtr *GitRepository) {
	client := resty.New()

	userResponse, err := client.R().
		EnableTrace().
		SetResult(GithubUser{}).
		SetJSONEscapeHTML(false).
		SetHeader("Authorization", fmt.Sprintf("token %s", token)).
		SetQueryParams(map[string]string{
			"access_token": token,
		}).
		Get("https://api.github.com/user")

	userName := userResponse.Result().(*GithubUser).Login

	repositoriesResponse, err := client.R().
		EnableTrace().
		SetResult(GithubRepository{}).
		SetJSONEscapeHTML(false).
		SetHeader("Authorization", fmt.Sprintf("token %s", token)).
		SetBody(map[string]string{
			"name":    repositoryPtr.Name,
			"private": "true",
		}).
		Post("https://api.github.com/user/repos")

	fmt.Println(repositoriesResponse.String())

	if err != nil {
		panic(err)
	}

	// //git clone
	// result := repositoriesResponse.Result().(*GithubRepositoryCollection)

	repositoryPtr.HTTPUrlToRepo = fmt.Sprintf("https://%s@github.com/%s/%s.git", token, userName, repositoryPtr.Name)
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
