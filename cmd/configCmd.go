package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/ahelmy/trunk-based-dev/app"
	"github.com/spf13/cobra"
)

func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Command 2 description",
		Long:  "Command 2 long description",
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
		Message: "Please enter release canidate name:",
	}
	survey.AskOne(prompt, &rcName)
	var releaseName string
	prompt = &survey.Input{
		Message: "Please enter release name:",
	}
	survey.AskOne(prompt, &releaseName)
	app.App.Save(releaseName, rcName)
}
