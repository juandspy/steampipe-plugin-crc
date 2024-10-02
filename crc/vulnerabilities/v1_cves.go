package vulnerabilities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/juandspy/steampipe-plugin-crc/crc/utils"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const V1CVEsTableName = "crc_openshift_insights_vulnerabilities_v1_cves"

type vulnerabilitiesV1CVEsResponse struct {
	Data []struct {
		ClustersExposed int     `json:"clusters_exposed"`
		CVSS2Score      float64 `json:"cvss2_score"`
		CVSS3Score      float64 `json:"cvss3_score"`
		Description     string  `json:"description"`
		Exploits        bool    `json:"exploits"`
		ImagesExposed   int     `json:"images_exposed"`
		PublishDate     string  `json:"publish_date"`
		Severity        string  `json:"severity"`
		Synopsis        string  `json:"synopsis"`
	} `json:"data"`
	Meta struct{} `json:"meta"`
}

func TableCVEsV1(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        V1CVEsTableName,
		Description: "Retrieves CVEs affecting the current workload.",
		List: &plugin.ListConfig{
			Hydrate: listVulnerabilitiesCVEsV1,
		},
		Columns: []*plugin.Column{
			{
				Name:        "synopsis",
				Type:        proto.ColumnType_STRING,
				Description: "Brief summary of the CVE.",
				Transform:   transform.FromField("Synopsis"),
			},
			{
				Name:        "clusters_exposed",
				Type:        proto.ColumnType_INT,
				Description: "Number of clusters exposed to this CVE.",
				Transform:   transform.FromField("ClustersExposed"),
			},
			{
				Name:        "cvss2_score",
				Type:        proto.ColumnType_DOUBLE,
				Description: "CVSS2 score of the CVE.",
				Transform:   transform.FromField("CVSS2Score"),
			},
			{
				Name:        "cvss3_score",
				Type:        proto.ColumnType_DOUBLE,
				Description: "CVSS3 score of the CVE.",
				Transform:   transform.FromField("CVSS3Score"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "Description of the CVE.",
				Transform:   transform.FromField("Description"),
			},
			{
				Name:        "exploits",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether the CVE has known exploits.",
				Transform:   transform.FromField("Exploits"),
			},
			{
				Name:        "images_exposed",
				Type:        proto.ColumnType_INT,
				Description: "Number of images exposed to this CVE.",
				Transform:   transform.FromField("ImagesExposed"),
			},
			{
				Name:        "publish_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date the CVE was published.",
				Transform:   transform.FromField("PublishDate"),
			},
			{
				Name:        "severity",
				Type:        proto.ColumnType_STRING,
				Description: "Severity level of the CVE.",
				Transform:   transform.FromField("Severity"),
			},
		},
	}
}

func listVulnerabilitiesCVEsV1(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// TODO: add pagination, sorting, filtering and so on.

	endpoint := "api/ocp-vulnerability/v1/cves"
	resp, err := utils.MakeAPIRequest(ctx, d, "GET", endpoint, nil, utils.DefaultTimeout)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V1CVEsTableName, "api_error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with status code %d", resp.StatusCode)
		utils.LogErrorUsingSteampipeLogger(ctx, V1CVEsTableName, "api_error", err)
		return nil, err
	}

	var cveResponse vulnerabilitiesV1CVEsResponse
	err = json.NewDecoder(resp.Body).Decode(&cveResponse)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V1CVEsTableName, "decode_error", err)
		return nil, err
	}

	for _, cve := range cveResponse.Data {
		d.StreamListItem(ctx, cve)
	}

	return nil, nil
}
