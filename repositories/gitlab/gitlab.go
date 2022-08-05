package gitlab

import (
	"fmt"

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

type gitlabUser struct {
	Username string
}

type gitlabRepository struct {
	Name          string `json:"path"`
	HTTPUrlToRepo string `json:"http_url_to_repo"`
	Archived      bool
}
