package krm

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/ardikabs/helmize/internal/helm"
	"helm.sh/helm/v3/pkg/releaseutil"
	"sigs.k8s.io/yaml"
)

func wrapErr(err error) error {
	return fmt.Errorf("%w\n", err)
}

func Process(rl *fn.ResourceList) (bool, error) {
	var generatedObjects fn.KubeObjects
	for _, manifest := range rl.Items {
		release := new(helm.Release)
		if err := yaml.Unmarshal([]byte(manifest.String()), &release); err != nil {
			return false, wrapErr(err)
		}

		if err := release.Validate(); err != nil {
			return false, wrapErr(err)
		}

		var buf bytes.Buffer
		if err := release.Render(&buf); err != nil {
			return false, wrapErr(err)
		}

		for _, manifest := range releaseutil.SplitManifests(buf.String()) {
			object, err := fn.ParseKubeObject([]byte(manifest))
			if err != nil {
				if strings.Contains(err.Error(), "expected exactly one object") {
					continue
				}
				return false, wrapErr(err)
			}

			generatedObjects = append(generatedObjects, object)
		}
	}

	rl.Items = generatedObjects
	return true, nil
}
