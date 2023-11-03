package helm

import "helm.sh/helm/v3/pkg/action"

type ClientOption func(client *action.Install)

func WithIncludeCRDs(include bool) ClientOption {
	return func(client *action.Install) {
		client.IncludeCRDs = include
	}
}

func WithChartVersion(version string) ClientOption {
	return func(client *action.Install) {
		client.ChartPathOptions.Version = version
	}
}

func WithInsecureSkipVerifyTLS(skip bool) ClientOption {
	return func(client *action.Install) {
		client.InsecureSkipTLSverify = skip
	}
}

func WithReleaseName(name string) ClientOption {
	return func(client *action.Install) {
		client.ReleaseName = name
	}
}

func WithReleaseNamespace(name string) ClientOption {
	return func(client *action.Install) {
		client.Namespace = name
	}
}
