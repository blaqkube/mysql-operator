package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
)

func NewRootCmd(in string) *cobra.Command {
	return &cobra.Command{
		Use:   "hugo",
		Short: "Hugo is a very fast static site generator",
		Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), in)
			return nil
		},
	}
}

func Test_ExecuteCommand(t *testing.T) {
	cmd := NewRootCmd("hi")
	cmd.Execute()
}
