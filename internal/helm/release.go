package helm

import (
	"context"
	"fmt"
	"io"

	"github.com/ardikabs/helmize/internal/errs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	APIVersion = "toolkit.ardikabs.com/v1alpha1"
	Kind       = "HelmRelease"
)

type Release struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ReleaseSpec `json:"spec,omitempty"`
}

type ReleaseSpec struct {
	Chart   string   `json:"chart,omitempty"`
	Repo    HelmRepo `json:"repo,omitempty"`
	Version string   `json:"version,omitempty"`

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
	if r.APIVersion != APIVersion {
		return fmt.Errorf("%w; APIVersion %s, it must be %s", errs.ErrInvalidObject, r.APIVersion, APIVersion)
	}

	if r.Kind != Kind {
		return fmt.Errorf("%w; Kind %s, it must be %s", errs.ErrInvalidObject, r.Kind, Kind)
	}

	if r.Name == "" {
		return fmt.Errorf("%w; Release name cannot be empty", errs.ErrInvalidObject)
	}

	if r.Namespace == "" {
		return fmt.Errorf("%w; Release namespace cannot be empty", errs.ErrInvalidObject)
	}

	if r.Spec.Repo.Name == nil &&
		r.Spec.Repo.URL == nil &&
		r.Spec.Repo.Path == nil {
		return fmt.Errorf("%w; Repo must have either name, url or path", errs.ErrInvalidObject)
	}

	return nil
}

func (r *Release) Render(w io.Writer) error {
	helmRenderer, err := NewHelmRenderer(w, r.Spec.Repo)
	if err != nil {
		return err
	}

	if err := helmRenderer.Render(context.Background(), r.toRenderParameter()); err != nil {
		return err
	}

	return nil
}

func (r *Release) toRenderParameter() RenderParameter {
	return RenderParameter{
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
