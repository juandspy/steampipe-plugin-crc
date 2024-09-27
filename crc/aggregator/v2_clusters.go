package aggregator

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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
				Transform:   transform.FromField("cluster_id"),
			},
			{
				Name:        "cluster_name",
				Type:        proto.ColumnType_STRING,
				Description: "Cluster name.",
				Transform:   transform.FromField("cluster_name"),
			},
			{
				Name:        "cluster_version",
				Type:        proto.ColumnType_STRING,
				Description: "Cluster version.",
				Transform:   transform.FromField("cluster_version"),
			},
			{
				Name:        "managed",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether the cluster is managed.",
				Transform:   transform.FromField("managed"),
			},
			{
				Name:        "last_checked_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time the cluster was last checked at.",
				Transform:   transform.FromField("last_checked_at"),
			},
			{
				Name:        "total_hit_count",
				Type:        proto.ColumnType_INT,
				Description: "The total hit count.",
				Transform:   transform.FromField("total_hit_count"),
			},
			{
				Name:        "hits_by_total_risk",
				Type:        proto.ColumnType_JSON,
				Description: "The total hits by risk.",
				Transform:   transform.FromField("hits_by_total_risk"),
			},
		},
	}
}

func listClustersV2(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	const functionName = "listClustersV2"
	timeout := 60 * time.Second // this API endpoint is very slow
	client, err := utils.GetConsoleDotClient(ctx, d, timeout)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V2ClustersTableName, functionName, "client_error", err)
		return nil, err
	}

	url := "https://console.redhat.com/api/insights-results-aggregator/v2/clusters"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V2ClustersTableName, functionName, "request_error", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V2ClustersTableName, functionName, "api_error", err)
		return nil, err
	}

	defer resp.Body.Close()

	clusterResponse, err := decodeClustersV2(resp.Body)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V2ClustersTableName, functionName, "decode_error", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = errors.New(clusterResponse.Status)
		utils.LogErrorUsingSteampipeLogger(ctx, V2ClustersTableName, functionName, "api_error", err)
		return nil, err
	}

	for _, cluster := range clusterResponse.Data {
		row := map[string]interface{}{}
		row["cluster_id"] = cluster.ClusterID
		row["cluster_name"] = cluster.ClusterName
		row["cluster_version"] = cluster.ClusterVersion
		row["managed"] = cluster.Managed
		row["last_checked_at"] = cluster.LastCheckedAt
		row["total_hit_count"] = cluster.TotalHitCount
		row["hits_by_total_risk"] = cluster.HitsByTotalRisk
		d.StreamListItem(ctx, row)
	}

	return nil, nil
}

func decodeClustersV2(body io.ReadCloser) (ClustersResponseV2, error) {
	var clusterResponse ClustersResponseV2
	err := json.NewDecoder(body).Decode(&clusterResponse)
	return clusterResponse, err
}
