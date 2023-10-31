package types

type HelmRepo struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

func (r HelmRepo) WithChartName(chartName string) string {
	return r.Name + "/" + chartName
}
