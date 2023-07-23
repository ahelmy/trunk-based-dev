package internal

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ahelmy/trunk-based-dev/app"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
func options() []string {
	return []string{"/f", "/q", "/s", "/v"}
}

func CreateRelease(semversion string) {
	gitRepo, err := NewGitRepository()
	if err != nil {
		fmt.Println("Failed to open current directory as git repo!")
		return
	}
	name, err := gitRepo.GetLatestAndNextRelease(app.App.Config.ReleaseCanidateName, semversion)
	if err != nil {
		fmt.Printf("Failed to read commit %v\n", err.Error())
	}
	fmt.Printf("Preparing release [%v] from [%v]...\n", name[1], name[0])
	var commit string
	var commits []string
	for {
		prompt := &survey.Input{
			Message: "Please enter commit(s) hash or /s to search /v to view selected commits or /f for finish",
		}
		survey.AskOne(prompt, &commit)
		if commit == "/f" {
			break
		}
		if commit == "/v" {
			if len(commits) > 0 {
				fmt.Println("Selected commits:")
			}
			for _, c := range commits {
				fmt.Println(c)
			}
		}
		if commit == "/s" {
			var searchedCommit string
			prompt := &survey.Input{
				Message: "Please enter text to search for commit(s):",
			}
			survey.AskOne(prompt, &searchedCommit)
			searchedCommits, err := gitRepo.ListGitRepoCommits(searchedCommit)
			if err != nil {
				fmt.Println("Failed to fetch commits")
			} else {
				var selectedCommits []string
				var searchedCommitsStrArr []string
				var searchedCommitsMsgsStrArr []string
				for _, sc := range searchedCommits {
					searchedCommitsStrArr = append(searchedCommitsStrArr, sc.Hash.String())
					searchedCommitsMsgsStrArr = append(searchedCommitsMsgsStrArr, sc.Message)
				}
				prompt := &survey.MultiSelect{
					Message: "Please select commits:",
					Options: searchedCommitsStrArr,
					Description: func(value string, index int) string {
						return searchedCommitsMsgsStrArr[index]
					},
				}
				survey.AskOne(prompt, &selectedCommits)
				commits = append(commits, selectedCommits...)
			}
		}
		if !contains(options(), commit) && gitRepo.IsValidCommitHash(commit) {
			commits = append(commits, commit)
		}

	}

	fmt.Printf("Creating release %v with %v commit(s)...\n", name[1], len(commits))
	err = gitRepo.CheckoutBranch(name[0], false)
	if err != nil {
		fmt.Printf("Failed to checkout %v!\n", name[0])
	}
	fmt.Printf("Checkout %v - Done\n", name[0])
	err = gitRepo.CheckoutBranch(name[1], true)
	if err != nil {
		fmt.Printf("Failed to checkout %v!\n", name[1])
	}
	fmt.Printf("Checkout %v - Done\n", name[1])
	gitRepo.CherryPickCommits(commits)

	err = gitRepo.PushBranch()
	if err != nil {
		fmt.Printf("Failed to push branch %v! %v\n", name[1], err.Error())
	}
}
