package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version represent an awssh version
var (
	Version   string
	GitCommit string
)

func MakeVersion() *cobra.Command {
	var command = &cobra.Command{
		Use:          "version",
		Short:        "Print the version number of kasque",
		Long:         `All software has versions. This is kasque's`,
		Example:      `  kasque version`,
		SilenceUsage: false,
	}

	command.Run = func(cmd *cobra.Command, args []string) {
		printVersion()
	}

	return command
}

func printVersion() {
	if len(Version) == 0 {
		fmt.Println("Version: dev")
	} else {
		fmt.Println("Version:", Version)
	}
	fmt.Println("Git Commit:", GitCommit)
}
