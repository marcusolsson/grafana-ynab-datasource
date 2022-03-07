package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/marcusolsson/grafana-ynab-datasource/pkg/ynab"
)

const pluginID = "marcusolsson-ynab-datasource"

func main() {
	backend.SetupPluginEnvironment(pluginID)

	if err := datasource.Manage(pluginID, NewDataSource, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}

func NewDataSource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var (
		accessToken = settings.DecryptedSecureJSONData["accessToken"]
		client      = ynab.NewCacheClient(ynab.NewClient(accessToken))
	)

	return NewYNABDatasource(client), nil
}
