package internal

import (
	"fmt"

	"github.com/ahelmy/trunk-based-dev/app"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type CreateRelease struct {
	SemVersion    string
	Commits       []string
	gitRepo       *GitRepository
	latestRelease string
	nextRelease   string
}

func (cr *CreateRelease) New() {
	gitRepo, err := NewGitRepository()
	if err != nil {
		fmt.Println("Failed to open current directory as git repo!")
		panic(err)
	}
	cr.gitRepo = gitRepo

	latestNext, err := cr.gitRepo.GetLatestAndNextRelease(app.App.Config.ReleaseCanidateName, cr.SemVersion)
	if err != nil {
		fmt.Printf("Failed to read releases %v\n", err.Error())
		panic(err)
	}
	cr.nextRelease, cr.latestRelease = latestNext[1], latestNext[0]

	fmt.Printf("Preparing release [%v] from [%v]...\n", latestNext[1], latestNext[0])
}

func (cr *CreateRelease) ListGitRepoCommits(searchCommit string) ([]object.Commit, error) {
	return cr.gitRepo.ListGitRepoCommits(searchCommit)
}

func (cr *CreateRelease) IsValidCommitHash(hash string) bool {
	return cr.gitRepo.IsValidCommitHash(hash)
}

func (cr *CreateRelease) Start(commits []string) {

	fmt.Printf("Creating release %v with %v commit(s)...\n", cr.nextRelease, len(commits))
	err := cr.gitRepo.CheckoutBranch(cr.latestRelease, false)
	if err != nil {
		fmt.Printf("Failed to checkout %v!\n", cr.latestRelease)
	}
	fmt.Printf("Checkout %v - Done\n", cr.latestRelease)
	err = cr.gitRepo.CheckoutBranch(cr.nextRelease, true)
	if err != nil {
		fmt.Printf("Failed to checkout %v!\n", cr.nextRelease)
	}
	fmt.Printf("Checkout %v - Done\n", cr.nextRelease)
	cr.gitRepo.CherryPickCommits(commits)

	err = cr.gitRepo.PushBranch()
	if err != nil {
		fmt.Printf("Failed to push branch %v! %v\n", cr.nextRelease, err.Error())
	}
}
