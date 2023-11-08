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
			Initially Helmize used as KRM function only, but it also support for direct use through flag
			as well as a Kustomize Legacy Plugin.

			> Legacy usage:

			export HELMIZE_PLUGIN_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/kustomize/plugin/toolkit.ardikabs.com/v1alpha1/helmrelease"
			mkdir -p $HELMIZE_PLUGIN_DIR
			mv /path/to/helmize-binary "${HELMIZE_PLUGIN_DIR}/HelmRelease"

			$ kustomize build --enable-alpha-plugins /path/to/kustomization_dir

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

	cmd.Args = cobra.MaximumNArgs(1)
	cmd.RunE = run

	cmd.Flags().StringArrayVarP(&files, "file", "f", []string{}, "specify HelmRelease in a YAML file (can specify multiple times)")
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	switch {
	case len(files) > 0:
		// Direct Use or Standalone
		if err := runAsStandalone(); err != nil {
			return err
		}
	case len(args) > 0:
		// Kustomize Legacy Plugin
		if err := runAsLegacyPlugin(args[0]); err != nil {
			return err
		}
	default:
		// Kustomize with KRM function
		stdinStat, _ := os.Stdin.Stat()
		// Check the StdIn content.
		if (stdinStat.Mode() & os.ModeCharDevice) != 0 {
			return cmd.Help()
		}

		if err := fn.AsMain(fn.ResourceListProcessorFunc(processor.ProcessResourceList)); err != nil {
			os.Exit(1)
		}
	}

	return nil
}

func runAsStandalone() error {
	var generatedObjects fn.KubeObjects
	for _, file := range files {

		rawObj, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		kubeObjects, err := fn.ParseKubeObjects(rawObj)
		if err != nil {
			return err
		}

		generated, err := processor.ProcessKubeObjects(kubeObjects)
		if err != nil {
			return err
		}
		generatedObjects = append(generatedObjects, generated...)
	}

	if _, err := os.Stdout.Write([]byte(generatedObjects.String())); err != nil {
		return err
	}

	return nil

}

func runAsLegacyPlugin(resourceFile string) error {
	rawObj, err := os.ReadFile(resourceFile)
	if err != nil {
		return err
	}

	kubeObjects, err := fn.ParseKubeObjects(rawObj)
	if err != nil {
		return err
	}

	var generatedObjects fn.KubeObjects
	generated, err := processor.ProcessKubeObjects(kubeObjects)
	if err != nil {
		return err
	}
	generatedObjects = append(generatedObjects, generated...)
	if _, err := os.Stdout.Write([]byte(generatedObjects.String())); err != nil {
		return err
	}
	return nil
}
