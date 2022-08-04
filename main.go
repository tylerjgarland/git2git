package main

import (
	"fmt"
	"os"
)

func main() {
	gitTokens := parseTokens(os.Args)

	if gitTokens.GithubToken == "" || gitTokens.GitlabToken == "" {
		fmt.Println("No git tokens provided.")
		gitTokens.Help = true
	}

	if gitTokens.Help {
		fmt.Println("This tool only copies from GitLab to GitHub. It does not delete anything.")
		fmt.Println("Create access tokens that have read access to repositories/projects and read access to the user.")
		fmt.Println("Usage:")
		fmt.Println("--gitlab-token <token>")
		fmt.Println("--github-token <token>")
		os.Exit(0)
	}

	//POC: Only Gitlab to Github export available.
	Gitlab2Github(gitTokens.GitlabToken, gitTokens.GithubToken)
}

func parseTokens(arguments []string) (parsedTokens Arguments) {
	parsedTokens = Arguments{}
	for index, arg := range arguments {
		switch arg {
		case "--gitlab-token":
			parsedTokens.GitlabToken = arguments[index+1]
			break
		case "--github-token":
			parsedTokens.GithubToken = arguments[index+1]
			break
		case "-h":
		case "--help":
			parsedTokens.Help = true
			break
		}

	}
	return parsedTokens
}

type Arguments struct {
	GitlabToken string
	GithubToken string
	Help        bool
}
