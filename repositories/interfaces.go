package repositories

type RepositoryType string

const (
	REPO_TYPE_GITLAB RepositoryType = "gitlab"
	REPO_TYPE_GITHUB RepositoryType = "github"
)

type GitRepository struct {
	Name          string
	HTTPUrlToRepo string
	Archived      bool
}
