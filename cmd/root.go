package cmd

import (
	"os"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/ardikabs/helmize/internal/krm"
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
)

func MakeRoot() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "helmize",
		Short: "helmize is a KRM Function to enable Helm on Kustomize with Glob support",
		Example: dedent.Dedent(`
			Helmize is intended to be used as KRM function only,
			thus a standalone usage is not supported.

			KRM usage:

			$ kustomize build --enable-alpha-plugins --enable-exec .
		`),
	}

	cmd.Args = cobra.MaximumNArgs(0)
	cmd.RunE = run
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	stdinStat, _ := os.Stdin.Stat()

	// Check the StdIn content.
	if (stdinStat.Mode() & os.ModeCharDevice) != 0 {
		cmd.Help()
		os.Exit(1)
	}

	if err := fn.AsMain(fn.ResourceListProcessorFunc(krm.Process)); err != nil {
		os.Exit(1)
	}

	return nil
}
