package main

import (
	"fmt"
	"os"

	"github.com/tylerjgarland/git2git/repositories/github"
	"github.com/tylerjgarland/git2git/repositories/gitlab"
)

func main() {
	parsedArguments := parseTokens(os.Args)

	if parsedArguments.TargetToken == "" || parsedArguments.OriginToken == "" {
		fmt.Println("No git tokens provided.")
		parsedArguments.Help = true
	}

	if parsedArguments.Help {
		fmt.Println("This tool only copies from GitLab to GitHub. It does not delete anything.")
		fmt.Println("Create access tokens that have read access to repositories/projects and read access to the user.")
		fmt.Println("Usage:")
		fmt.Println("--origin-token <token>")
		fmt.Println("--target-token <token>")
		fmt.Println("--origin github,gitlab")
		fmt.Println("--target github,gitlab")
		os.Exit(0)
	}

	switch parsedArguments.Origin + "-" + parsedArguments.Target {

	case "gitlab-github":
		SyncRepositories(parsedArguments.OriginToken, parsedArguments.TargetToken, gitlab.GetRepositories, github.GetRepositories, github.CreateRepository)
	case "github-gitlab":
		SyncRepositories(parsedArguments.OriginToken, parsedArguments.TargetToken, github.GetRepositories, gitlab.GetRepositories, gitlab.CreateRepository)
	case "github-github":
		SyncRepositories(parsedArguments.OriginToken, parsedArguments.TargetToken, github.GetRepositories, github.GetRepositories, github.CreateRepository)
	case "gitlab-gitlab":
		SyncRepositories(parsedArguments.OriginToken, parsedArguments.TargetToken, gitlab.GetRepositories, gitlab.GetRepositories, gitlab.CreateRepository)
	default:
		fmt.Println("Combination " + parsedArguments.Origin + "-" + parsedArguments.Target + " not supported.")
	}

}

func parseTokens(arguments []string) (parsedArgs Arguments) {
	parsedArgs = Arguments{}
	for index, arg := range arguments {
		switch arg {
		case "--origin-token":
			parsedArgs.OriginToken = arguments[index+1]
		case "--target-token":
			parsedArgs.TargetToken = arguments[index+1]
		case "--target":
			parsedArgs.Target = arguments[index+1]
		case "--origin":
			parsedArgs.Origin = arguments[index+1]
		case "-h":
			parsedArgs.Help = true
		case "--help":
			parsedArgs.Help = true
		}

	}
	return parsedArgs
}

type Arguments struct {
	OriginToken string
	TargetToken string
	Help        bool
	Target      string
	Origin      string
}
