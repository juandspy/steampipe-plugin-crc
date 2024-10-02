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

const V1ClusterExposedImagesTableName = "crc_openshift_insights_vulnerabilities_v1_cluster_exposed_images"

type vulnerabilitiesV1ClusterExposedImagesResponse struct {
	Data []struct {
		Name     string `json:"name"`
		Registry string `json:"registry"`
		Version  string `json:"version"`
	} `json:"data"`
	Meta struct{} `json:"meta"`
}

func TableClusterExposedImagesV1(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        V1ClusterExposedImagesTableName,
		Description: "Retrieves exposed images for a specific Cluster ID.",
		List: &plugin.ListConfig{
			Hydrate:    getVulnerabilitiesClusterExposedImagesV1,
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

func getVulnerabilitiesClusterExposedImagesV1(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	clusterID := d.EqualsQualString("cluster_id")

	if clusterID == "" {
		err := errors.New("you must specify a Cluster ID")
		utils.LogErrorUsingSteampipeLogger(ctx, V1ClusterExposedImagesTableName, "query_error", err)
		return nil, err
	}

	endpoint := fmt.Sprintf("api/ocp-vulnerability/v1/clusters/%s/exposed_images", clusterID)
	resp, err := utils.MakeAPIRequest(ctx, d, "GET", endpoint, nil, utils.DefaultTimeout)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V1ClusterExposedImagesTableName, "api_error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with status code %d", resp.StatusCode)
		utils.LogErrorUsingSteampipeLogger(ctx, V1ClusterExposedImagesTableName, "api_error", err)
		return nil, err
	}

	exposedImagesResponse, err := decodeVulnerabilitiesClusterExposedImagesV1Response(resp.Body)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V1ClusterExposedImagesTableName, "decode_error", err)
		return nil, err
	}

	for _, image := range exposedImagesResponse.Data {
		d.StreamListItem(ctx, image)
	}

	return nil, nil
}

func decodeVulnerabilitiesClusterExposedImagesV1Response(body io.ReadCloser) (vulnerabilitiesV1ClusterExposedImagesResponse, error) {
	var exposedImagesResponse vulnerabilitiesV1ClusterExposedImagesResponse
	err := json.NewDecoder(body).Decode(&exposedImagesResponse)
	return exposedImagesResponse, err
}
