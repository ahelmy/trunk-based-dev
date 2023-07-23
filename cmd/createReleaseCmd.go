package cmd

import (
	. "github.com/ahelmy/trunk-based-dev/internal"
	"github.com/spf13/cobra"
)

func CreateReleaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-release",
		Short: "Command 2 description",
		Long:  "Command 2 long description",
		RunE: func(cmd *cobra.Command, args []string) error {
			CreateRelease(cmd.Flag("sv").Value.String())
			return nil
		},
	}

	// Add any flags or arguments specific to this command
	cmd.PersistentFlags().String("sv", "patch", "Incremental semversion (patch, minor, major)")
	return cmd
}
