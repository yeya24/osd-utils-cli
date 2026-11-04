package cad

import (
	"github.com/spf13/cobra"
)

func NewCmdCad() *cobra.Command {
	cadCmd := &cobra.Command{
		Use:               "cad",
		Short:             "Provides commands to run CAD tasks",
		Args:              cobra.NoArgs,
		DisableAutoGenTag: true,
	}

	cadCmd.AddCommand(newCmdRun())
	return cadCmd
}
