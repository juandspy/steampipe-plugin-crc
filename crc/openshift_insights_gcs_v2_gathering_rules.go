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

const openshiftInsightsGCSV2RemoteConfiguration = "openshift_insights_gcs_v2_gathering_rules"

type gatheringRulesV2 struct {
	Version                   string        `json:"version,omitempty"`
	ConditionalGatheringRules []interface{} `json:"conditional_gathering_rules,omitempty"`
	ContainerLogs             interface{}   `json:"container_logs,omitempty"`
}

func tableInsightsGatheringRulesV2(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        openshiftInsightsGCSV2RemoteConfiguration,
		Description: "Return the gathering rules for a given OCP version.",
		Get: &plugin.GetConfig{
			Hydrate:    getGatheringRulesV2,
			KeyColumns: plugin.SingleColumn("ocp_version"),
		},
		Columns: []*plugin.Column{
			{
				Name:        "ocp_version",
				Type:        proto.ColumnType_STRING,
				Description: "Cluster version.",
				Transform:   transform.FromQual("ocp_version"),
			},
			{
				Name:        "version",
				Type:        proto.ColumnType_STRING,
				Description: "Gathering rules version.",
				// Transform:   transform.FromField("version"),
			},
			{
				Name:        "conditional_gathering_rules",
				Type:        proto.ColumnType_JSON,
				Description: "The conditions that trigger the gathering functions.",
				// Transform:   transform.FromField("conditional_gathering_rules"),
			},
			{
				Name:        "container_logs",
				Type:        proto.ColumnType_JSON,
				Description: "The container logs filtering.",
				// Transform:   transform.FromField("container_logs"),
			},
		},
	}
}

// getGatheringRulesV2 populates the table with all the gathering rules for a specific OCP version
func getGatheringRulesV2(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	const functionName = "getGatheringRulesV2"

	// get the OCP Version
	ocpVersion := d.EqualsQualString("ocp_version")

	if ocpVersion == "" {
		err := errors.New("you must specify an OCP version")
		pluginLogError(ctx, openshiftInsightsGCSV2RemoteConfiguration, functionName, "query_error", err)
		return nil, err
	}

	client, err := connect(ctx, d)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsGCSV2RemoteConfiguration, functionName, "client_error", err)
		return nil, err
	}

	url := fmt.Sprintf("https://console.redhat.com/api/gathering/v2/%s/gathering_rules", ocpVersion)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsGCSV2RemoteConfiguration, functionName, "request_error", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsGCSV2RemoteConfiguration, functionName, "api_error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		pluginLogError(ctx, openshiftInsightsGCSV2RemoteConfiguration, functionName, "api_error", err)
		return nil, err
	}

	rules, err := decodeGatheringRulesV2(resp.Body)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsGCSV2RemoteConfiguration, functionName, "decode_error", err)
		return nil, err
	}

	plugin.Logger(ctx).Warn(functionName, "rules", rules)

	return rules, nil
}

func decodeGatheringRulesV2(body io.ReadCloser) (gatheringRulesV2, error) {
	var rules gatheringRulesV2
	err := json.NewDecoder(body).Decode(&rules)
	return rules, err
}
