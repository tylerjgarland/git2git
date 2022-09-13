package gitlab

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/thoas/go-funk"

	"github.com/tylerjgarland/git2git/repositories"
)

func GetRepositories(token string) ([]repositories.GitRepository, bool) {

	client := resty.New()

	userResp, err := client.R().
		EnableTrace().
		SetResult(gitlabUser{}).
		SetJSONEscapeHTML(false).
		SetQueryParams(map[string]string{
			"private_token": token,
		}).
		Get("https://gitlab.com/api/v4/user")

	userName := userResp.Result().(*gitlabUser).Username

	resp, err := client.R().
		EnableTrace().
		SetResult([]gitlabRepository{}).
		SetJSONEscapeHTML(false).
		SetQueryParams(map[string]string{
			"private_token": token,
			"owned":         "true",
		}).
		Get("https://gitlab.com/api/v4/projects")

	if err != nil {
		panic(err)
	}

	result := resp.Result().(*[]gitlabRepository)

	//git clone `https://oauth2:TOKEN@github.com/username/repo.git`
	// git clone

	return funk.Map(*result, func(repository gitlabRepository) repositories.GitRepository {
		return repositories.GitRepository{
			Name: repository.Name,
			// HTTPUrlToRepo: fmt.Sprintf("https://oauth2:%s@github.com/%s/%s.git", token, userName, repository.Name),
			HTTPUrlToRepo: fmt.Sprintf("https://%s:%s@gitlab.com/%s/%s.git", userName, token, userName, repository.Name),
			Archived:      repository.Archived,
		}
	}).([]repositories.GitRepository), true
}

func CreateRepository(token string, repositoryPtr *repositories.GitRepository) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Default().Print(err)
			panic("error creating github repository")
		}
	}()

	client := resty.New()

	userResp, _ := client.R().
		EnableTrace().
		SetResult(gitlabUser{}).
		SetJSONEscapeHTML(false).
		SetQueryParams(map[string]string{
			"private_token": token,
		}).
		Get("https://gitlab.com/api/v4/user")

	userName := userResp.Result().(*gitlabUser).Username

	repositoryExists, _ := client.R().
		EnableTrace().
		SetResult([]gitlabRepository{}).
		SetJSONEscapeHTML(false).
		SetHeader("Authorization", fmt.Sprintf("token %s", token)).
		SetQueryParams(map[string]string{
			"private_token": token,
			"searchName":    repositoryPtr.Name,
		}).
		SetPathParam("user", userName).
		SetPathParam("repo", repositoryPtr.Name).
		Get("https://gitlab.com/api/v4/users/:user/projects?private_token=:private_token&search=:searchName")

	if repositoryExists.StatusCode() != 404 {
		res := repositoryExists.Result().([]gitlabRepository)

		if funk.Any(res, func(repository gitlabRepository) bool {
			return repository.Name == repositoryPtr.Name
		}) {
			return false
		}
	}

	_, err := client.R().
		EnableTrace().
		SetQueryParams(map[string]string{
			"private_token": token,
		}).
		SetResult(gitlabRepository{}).
		SetJSONEscapeHTML(false).
		SetHeader("Authorization", fmt.Sprintf("token %s", token)).
		SetBody(map[string]string{
			"name":       repositoryPtr.Name,
			"visibility": "true",
		}).
		Post("https://gitlab.com/api/v4/projects")

	if err != nil {
		return false
	}

	repositoryPtr.HTTPUrlToRepo = fmt.Sprintf("https://%s:%s@gitlab.com/%s/%s.git", userName, token, userName, repositoryPtr.Name)

	return true
}

type gitlabUser struct {
	Username string
}

type gitlabRepository struct {
	Name          string `json:"path"`
	HTTPUrlToRepo string `json:"http_url_to_repo"`
	Archived      bool
}
