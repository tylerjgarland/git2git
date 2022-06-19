package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	gitTokens := parseTokens(os.Args[1:])

	//POC: Only Gitlab to Github export available.
	Gitlab2Github(gitTokens.GitlabToken, gitTokens.GithubToken)

	fmt.Println(gitTokens)
}

func parseTokens(arguments []string) (parsedTokens Arguments) {
	parsedTokens = Arguments{}
	for _, arg := range arguments {
		splitArg := strings.Split(arg, " ")
		switch splitArg[0] {
		case "--gitlab-token":
			parsedTokens.GitlabToken = splitArg[1]
			break
		case "--github-token":
			parsedTokens.GithubToken = splitArg[1]
			break
		}
	}
	return parsedTokens
}

type Arguments struct {
	GitlabToken string
	GithubToken string
}
