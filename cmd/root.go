package cmd

import (
	"os"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/ardikabs/kasque/internal/release"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

func MakeRoot() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "kasque",
		Short: "kasque is ...",
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

	if err := fn.AsMain(fn.ResourceListProcessorFunc(generator)); err != nil {
		os.Exit(1)
	}

	return nil
}

func generator(rl *fn.ResourceList) (bool, error) {
	var generatedObjects fn.KubeObjects
	for _, manifest := range rl.Items {
		release := new(release.Release)
		if err := yaml.Unmarshal([]byte(manifest.String()), &release); err != nil {
			rl.LogResult(err)
			return false, err
		}

		if err := release.Validate(); err != nil {
			rl.LogResult(err)
			return false, err
		}

		generated, err := release.Render()
		if err != nil {
			rl.LogResult(err)
			return false, err
		}

		objects, err := fn.ParseKubeObjects(generated)
		if err != nil {
			rl.LogResult(err)
			return false, err
		}

		generatedObjects = append(generatedObjects, objects...)
	}

	rl.Items = generatedObjects
	return true, nil
}
