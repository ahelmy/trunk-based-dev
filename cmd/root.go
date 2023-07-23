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
		Short: "Your app description",
		Long:  "Your app long description",
	}
	if app.App.IsEmpty() {
		fmt.Println("You need to config it first!")
		Init()
	}
	rootCmd.AddCommand(ConfigCmd())
	rootCmd.AddCommand(CreateReleaseCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
