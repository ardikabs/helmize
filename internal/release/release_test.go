package release_test

import (
	"bytes"
	"testing"

	"github.com/ardikabs/helmize/internal/helm"
	"github.com/ardikabs/helmize/internal/release"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func TestRelease_Render(t *testing.T) {

	t.Run("simple chart", func(t *testing.T) {
		rel := &release.Release{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "toolkit.ardikabs.com/v1alpha1",
				Kind:       "Release",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "default",
			},
			Spec: release.ReleaseSpec{
				Chart: "common-app",
				Repo: helm.HelmRepo{
					Name: ptr.To[string]("ardikabs"),
					URL:  ptr.To[string]("https://charts.ardikabs.com"),
				},
				Version: "0.4.1",
			},
		}

		var buf bytes.Buffer
		err := rel.Render(&buf)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "0.4.1")
	})

	t.Run("chart with OCI with includeCRDs", func(t *testing.T) {
		rel := &release.Release{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "toolkit.ardikabs.com/v1alpha1",
				Kind:       "Release",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "default",
			},
			Spec: release.ReleaseSpec{
				Repo: helm.HelmRepo{
					URL: ptr.To[string]("oci://docker.io/envoyproxy/gateway-helm"),
				},
				Version:     "v0.5.0",
				IncludeCRDs: true,
			},
		}

		var buf bytes.Buffer
		err := rel.Render(&buf)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "v0.5.0")
		assert.Contains(t, buf.String(), "CustomResourceDefinition")
	})
}
