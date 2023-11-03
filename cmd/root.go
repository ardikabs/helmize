package cmd

import (
	"os"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/ardikabs/helmize/pkg/processor"
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
)

var files []string

func MakeRoot() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "helmize",
		Short: "helmize is a KRM Function to enable Helm on Kustomize with Glob support",
		Example: dedent.Dedent(`
			Initially Helmize used as KRM function only, but it also support for direct use through flag.

			> KRM usage:

			Within "generators" fields through annotations [config.kubernetes.io/function]

			annotations:
			  config.kubernetes.io/function: |
			    exec:
			      path: helmize

			$ kustomize build --enable-alpha-plugins --enable-exec /path/to/kustomization_dir

			> Direct usage:

			$ helmize [-f HELM_RELEASE_FILE]
		`),
	}

	cmd.Args = cobra.MaximumNArgs(0)
	cmd.RunE = run

	cmd.Flags().StringArrayVarP(&files, "file", "f", []string{}, "specify HelmRelease in a YAML file (can specify multiple times)")
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	// Direct usage
	if len(files) > 0 {
		var generatedObjects fn.KubeObjects
		for _, file := range files {

			raw, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			kubeObjects, err := fn.ParseKubeObjects(raw)
			if err != nil {
				return err
			}

			generated, err := processor.ProcessKubeObjects(kubeObjects)
			if err != nil {
				return err
			}
			generatedObjects = append(generatedObjects, generated...)
		}

		os.Stdout.Write([]byte(generatedObjects.String()))
		return nil
	}

	// KRM usage
	stdinStat, _ := os.Stdin.Stat()
	// Check the StdIn content.
	if (stdinStat.Mode() & os.ModeCharDevice) != 0 {
		return cmd.Help()
	}

	if err := fn.AsMain(fn.ResourceListProcessorFunc(processor.ProcessResourceList)); err != nil {
		os.Exit(1)
	}

	return nil
}
