package crc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const openshiftInsightsGCSV1GatheringRules = "openshift_insights_gcs_v1_gathering_rules"

type gatheringRules struct {
	Version string `json:"version"`
	Rules   []struct {
		Conditions         []interface{} `json:"conditions"`
		GatheringFunctions interface{}   `json:"gathering_functions"`
	} `json:"rules"`
}

// Define the table
func tableInsightsGatheringRules(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        openshiftInsightsGCSV1GatheringRules,
		Description: "Return the gathering rules for a given version.",
		List: &plugin.ListConfig{
			Hydrate: listGatheringRules,
		},
		Columns: []*plugin.Column{
			{
				Name:        "version",
				Type:        proto.ColumnType_STRING,
				Description: "Gathering rules version.",
				Transform:   transform.FromField("version"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "The conditions that trigger the gathering functions.",
				Transform:   transform.FromField("conditions"),
			},
			{
				Name:        "gathering_functions",
				Type:        proto.ColumnType_JSON,
				Description: "The gathering mechanisms.",
				Transform:   transform.FromField("gathering_functions"),
			},
		},
	}
}

// listGatheringRules populates the openshift_insights_gathering_conditions_service table with all the gathering rules in the API
func listGatheringRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	const functionName = "listGatheringRules"
	client, err := connect(ctx, d)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsGCSV1GatheringRules, functionName, "client_error", err)
		return nil, err
	}

	url := "https://console.redhat.com/api/gathering/v1/gathering_rules"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsGCSV1GatheringRules, functionName, "request_error", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsGCSV1GatheringRules, functionName, "api_error", err)
		return nil, err
	}
	defer resp.Body.Close()

	rules, err := decodeGatheringRules(resp.Body)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsGCSV1GatheringRules, functionName, "decode_error", err)
		return nil, err
	}

	for _, rule := range rules.Rules {
		row := map[string]interface{}{}
		row["version"] = rules.Version
		row["conditions"] = rule.Conditions
		row["gathering_functions"] = rule.GatheringFunctions
		d.StreamListItem(ctx, row)
	}

	return nil, nil
}

func decodeGatheringRules(body io.ReadCloser) (gatheringRules, error) {
	var rules gatheringRules
	err := json.NewDecoder(body).Decode(&rules)
	return rules, err
}
