package vulnerabilities

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/juandspy/steampipe-plugin-crc/crc/utils"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const V1ClusterCVEsTableName = "openshift_insights_vulnerabilities_v1_cluster_cves"

type vulnerabilitiesV1ClusterCVEsResponse struct {
	Data []struct {
		CVSS2Score  float64 `json:"cvss2_score"`
		CVSS3Score  float64 `json:"cvss3_score"`
		Description string  `json:"description"`
		Exploits    bool    `json:"exploits"`
		PublishDate string  `json:"publish_date"`
		Severity    string  `json:"severity"`
		Synopsis    string  `json:"synopsis"`
	} `json:"data"`
	Meta struct{} `json:"meta"`
}

func TableClusterCVEsV1(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        V1ClusterCVEsTableName,
		Description: "Retrieves CVE details for a specific Cluster ID.",
		List: &plugin.ListConfig{
			Hydrate:    getVulnerabilitiesClusterCVEsV1,
			KeyColumns: plugin.SingleColumn("cluster_id"),
		},
		Columns: []*plugin.Column{
			{
				Name:        "cluster_id",
				Type:        proto.ColumnType_STRING,
				Description: "The Cluster ID.",
				Transform:   transform.FromQual("cluster_id"),
			},
			{
				Name:        "cvss2_score",
				Type:        proto.ColumnType_DOUBLE,
				Description: "CVSS2 score of the CVE.",
				Transform:   transform.FromField("cvss2_score"),
			},
			{
				Name:        "cvss3_score",
				Type:        proto.ColumnType_DOUBLE,
				Description: "CVSS3 score of the CVE.",
				Transform:   transform.FromField("cvss3_score"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "Description of the CVE.",
				Transform:   transform.FromField("description"),
			},
			{
				Name:        "exploits",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether the CVE has known exploits.",
				Transform:   transform.FromField("exploits"),
			},
			{
				Name:        "publish_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date the CVE was published.",
				Transform:   transform.FromField("publish_date"),
			},
			{
				Name:        "severity",
				Type:        proto.ColumnType_STRING,
				Description: "Severity level of the CVE.",
				Transform:   transform.FromField("severity"),
			},
			{
				Name:        "synopsis",
				Type:        proto.ColumnType_STRING,
				Description: "Brief summary of the CVE.",
				Transform:   transform.FromField("synopsis"),
			},
		},
	}
}

func getVulnerabilitiesClusterCVEsV1(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	clusterID := d.EqualsQualString("cluster_id")

	if clusterID == "" {
		err := errors.New("you must specify a Cluster ID")
		utils.LogErrorUsingSteampipeLogger(ctx, V1ClusterCVEsTableName, "query_error", err)
		return nil, err
	}

	endpoint := fmt.Sprintf("api/ocp-vulnerability/v1/clusters/%s/cves", clusterID)
	resp, err := utils.MakeAPIRequest(ctx, d, "GET", endpoint, nil, utils.DefaultTimeout)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V1ClusterCVEsTableName, "api_error", err)
		return nil, err
	}
	defer resp.Body.Close()

	cveResponse, err := decodeVulnerabilitiesClusterCVEsV1Response(resp.Body)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V1ClusterCVEsTableName, "decode_error", err)
		return nil, err
	}

	for _, cve := range cveResponse.Data {
		// TODO: Simplify this
		row := map[string]interface{}{}
		row["cluster_id"] = clusterID
		row["synopsis"] = cve.Synopsis
		row["cvss2_score"] = cve.CVSS2Score
		row["cvss3_score"] = cve.CVSS3Score
		row["description"] = cve.Description
		row["exploits"] = cve.Exploits
		row["publish_date"] = cve.PublishDate
		row["severity"] = cve.Severity
		d.StreamListItem(ctx, row)
	}

	return nil, nil
}

func decodeVulnerabilitiesClusterCVEsV1Response(body io.ReadCloser) (vulnerabilitiesV1ClusterCVEsResponse, error) {
	var cveResponse vulnerabilitiesV1ClusterCVEsResponse
	err := json.NewDecoder(body).Decode(&cveResponse)
	return cveResponse, err
}
