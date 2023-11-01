package helm

import (
	"helm.sh/helm/v3/pkg/action"
)

type clientFactory struct{}

func newClientFactory() *clientFactory {
	return &clientFactory{}
}

func (h clientFactory) GetClient(opts ...ClientOption) *action.Install {
	actionConfig := new(action.Configuration)
	installClient := action.NewInstall(actionConfig)
	installClient.DryRun = true
	installClient.ClientOnly = true
	installClient.UseReleaseName = true

	for _, o := range opts {
		o(installClient)
	}

	return installClient
}
