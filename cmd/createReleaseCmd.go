package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	. "github.com/ahelmy/trunk-based-dev/internal"
	"github.com/ahelmy/trunk-based-dev/utils"
	"github.com/spf13/cobra"
)

func CreateReleaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-release",
		Short: "Create release from latest release",
		Long:  "Create release from latest release, run --help for more options.",
		RunE: func(cmd *cobra.Command, args []string) error {
			newRelease := &CreateRelease{SemVersion: cmd.Flag("sv").Value.String()}
			newRelease.New()
			newRelease.Start(getCommits(newRelease))
			return nil
		},
	}

	// Add any flags or arguments specific to this command
	cmd.PersistentFlags().String("sv", "patch", "Incremental semversion (patch, minor, major)")
	return cmd
}

func options() []string {
	return []string{"/f", "/q", "/s", "/v"}
}

func checkFinish(command string) bool {
	return command == "/f"
}
func checkView(command string, commits []string) bool {
	if command == "/v" {
		if len(commits) > 0 {
			fmt.Println("Selected commits:")
		}
		for _, c := range commits {
			fmt.Println(c)
		}
		return true
	}
	return false
}
func checkSearch(command string, cr *CreateRelease, commits []string) (bool, []string) {
	newCommits := commits
	if command == "/s" {
		var searchedCommit string
		prompt := &survey.Input{
			Message: "Please enter text to search for commit(s):",
		}
		survey.AskOne(prompt, &searchedCommit)

		searchedCommits, err := cr.ListGitRepoCommits(searchedCommit)
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
			newCommits = append(newCommits, selectedCommits...)
		}
		return true, newCommits
	}
	return false, newCommits
}

func getCommits(cr *CreateRelease) []string {
	var commit string
	var commits []string
	firstIteration := true
	for {
		defaultCmd := "/f"
		if firstIteration {
			defaultCmd = "/s"
			firstIteration = false
		}
		prompt := &survey.Input{
			Message: "Please enter commit(s) hash or /s to search /v to view selected commits or /f for finish",
			Default: defaultCmd,
		}
		survey.AskOne(prompt, &commit)

		if checkFinish(commit) {
			break
		}

		_ = checkView(commit, commits)

		_, commits = checkSearch(commit, cr, commits)

		if !utils.Contains(options(), commit) && cr.IsValidCommitHash(commit) {
			commits = append(commits, commit)
		}

	}
	return commits
}
