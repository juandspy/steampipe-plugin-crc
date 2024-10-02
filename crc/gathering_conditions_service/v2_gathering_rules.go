package gathering_conditions_service

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

const V2RemoteConfigurationTableName = "crc_openshift_insights_gcs_v2_gathering_rules"

type gatheringRulesV2 struct {
	Version                   string        `json:"version,omitempty"`
	ConditionalGatheringRules []interface{} `json:"conditional_gathering_rules,omitempty"`
	ContainerLogs             interface{}   `json:"container_logs,omitempty"`
}

func TableGatheringRulesV2(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        V2RemoteConfigurationTableName,
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
	// get the OCP Version
	ocpVersion := d.EqualsQualString("ocp_version")

	if ocpVersion == "" {
		err := errors.New("you must specify an OCP version")
		utils.LogErrorUsingSteampipeLogger(ctx, V2RemoteConfigurationTableName, "query_error", err)
		return nil, err
	}
	endpoint := fmt.Sprintf("api/gathering/v2/%s/gathering_rules", ocpVersion)
	resp, err := utils.MakeAPIRequest(ctx, d, "GET", endpoint, nil, utils.DefaultTimeout)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V2RemoteConfigurationTableName, "api_error", err)
		return nil, err
	}

	defer resp.Body.Close()

	rules, err := decodeGatheringRulesV2(resp.Body)
	if err != nil {
		utils.LogErrorUsingSteampipeLogger(ctx, V2RemoteConfigurationTableName, "decode_error", err)
		return nil, err
	}

	return rules, nil
}

func decodeGatheringRulesV2(body io.ReadCloser) (gatheringRulesV2, error) {
	var rules gatheringRulesV2
	err := json.NewDecoder(body).Decode(&rules)
	return rules, err
}
