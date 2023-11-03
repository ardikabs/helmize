package processor

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/ardikabs/helmize/pkg/helm"
	"helm.sh/helm/v3/pkg/releaseutil"
	"sigs.k8s.io/yaml"
)

func wrapErr(err error) error {
	return fmt.Errorf("%w\n", err)
}

func ProcessResourceList(rl *fn.ResourceList) (bool, error) {
	generatedObjects, err := ProcessKubeObjects(rl.Items)
	if err != nil {
		return false, err
	}

	rl.Items = generatedObjects
	return true, nil
}

func ProcessKubeObjects(objs fn.KubeObjects) (fn.KubeObjects, error) {
	var generatedObjects fn.KubeObjects

	for _, manifest := range objs {
		release := new(helm.Release)
		if err := yaml.Unmarshal([]byte(manifest.String()), &release); err != nil {
			return nil, wrapErr(err)
		}

		if err := release.Validate(); err != nil {
			return nil, wrapErr(err)
		}

		var buf bytes.Buffer
		if err := release.Render(&buf); err != nil {
			return nil, wrapErr(err)
		}

		for _, manifest := range releaseutil.SplitManifests(buf.String()) {
			object, err := fn.ParseKubeObject([]byte(manifest))
			if err != nil {
				if strings.Contains(err.Error(), "expected exactly one object") {
					continue
				}
				return nil, wrapErr(err)
			}

			generatedObjects = append(generatedObjects, object)
		}
	}

	return generatedObjects, nil
}
