package release_test

import (
	"testing"

	"github.com/ardikabs/kasque/internal/release"
	"github.com/ardikabs/kasque/internal/types"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
				Repo: types.HelmRepo{
					Name: "ardikabs",
					URL:  "https://charts.ardikabs.com",
				},
				Version: "0.4.1",
			},
		}

		bytes, err := rel.Render()
		assert.NoError(t, err)
		assert.Contains(t, string(bytes), "0.4.1")
	})

	t.Run("chart with OCI and enable includeCRDs", func(t *testing.T) {
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
				Repo: types.HelmRepo{
					URL: "oci://docker.io/envoyproxy/gateway-helm",
				},
				Version:     "v0.5.0",
				IncludeCRDs: true,
			},
		}

		bytes, err := rel.Render()
		assert.NoError(t, err)
		assert.Contains(t, string(bytes), "v0.5.0")
		assert.Contains(t, string(bytes), "CustomResourceDefinition")
	})
}
