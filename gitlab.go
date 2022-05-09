package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func GetRepositories(token string) {

	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		SetResult([]GitlabRepository{}).
		SetJSONEscapeHTML(false).
		SetQueryParams(map[string]string{
			"private_token": token,
			"owned":         "true",
		}).
		Get("https://gitlab.com/api/v4/projects")

	if err != nil {
		panic(err)
	}

	// var objmap []map[string]json.Number
	// err = json.Unmarshal(resp.Body(), &objmap)

	result := resp.Result().(*[]GitlabRepository)

	fmt.Println(result)
	// fmt.Println(ok)

	// for _, repo := range  {
	// 	repositoriesResponse = append(repositoriesResponse, repo)
	// }
}

type GitlabRepository struct {
	ID            float64
	Name          string
	HTTPUrlToRepo string `json:"http_url_to_repo"`
	Archived      bool
	Visibility    string
}
