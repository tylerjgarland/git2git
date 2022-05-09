package main

type Repository struct {
	URL      string
	Name     string
	RepoType RepositoryType
}

type RepositoryType string

const (
	REPO_TYPE_GITLAB RepositoryType = "gitlab"
	REPO_TYPE_GITHUB RepositoryType = "github"
)
