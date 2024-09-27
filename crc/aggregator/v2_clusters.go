package aggregator

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/juandspy/steampipe-plugin-crc/crc/utils"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const V2ClustersTableName = "openshift_insights_aggregator_v2_clusters"

type ClustersResponseV2 struct {
	Data []struct {
		ClusterID       string    `json:"cluster_id"`
		ClusterName     string    `json:"cluster_name"`
		Managed         bool      `json:"managed"`
		LastCheckedAt   time.Time `json:"last_checked_at,omitempty"`
		TotalHitCount   int       `json:"total_hit_count"`
		HitsByTotalRisk struct {
			Low      int `json:"1"`
			Moderate int `json:"2"`
			High     int `json:"3"`
			Critical int `json:"4"`
		} `json:"hits_by_total_risk"`
		ClusterVersion string `json:"cluster_version,omitempty"`
	} `json:"data"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
	Status string `json:"status"`
}

func TableClustersV2(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        V2ClustersTableName,
		Description: "Retrieves all clusters for given organization, retrieves the impacting rules for each cluster and calculates the count of impacting rules by total risk (severity == critical, high, moderate, low).",
		List: &plugin.ListConfig{
			Hydrate: listClustersV2,
		},
		Columns: []*plugin.Column{
			{
				Name:        "cluster_id",
				Type:        proto.ColumnType_STRING,
				Description: "Cluster ID.",
				Transform:   transform.FromField("ClusterID"),
			},
			{
				Name:        "cluster_name",
				Type:        proto.ColumnType_STRING,
				Description: "Cluster name.",
				Transform:   transform.FromField("ClusterName"),
			},
			{
				Name:        "cluster_version",
				Type:        proto.ColumnType_STRING,
				Description: "Cluster version.",
				Transform:   transform.FromField("ClusterVersion"),
			},
			{
				Name:        "managed",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether the cluster is managed.",
				Transform:   transform.FromField("Managed"),
			},
			{
				Name:        "last_checked_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time the cluster was last checked at.",
				Transform:   transform.FromField("LastCheckedAt"),
			},
			{
				Name:        "total_hit_count",
				Type:        proto.ColumnType_INT,
				Description: "The total hit count.",
				Transform:   transform.FromField("TotalHitCount"),
			},
			{
				Name:        "hits_by_total_risk",
				Type:        proto.ColumnType_JSON,
				Description: "The total hits by risk.",
				Transform:   transform.FromField("HitsByTotalRisk"),
			},
		},
	}
}

func listClustersV2(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	timeout := 60 * time.Second // this API endpoint is very slow

	endpoint := "api/insights-results-aggregator/v2/clusters"
	resp, err := utils.MakeAPIRequest(ctx, d, "GET", endpoint, nil, timeout)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V2ClustersTableName, "api_error", err)
		return nil, err
	}

	defer resp.Body.Close()

	clusterResponse, err := decodeClustersV2(resp.Body)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V2ClustersTableName, "decode_error", err)
		return nil, err
	}

	for _, cluster := range clusterResponse.Data {
		d.StreamListItem(ctx, cluster)
	}

	return nil, nil
}

func decodeClustersV2(body io.ReadCloser) (ClustersResponseV2, error) {
	var clusterResponse ClustersResponseV2
	err := json.NewDecoder(body).Decode(&clusterResponse)
	return clusterResponse, err
}
