package crc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const openshiftInsightsVulnerabilitiesV1CVEsExposedImages = "openshift_insights_vulnerabilities_v1_cves_exposed_images"

type vulnerabilitiesV1CVEsExposedImagesResponse struct {
	Data []struct {
		ClustersExposed int    `json:"clusters_exposed"`
		Name            string `json:"name"`
		Registry        string `json:"registry"`
		Version         string `json:"version"`
	} `json:"data"`
	Meta struct{} `json:"meta"`
}

func tableVulnerabilitiesCVEsExposedImagesV1(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        openshiftInsightsVulnerabilitiesV1CVEsExposedImages,
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
	const functionName = "getVulnerabilitiesCVEsExposedImagesV1"

	cveName := d.EqualsQualString("cve_name")

	if cveName == "" {
		err := errors.New("you must specify a CVE name")
		pluginLogError(ctx, openshiftInsightsVulnerabilitiesV1CVEsExposedImages, functionName, "query_error", err)
		return nil, err
	}

	client, err := connect(ctx, d, defaultTimeout)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsVulnerabilitiesV1CVEsExposedImages, functionName, "client_error", err)
		return nil, err
	}

	url := fmt.Sprintf("https://console.redhat.com/api/ocp-vulnerability/v1/cves/%s/exposed_images", cveName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsVulnerabilitiesV1CVEsExposedImages, functionName, "request_error", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsVulnerabilitiesV1CVEsExposedImages, functionName, "api_error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with status code %d", resp.StatusCode)
		pluginLogError(ctx, openshiftInsightsVulnerabilitiesV1CVEsExposedImages, functionName, "api_error", err)
		return nil, err
	}

	var exposedImagesResponse vulnerabilitiesV1CVEsExposedImagesResponse
	err = json.NewDecoder(resp.Body).Decode(&exposedImagesResponse)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsVulnerabilitiesV1CVEsExposedImages, functionName, "decode_error", err)
		return nil, err
	}

	for _, image := range exposedImagesResponse.Data {
		d.StreamListItem(ctx, image)
	}

	return nil, nil
}
