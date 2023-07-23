package main

import (
	"fmt"
	"log"

	"github.com/ahelmy/trunk-based-dev/cmd"
	"github.com/ahelmy/trunk-based-dev/app"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func main() {
	app.InitApp()
	cmd.Execute()
}

func listReleaseBranches() {
	r, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal("Failed to open current dir")
	}
	// Get all references in the repository
	refs, err := r.References()
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over the references and print branch names
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			fmt.Println(ref.Name().Short())
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
