package helm

import (
	"context"
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func registerRepo(ctx context.Context, name, url string) error {
	errChan := make(chan error)

	go func() {
		b, err := os.ReadFile(settings.RepositoryConfig)
		if err != nil && !os.IsNotExist(err) {
			errChan <- err
			return
		}

		var repoFile repo.File
		if err := yaml.Unmarshal(b, &repoFile); err != nil {
			errChan <- err
			return
		}

		entry := repo.Entry{
			Name: name,
			URL:  url,
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
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
