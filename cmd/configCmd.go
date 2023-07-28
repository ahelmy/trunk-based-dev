package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ahelmy/trunk-based-dev/app"
	"github.com/spf13/cobra"
)

func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize trunk configuration",
		Long:  "Use this command to setup tunk configuration, creating .trunk file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			Init()
			return nil
		},
	}
	return cmd
}

func Init() {
	var rcName string
	prompt := &survey.Input{
		Message: "Please enter release canidate prefix:",
		Default: "release-v",
	}
	survey.AskOne(prompt, &rcName)
	var releaseName string
	prompt = &survey.Input{
		Message: "Please enter release name prefix:",
		Default: "v-",
	}
	survey.AskOne(prompt, &releaseName)
	err := app.App.Save(releaseName, rcName)
	if err != nil {
		fmt.Println("Failed to save .trunk file!")
		panic(err)
	}
}
