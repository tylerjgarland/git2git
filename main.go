package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := parseArguments(os.Args[1:])

	GetRepositories(args.GitlabToken)

	fmt.Println(args)
}

func parseArguments(arguments []string) (parsedArgs Arguments) {
	parsedArgs = Arguments{}
	for _, arg := range arguments {
		splitArg := strings.Split(arg, " ")
		switch splitArg[0] {
		case "--gitlab-token":
			parsedArgs.GitlabToken = splitArg[1]
			break
		case "--github-token":
			parsedArgs.GithubToken = splitArg[1]
			break
		}
	}
	return parsedArgs
}

type Arguments struct {
	GitlabToken string
	GithubToken string
}
