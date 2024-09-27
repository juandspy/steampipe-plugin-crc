package crc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const openshiftInsightsVulnerabilitiesV1ClusterExposedImages = "openshift_insights_vulnerabilities_v1_cluster_exposed_images"

type vulnerabilitiesV1ClusterExposedImagesResponse struct {
	Data []struct {
		Name     string `json:"name"`
		Registry string `json:"registry"`
		Version  string `json:"version"`
	} `json:"data"`
	Meta struct{} `json:"meta"`
}

func tableVulnerabilitiesClusterExposedImagesV1(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        openshiftInsightsVulnerabilitiesV1ClusterExposedImages,
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
	const functionName = "getVulnerabilitiesClusterExposedImagesV1"

	clusterID := d.EqualsQualString("cluster_id")

	if clusterID == "" {
		err := errors.New("you must specify a Cluster ID")
		pluginLogError(ctx, openshiftInsightsVulnerabilitiesV1ClusterExposedImages, functionName, "query_error", err)
		return nil, err
	}

	client, err := connect(ctx, d, defaultTimeout)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsVulnerabilitiesV1ClusterExposedImages, functionName, "client_error", err)
		return nil, err
	}

	url := fmt.Sprintf("https://console.redhat.com/api/ocp-vulnerability/v1/clusters/%s/exposed_images", clusterID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsVulnerabilitiesV1ClusterExposedImages, functionName, "request_error", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsVulnerabilitiesV1ClusterExposedImages, functionName, "api_error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with status code %d", resp.StatusCode)
		pluginLogError(ctx, openshiftInsightsVulnerabilitiesV1ClusterExposedImages, functionName, "api_error", err)
		return nil, err
	}

	exposedImagesResponse, err := decodeVulnerabilitiesClusterExposedImagesV1Response(resp.Body)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsVulnerabilitiesV1ClusterExposedImages, functionName, "decode_error", err)
		return nil, err
	}

	if len(exposedImagesResponse.Data) == 0 {
		return nil, nil
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
