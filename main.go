package main

import (
	"fmt"
	"github.com/JoelOtter/terraform-workspace-i/internal/terraform"
	"github.com/JoelOtter/terraform-workspace-i/internal/ui"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func main() {
	var debug bool

	cmd := &cobra.Command{
		Use: "git-branch-i",
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaces, err := terraform.GetWorkspaces()
			if err != nil {
				return err
			}
			if err := ui.ShowUI(workspaces); err != nil {
				return fmt.Errorf("failed to show UI: %w", err)
			}
			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.PersistentFlags().BoolVar(
		&debug,
		"debug",
		false,
		"Show debug output",
	)

	if err := cmd.Execute(); err != nil {
		if debug {
			log.Fatalln(err)
		}
		os.Exit(1)
	}
}
