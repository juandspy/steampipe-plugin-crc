package vulnerabilities

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/juandspy/steampipe-plugin-crc/crc/utils"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const V1ClustersTableName = "openshift_insights_vulnerabilities_v1_clusters"

type VulnerabilitiesV1ClustersResponse struct {
	Data []struct {
		CvesSeverity struct {
			Critical  int `json:"critical"`
			Important int `json:"important"`
			Low       int `json:"low"`
			Moderate  int `json:"moderate"`
		} `json:"cves_severity"`
		DisplayName string `json:"display_name"`
		ID          string `json:"id"`
		LastSeen    string `json:"last_seen"`
		Provider    string `json:"provider"`
		Status      string `json:"status"`
		Type        string `json:"type"`
		Version     string `json:"version"`
	} `json:"data"`
	Meta struct{} `json:"meta"`
}

func TableClustersV1(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        V1ClustersTableName,
		Description: "Retrieves all clusters for given organization, retrieves the impacting rules for each cluster and the count of impacting CVEs.",
		List: &plugin.ListConfig{
			Hydrate: listVulnerabilitiesClustersV1,
		},
		Columns: []*plugin.Column{
			{
				Name:        "cluster_id",
				Type:        proto.ColumnType_STRING,
				Description: "Cluster ID.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "Cluster display name.",
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "version",
				Type:        proto.ColumnType_STRING,
				Description: "Cluster version.",
				Transform:   transform.FromField("Version"),
			},
			{
				Name:        "provider",
				Type:        proto.ColumnType_STRING,
				Description: "Provider of the cluster.",
				Transform:   transform.FromField("Provider"),
			},
			{
				Name:        "last_seen",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time the cluster was last checked at.",
				Transform:   transform.FromField("LastSeen"),
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Status of the cluster.",
				Transform:   transform.FromField("Status"),
			},
			{
				Name:        "low_cves",
				Type:        proto.ColumnType_INT,
				Description: "The total low CVEs.",
				Transform:   transform.FromField("CvesSeverity.Low"),
			},
			{
				Name:        "moderate_cves",
				Type:        proto.ColumnType_INT,
				Description: "The total moderate CVEs.",
				Transform:   transform.FromField("CvesSeverity.Moderate"),
			},
			{
				Name:        "important_cves",
				Type:        proto.ColumnType_INT,
				Description: "The total important CVEs.",
				Transform:   transform.FromField("CvesSeverity.Important"),
			},
			{
				Name:        "critical_cves",
				Type:        proto.ColumnType_INT,
				Description: "The total critical CVEs.",
				Transform:   transform.FromField("CvesSeverity.Critical"),
			},
		},
	}
}

func listVulnerabilitiesClustersV1(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	endpoint := "api/ocp-vulnerability/v1/clusters"
	resp, err := utils.MakeAPIRequest(ctx, d, "GET", endpoint, nil, utils.DefaultTimeout)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V1ClustersTableName, "api_error", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
		utils.LogErrorUsingSteampipeLogger(ctx, V1ClustersTableName, "api_error", err)
		return nil, err
	}

	clusterResponse, err := decodeVulnerabilitiesClustersV1(resp.Body)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V1ClustersTableName, "decode_error", err)
		return nil, err
	}

	for _, cluster := range clusterResponse.Data {
		d.StreamListItem(ctx, cluster)
	}

	return nil, nil
}

func decodeVulnerabilitiesClustersV1(body io.ReadCloser) (VulnerabilitiesV1ClustersResponse, error) {
	var clusterResponse VulnerabilitiesV1ClustersResponse
	err := json.NewDecoder(body).Decode(&clusterResponse)
	return clusterResponse, err
}
