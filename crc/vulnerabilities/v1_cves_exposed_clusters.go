package vulnerabilities

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/juandspy/steampipe-plugin-crc/crc/utils"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const V1CVEsExposedClustersTableName = "crc_openshift_insights_vulnerabilities_v1_cves_exposed_clusters"

type vulnerabilitiesV1CVEsExposedClustersResponse struct {
	Data []struct {
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

func TableCVEsExposedClustersV1(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        V1CVEsExposedClustersTableName,
		Description: "Retrieves exposed clusters for a specific CVE.",
		List: &plugin.ListConfig{
			Hydrate:    getVulnerabilitiesCVEsExposedClustersV1,
			KeyColumns: plugin.SingleColumn("cve_name"),
		},
		Columns: []*plugin.Column{
			{
				Name:        "cve_name",
				Type:        proto.ColumnType_STRING,
				Description: "The CVE name.",
				Transform:   transform.FromQual("cve_name"),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "Display name of the exposed cluster.",
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the exposed cluster.",
			},
			{
				Name:        "last_seen",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Last seen timestamp of the exposed cluster.",
			},
			{
				Name:        "provider",
				Type:        proto.ColumnType_STRING,
				Description: "Provider of the exposed cluster.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Status of the exposed cluster.",
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "Type of the exposed cluster.",
			},
			{
				Name:        "version",
				Type:        proto.ColumnType_STRING,
				Description: "Version of the exposed cluster.",
			},
		},
	}
}

func getVulnerabilitiesCVEsExposedClustersV1(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cveName := d.EqualsQualString("cve_name")

	if cveName == "" {
		err := errors.New("you must specify a CVE name")
		utils.LogErrorUsingSteampipeLogger(ctx, V1CVEsExposedClustersTableName, "query_error", err)
		return nil, err
	}

	endpoint := fmt.Sprintf("api/ocp-vulnerability/v1/cves/%s/exposed_clusters", cveName)
	resp, err := utils.MakeAPIRequest(ctx, d, "GET", endpoint, nil, utils.DefaultTimeout)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V1CVEsExposedClustersTableName, "api_error", err)
		return nil, err
	}
	defer resp.Body.Close()

	var exposedClustersResponse vulnerabilitiesV1CVEsExposedClustersResponse
	err = json.NewDecoder(resp.Body).Decode(&exposedClustersResponse)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V1CVEsExposedClustersTableName, "decode_error", err)
		return nil, err
	}

	for _, cluster := range exposedClustersResponse.Data {
		d.StreamListItem(ctx, cluster)
	}

	return nil, nil
}
