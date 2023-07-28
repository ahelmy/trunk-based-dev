package cmd

import (
	"fmt"
	"os"

	"github.com/ahelmy/trunk-based-dev/app"
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "trunk",
		Short: "Trunk based development tool",
		Long:  "Create release in seconds!",
	}
	if app.App.IsEmpty() {
		fmt.Println("No configuration found!")
		Init()
	}

	rootCmd.AddCommand(ConfigCmd())
	rootCmd.AddCommand(CreateReleaseCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
