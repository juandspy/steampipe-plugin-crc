package gathering_conditions_service

import (
	"context"
	"encoding/json"
	"io"

	"github.com/juandspy/steampipe-plugin-crc/crc/utils"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const V1GatheringRulesTableName = "crc_openshift_insights_gcs_v1_gathering_rules"

type gatheringRulesV1 struct {
	Version string `json:"version"`
	Rules   []struct {
		Conditions         []interface{} `json:"conditions"`
		GatheringFunctions interface{}   `json:"gathering_functions"`
	} `json:"rules"`
}

func TableGatheringRulesV1(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        V1GatheringRulesTableName,
		Description: "Return a list of versioned gathering rules.",
		List: &plugin.ListConfig{
			Hydrate: listGatheringRulesV1,
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

// listGatheringRulesV1 populates the table with all the gathering rules in the API
func listGatheringRulesV1(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	endpoint := "api/gathering/v1/gathering_rules"
	resp, err := utils.MakeAPIRequest(ctx, d, "GET", endpoint, nil, utils.DefaultTimeout)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V1GatheringRulesTableName, "api_error", err)
		return nil, err
	}
	defer resp.Body.Close()

	rules, err := decodeGatheringRulesV1(resp.Body)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V1GatheringRulesTableName, "decode_error", err)
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

func decodeGatheringRulesV1(body io.ReadCloser) (gatheringRulesV1, error) {
	var rules gatheringRulesV1
	err := json.NewDecoder(body).Decode(&rules)
	return rules, err
}
