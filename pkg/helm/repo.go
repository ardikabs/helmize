package helm

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ardikabs/helmize/pkg/errs"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/utils/ptr"
)

type HelmRepo struct {
	// Name represents the repo entry name to refer for
	Name *string `json:"name,omitempty"`

	// URL represents the remote target to fetch the helm chart.
	// It only supports https://. http://, and oci:// protocol.
	URL *string `json:"url,omitempty"`

	// Path represents as local path of the helm chart.
	Path *string `json:"path,omitempty"`
}

func (r HelmRepo) Init(ctx context.Context, client *action.Install) error {
	if r.Path != nil {
		return nil
	}

	repoURL := ptr.Deref[string](r.URL, "")
	if registry.IsOCI(repoURL) {
		registryClient, err := registry.NewClient()
		if err != nil {
			return err
		}
		client.SetRegistryClient(registryClient)
		return nil
	}

	repoName := ptr.Deref[string](r.Name, "")
	errChan := make(chan error)
	go func(name, url string) {
		repoFile, err := r.loadFile(settings.RepositoryConfig)
		if err != nil {
			errChan <- err
			return
		}

		entry := repo.Entry{
			Name:                  name,
			URL:                   url,
			InsecureSkipTLSverify: client.InsecureSkipTLSverify,
		}

		chartRepo, err := repo.NewChartRepository(&entry, getter.All(settings))
		if err != nil {
			errChan <- err
			return
		}

		if _, err := chartRepo.DownloadIndexFile(); err != nil {
			errChan <- fmt.Errorf("chart repository couldn't be reached or %q is an invalid chart repository, %w", url, err)
			return
		}

		repoFile.Update(&entry)
		if err := repoFile.WriteFile(settings.RepositoryConfig, 0644); err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}(repoName, repoURL)

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r HelmRepo) loadFile(path string) (file *repo.File, err error) {
	file, err = repo.LoadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return
	}

	if err = file.WriteFile(path, 0644); err != nil {
		return
	}

	return
}

// NameAndChart returns the name and chart that should be used.
// The precedence during the NameAndChart generation are as follows:
// 1. Path will be used when available
// 2. Repo with OCI protocol will refer to the its URL
// 3. Repo with HTTP protocol will refer to the combined of repo Name and Chart Name, e.g., repoName/chartName.
func (r HelmRepo) NameAndChart(chartName string) (string, error) {
	if r.Path != nil {
		return ptr.Deref[string](r.Path, "."), nil
	}

	if r.Name == nil && r.URL == nil {
		return ".", nil
	}

	repoURL := ptr.Deref[string](r.URL, ".")
	if registry.IsOCI(repoURL) {
		return repoURL, nil
	}

	if r.Name == nil {
		return "", fmt.Errorf("%w; Repo name cannot be empty when repo refers to HTTP protocol", errs.ErrInvalidObject)
	}

	return ptr.Deref[string](r.Name, "") + "/" + chartName, nil
}
