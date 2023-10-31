package release

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ardikabs/kasque/internal/helm"
	"github.com/ardikabs/kasque/internal/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	apiVersion = "toolkit.ardikabs.com/v1alpha1"
	kind       = "Release"
)

type Release struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ReleaseSpec `json:"spec,omitempty"`
}

type ReleaseSpec struct {
	Chart                 string         `json:"chart,omitempty"`
	Repo                  types.HelmRepo `json:"repo,omitempty"`
	Version               string         `json:"version,omitempty"`
	Values                []string       `json:"values,omitempty"`
	IncludeCRDs           bool           `json:"includeCRDs,omitempty"`
	IgnoreDeprecatedChart bool           `json:"ignoreDeprecatedChart,omitempty"`
}

func (r *Release) Validate() error {
	if r.APIVersion != apiVersion {
		return fmt.Errorf("invalid APIVersion %s, it must be %s\n", r.APIVersion, apiVersion)
	}

	if r.Kind != kind {
		return fmt.Errorf("invalid Kind %s, it must be %s\n", r.Kind, kind)
	}

	if r.Namespace == "" {
		return fmt.Errorf("namespace cannot be empty\n")
	}

	return nil
}

func (r *Release) Render() ([]byte, error) {
	buf := new(bytes.Buffer)
	helmRenderer := helm.NewHelmRenderer(buf, r.Spec.Repo)

	if err := helmRenderer.Render(context.Background(), r.toRenderParameter()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Release) toRenderParameter() helm.RenderParameter {
	return helm.RenderParameter{
		ReleaseName:           r.Name,
		ReleaseNamespace:      r.Namespace,
		IncludeCRDs:           r.Spec.IncludeCRDs,
		ChartName:             r.Spec.Chart,
		Version:               r.Spec.Version,
		Values:                r.Spec.Values,
		IgnoreDeprecatedChart: r.Spec.IgnoreDeprecatedChart,
	}
}
