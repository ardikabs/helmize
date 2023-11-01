package release

import (
	"context"
	"fmt"
	"io"

	"github.com/ardikabs/helmize/internal/helm"
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
	Chart   string        `json:"chart,omitempty"`
	Repo    helm.HelmRepo `json:"repo,omitempty"`
	Version string        `json:"version,omitempty"`

	// Values represent list of files specified for the Helm release, it is similar with -f/--values Helm template flag
	Values []string `json:"values,omitempty"`
	// ValuesInline represent an inline values, it will take precedence over Values files
	ValuesInline map[string]interface{} `json:"valuesInline,omitempty"`

	IncludeCRDs           bool `json:"includeCRDs,omitempty"`
	InsecureSkipTLSVerify bool `json:"insecureSkipTLSVerify,omitempty"`
	IgnoreDeprecatedChart bool `json:"ignoreDeprecatedChart,omitempty"`
	CreateNamespace       bool `json:"createNamespace,omitempty"`
}

func (r *Release) Validate() error {
	if r.APIVersion != apiVersion {
		return fmt.Errorf("invalid APIVersion %s, it must be %s", r.APIVersion, apiVersion)
	}

	if r.Kind != kind {
		return fmt.Errorf("invalid Kind %s, it must be %s", r.Kind, kind)
	}

	if r.Name == "" {
		return fmt.Errorf("release name cannot be empty")
	}

	if r.Namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}

	if r.Spec.Repo.Name == nil &&
		r.Spec.Repo.URL == nil &&
		r.Spec.Repo.Path == nil {
		return fmt.Errorf("invalid repo; repo must have name, url, or path")
	}

	return nil
}

func (r *Release) Render(w io.Writer) error {
	helmRenderer, err := helm.NewHelmRenderer(w, r.Spec.Repo)
	if err != nil {
		return err
	}

	if err := helmRenderer.Render(context.Background(), r.toRenderParameter()); err != nil {
		return err
	}

	return nil
}

func (r *Release) toRenderParameter() helm.RenderParameter {
	return helm.RenderParameter{
		ReleaseName:           r.Name,
		ReleaseNamespace:      r.Namespace,
		ChartName:             r.Spec.Chart,
		Version:               r.Spec.Version,
		Values:                r.Spec.Values,
		ValuesInline:          r.Spec.ValuesInline,
		IncludeCRDs:           r.Spec.IncludeCRDs,
		InsecureSkipVerifyTLS: r.Spec.InsecureSkipTLSVerify,
		IgnoreDeprecatedChart: r.Spec.IgnoreDeprecatedChart,
		CreateNamespace:       r.Spec.CreateNamespace,
	}
}
