package crc

import (
	"context"

	"github.com/juandspy/steampipe-plugin-crc/crc/aggregator"
	gcs "github.com/juandspy/steampipe-plugin-crc/crc/gathering_conditions_service"
	"github.com/juandspy/steampipe-plugin-crc/crc/utils"
	"github.com/juandspy/steampipe-plugin-crc/crc/vulnerabilities"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-crc",
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: utils.ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			gcs.V1GatheringRules:                            gcs.TableGatheringRulesV1(ctx),
			gcs.V2RemoteConfiguration:                       gcs.TableGatheringRulesV2(ctx),
			aggregator.V2Clusters:                           aggregator.TableClustersV2(ctx),
			aggregator.V2ClusterReports:                     aggregator.TableClusterReportsV2(ctx),
			vulnerabilities.V1ClustersTableName:             vulnerabilities.TableClustersV1(ctx),
			vulnerabilities.V1ClusterCVEsTableName:          vulnerabilities.TableClusterCVEsV1(ctx),
			vulnerabilities.V1ClusterExposedImagesTableName: vulnerabilities.TableClusterExposedImagesV1(ctx),
			vulnerabilities.V1CVEsTableName:                 vulnerabilities.TableCVEsV1(ctx),
			vulnerabilities.V1CVEsExposedClustersTableName:  vulnerabilities.TableCVEsExposedClustersV1(ctx),
			vulnerabilities.V1CVEsExposedImagesTableName:    vulnerabilities.TableCVEsExposedImagesV1(ctx),
		},
	}
	return p
}
