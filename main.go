package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xdrive/githubclean/cleanup"
	"github.com/xdrive/githubclean/github_api"
)

func main() {
	githubUser := getEnv("GITHUB_USER")
	githubRepository := getEnv("GITHUB_REPOSITORY")
	userToken := getEnv("USER_TOKEN")

	githubClient := github_api.NewGithubClient(githubUser, githubRepository, userToken)

	refService := github_api.NewReferenceService(githubClient)
	refs, err := refService.ListReferences()
	if err != nil {
		log.Fatal("ListReferences: ", err)
	}

	pullsService := github_api.NewPullRequestService(githubClient)
	pulls, err := pullsService.ListPullRequests()
	if err != nil {
		log.Fatal("ListPullRequests: ", err)
	}

	refs = cleanup.FindRefsWithoutPulls(refs, pulls)

	repoService := github_api.NewRepositoryService(githubClient)
	timeBefore, _ := time.Parse("2006-01-02", "2017-01-10")
	refsProcessor := cleanup.NewReferencesProcessor(repoService, refService)
	refs = refsProcessor.FilterReference(refs, &cleanup.FilterOptions{OlderThan: timeBefore})
	if len(refs) == 0 {
		fmt.Println("There are no outdated references.")
		return
	}
	fmt.Println("Here is the list of outdated references:")
	for _,ref := range refs {
		fmt.Println(ref.URL)
	}

	fmt.Println("Are you sure you what to delete all of them? Type \"yes\" to proceed with clean-up.")
	startCleanup := ""
	fmt.Scanln(&startCleanup)

	if startCleanup != "yes" {
		return
	}

	fmt.Println("Starting cleanup")
	//fmt.Println(refs[0].Ref)
	refsProcessor.DeleteReferences(refs)


	return





	for _, ref := range refs {
		if ref.Object.Sha == "7048aee8d8e6191235a7da092ea403d76f47a669" {
			fmt.Println(ref.Ref)
			break
		}
	}
	//fmt.Printf("%v", *refs[0])
	return

	repositoryService := github_api.NewRepositoryService(githubClient)
	branches, err := repositoryService.ListBranches()
	if err != nil {
		log.Fatal("ListBranches: ", err)
	}

	i := 0
	for _, branch := range branches {
		branchCommit, err := repositoryService.GetCommit(branch.Commit.Sha)
		if err != nil {
			log.Fatal("GetCommmit: ", err)
		}
		fmt.Println(branch.Name, branchCommit.Commit.Author.Date, "before?", branchCommit.Commit.Author.Date.Before(timeBefore))
		i++
		if i > 0 {
			break
		}
	}
	return

	commit, err := repositoryService.GetCommit("0b960628f17a3dd3683076dd1ca8fd9548a8570c")
	if err != nil {
		log.Fatal("GetCommmit: ", err)
	}

	fmt.Println(commit.Commit.Author.Date)
	return
}

func getEnv(varName string) string {
	varValue, ok := os.LookupEnv(varName)
	if !ok {
		log.Fatal(varName, "env variable should be set")
	}

	return varValue
}
