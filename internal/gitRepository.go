package internal

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type NoLatestRelease struct {
	message string
}

func (e *NoLatestRelease) Error() string {
	return e.message
}

type GitRepository struct {
	*git.Repository
}

func NewGitRepository() (*GitRepository, error) {
	r, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal("Failed to open current dir")
		return nil, err
	}
	return &GitRepository{r}, nil
}

func (r *GitRepository) IsValidCommitHash(hash string) bool {
	_, err := r.CommitObject(plumbing.NewHash(hash))
	return err == nil
}
func (r *GitRepository) IsBranchExists(name string) bool {
	_, err := r.Branch(name)
	return err == nil
}

func (r *GitRepository) GetLatestAndNextRelease(releaseNamePrefix string, semversion string) ([]string, error) {
	latestRelease, err := r.getLatestReleaseName(releaseNamePrefix)
	if err != nil && !errors.Is(err, &NoLatestRelease{}) {
		return nil, err
	}
	v, err := semver.NewVersion(latestRelease[1])
	if err != nil {
		return nil, err
	}
	versions := []string{latestRelease[0]}
	// Get the next semantic version.
	if semversion == "major" {
		versions = append(versions, releaseNamePrefix+v.IncMajor().String())
	} else if semversion == "minor" {
		versions = append(versions, releaseNamePrefix+v.IncMinor().String())
	} else {
		versions = append(versions, releaseNamePrefix+v.IncPatch().String())
	}
	return versions, nil
}

func (r *GitRepository) getLatestReleaseName(releaseNamePrefix string) ([]string, error) {
	branches, err := r.Branches()
	if err != nil {
		return nil, err
	}
	var releaseName = releaseNamePrefix + "0.0.0"
	var releaseNameVersion = "0.0.0"
	branches.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			if strings.Contains(ref.Name().Short(), releaseNamePrefix) {
				branchReleaseVersion := strings.Replace(ref.Name().Short(), releaseNamePrefix, "", -1)
				v1, _ := semver.NewVersion(branchReleaseVersion)
				v2, _ := semver.NewVersion(releaseNameVersion)
				if v1.GreaterThan(v2) {
					releaseName = ref.Name().Short()
					releaseNameVersion = branchReleaseVersion
				}
			}
		}
		return nil
	})
	if len(releaseName) == 0 {
		return nil, &NoLatestRelease{message: "Failed to get latest release!"}
	}
	return []string{releaseName, releaseNameVersion}, nil

}

func (r *GitRepository) CheckoutBranch(name string, isCreate bool) error {
	worktree, err := r.Worktree()
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Create: isCreate,
		Branch: plumbing.NewBranchReferenceName(name),
	})
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
func (r *GitRepository) PushBranch() error {
	cmd := exec.Command("git", "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
func (r *GitRepository) CherryPickCommits(commits []string) {
	for _, commit := range commits {
		cmd := exec.Command("git", "cherry-pick", commit)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
			fmt.Printf("Failed to pick commit: %v\n", commit)
		}
	}
}
func (r *GitRepository) ListGitRepoCommits(filter string) ([]object.Commit, error) {
	commits := []object.Commit{}
	// Get all references in the repository
	commitObjects, err := r.CommitObjects()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	commitObjects.ForEach(func(commit *object.Commit) error {
		filter = strings.Trim(filter, " ")
		if (len(filter) > 0 && (strings.Contains(commit.Message, filter) || strings.Contains(commit.Hash.String(), filter))) || len(filter) == 0 {
			commits = append(commits, *commit)
		}
		return nil
	})
	return commits, err
}
