package helm

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/ardikabs/kasque/internal/types"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
)

var settings = cli.New()

type DefaultHelmRenderer struct {
	writer io.Writer
	client *clientFactory
	repo   types.HelmRepo
}

func NewHelmRenderer(writer io.Writer, helmRepo types.HelmRepo) *DefaultHelmRenderer {
	r := &DefaultHelmRenderer{
		writer: writer,
		client: newClientFactory(),
		repo:   helmRepo,
	}

	if !registry.IsOCI(helmRepo.URL) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*30)
		defer cancel()
		if err := registerRepo(ctx, helmRepo.Name, helmRepo.URL); err != nil {
			fmt.Println(err)
		}
	}

	return r
}

type RenderParameter struct {
	IncludeCRDs           bool
	ReleaseName           string
	ReleaseNamespace      string
	ChartName             string
	Version               string
	Values                []string
	IgnoreDeprecatedChart bool
}

func (r DefaultHelmRenderer) Render(ctx context.Context, param RenderParameter) error {
	client := r.client.GetClient(param.IncludeCRDs)
	if param.Version != "" {
		client.ChartPathOptions.Version = param.Version
	}

	client.ReleaseName = param.ReleaseName
	client.Namespace = param.ReleaseNamespace

	var chartName string
	if registry.IsOCI(r.repo.URL) {
		chartName = r.repo.URL
		registryClient, err := registry.NewClient()
		if err != nil {
			return fmt.Errorf("helm render: %w", err)
		}
		client.SetRegistryClient(registryClient)
	} else {
		chartName = r.repo.WithChartName(param.ChartName)
	}

	cp, err := client.LocateChart(chartName, settings)
	if err != nil {
		return fmt.Errorf("helm render: %w", err)
	}

	p := getter.All(settings)
	valuesOpts := new(values.Options)

	for _, v := range param.Values {
		matches, err := filepath.Glob(v)
		if err != nil {
			continue
		}

		valuesOpts.ValueFiles = append(valuesOpts.ValueFiles, matches...)
	}
	vals, err := valuesOpts.MergeValues(p)
	if err != nil {
		return fmt.Errorf("helm render: %w", err)
	}

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return fmt.Errorf("helm render: %w", err)
	}

	if err := isChartInstallable(chartRequested); err != nil {
		return fmt.Errorf("helm render: %w", err)
	}

	if chartRequested.Metadata.Deprecated && !param.IgnoreDeprecatedChart {
		return fmt.Errorf("Chart %s is deprecated", chartName)
	}

	helmRelease, err := client.Run(chartRequested, vals)
	if err != nil {
		return fmt.Errorf("helm render: %w", err)
	}

	if _, err := r.writer.Write([]byte(strings.TrimSpace(helmRelease.Manifest))); err != nil {
		return err
	}

	return nil
}

func isChartInstallable(ch *chart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return fmt.Errorf("%s charts are not installable", ch.Metadata.Type)
}
