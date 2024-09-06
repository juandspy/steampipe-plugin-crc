package crc

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-crc",
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			openshiftInsightsGCSV1GatheringRules:      tableInsightsGatheringRulesV1(ctx),
			openshiftInsightsGCSV2RemoteConfiguration: tableInsightsGatheringRulesV2(ctx),
		},
	}
	return p
}
