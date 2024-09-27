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

const V1CVEsExposedImagesTableName = "openshift_insights_vulnerabilities_v1_cves_exposed_images"

type vulnerabilitiesV1CVEsExposedImagesResponse struct {
	Data []struct {
		ClustersExposed int    `json:"clusters_exposed"`
		Name            string `json:"name"`
		Registry        string `json:"registry"`
		Version         string `json:"version"`
	} `json:"data"`
	Meta struct{} `json:"meta"`
}

func TableCVEsExposedImagesV1(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        V1CVEsExposedImagesTableName,
		Description: "Retrieves exposed images for a specific CVE.",
		List: &plugin.ListConfig{
			Hydrate:    getVulnerabilitiesCVEsExposedImagesV1,
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
				Name:        "clusters_exposed",
				Type:        proto.ColumnType_INT,
				Description: "Number of clusters exposed to this image.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the exposed image.",
			},
			{
				Name:        "registry",
				Type:        proto.ColumnType_STRING,
				Description: "Registry of the exposed image.",
			},
			{
				Name:        "version",
				Type:        proto.ColumnType_STRING,
				Description: "Version of the exposed image.",
			},
		},
	}
}

func getVulnerabilitiesCVEsExposedImagesV1(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cveName := d.EqualsQualString("cve_name")

	if cveName == "" {
		err := errors.New("you must specify a CVE name")
		utils.LogErrorUsingSteampipeLogger(ctx, V1CVEsExposedImagesTableName, "query_error", err)
		return nil, err
	}

	endpoint := fmt.Sprintf("api/ocp-vulnerability/v1/cves/%s/exposed_images", cveName)
	resp, err := utils.MakeAPIRequest(ctx, d, "GET", endpoint, nil, utils.DefaultTimeout)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V1CVEsExposedImagesTableName, "api_error", err)
		return nil, err
	}
	defer resp.Body.Close()

	var exposedImagesResponse vulnerabilitiesV1CVEsExposedImagesResponse
	err = json.NewDecoder(resp.Body).Decode(&exposedImagesResponse)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V1CVEsExposedImagesTableName, "decode_error", err)
		return nil, err
	}

	for _, image := range exposedImagesResponse.Data {
		d.StreamListItem(ctx, image)
	}

	return nil, nil
}
