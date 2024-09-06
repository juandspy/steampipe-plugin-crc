package crc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const openshiftInsightsAggregatorV2ClusterReports = "openshift_insights_aggregator_v2_cluster_reports"

type ClusterReportsResponseV2 struct {
	Report struct {
		Meta struct {
			ClusterName   string    `json:"cluster_name"`
			Managed       bool      `json:"managed"`
			Count         int       `json:"count"`
			LastCheckedAt time.Time `json:"last_checked_at"`
			GatheredAt    time.Time `json:"gathered_at"`
		} `json:"meta"`
		Data []struct {
			ClusterID       string    `json:"cluster_name,omitempty"` // added manually
			RuleID          string    `json:"rule_id"`
			CreatedAt       time.Time `json:"created_at"`
			Description     string    `json:"description"`
			Details         string    `json:"details"`
			Reason          string    `json:"reason"`
			Resolution      string    `json:"resolution"`
			MoreInfo        string    `json:"more_info"`
			TotalRisk       int       `json:"total_risk"`
			Disabled        bool      `json:"disabled"`
			DisableFeedback string    `json:"disable_feedback"`
			DisabledAt      string    `json:"disabled_at"`
			Internal        bool      `json:"internal"`
			UserVote        int       `json:"user_vote"`
			ExtraData       struct {
				ErrorKey      string   `json:"error_key"`
				InvalidInfras []string `json:"invalid_infras"`
				OcpVersion    string   `json:"ocp_version"`
				Type          string   `json:"type"`
			} `json:"extra_data"`
			Tags     []string  `json:"tags"`
			Impacted time.Time `json:"impacted"`
		} `json:"data"`
	} `json:"report"`
	Status string `json:"status"`
}

func tableAggregatorClusterReportsV2(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_insights_aggregator_v2_cluster_reports",
		Description: "Returns the latest report for the given cluster.",
		List: &plugin.ListConfig{
			Hydrate:    listClusterReportsV2,
			KeyColumns: plugin.SingleColumn("cluster_id"),
		},
		Columns: []*plugin.Column{
			{
				Name:        "cluster_id",
				Type:        proto.ColumnType_STRING,
				Description: "Cluster ID.",
				Transform:   transform.FromQual("cluster_id"),
			},
			{
				Name:        "rule_id",
				Type:        proto.ColumnType_STRING,
				Description: "Unique identifier for the rule.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time when the report was created.",
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "Description of the report.",
			},
			{
				Name:        "details",
				Type:        proto.ColumnType_STRING,
				Description: "Details about the report.",
			},
			{
				Name:        "reason",
				Type:        proto.ColumnType_STRING,
				Description: "Reason for the report.",
			},
			{
				Name:        "resolution",
				Type:        proto.ColumnType_STRING,
				Description: "Resolution of the issue described in the report.",
			},
			{
				Name:        "more_info",
				Type:        proto.ColumnType_STRING,
				Description: "Additional information related to the report.",
			},
			{
				Name:        "total_risk",
				Type:        proto.ColumnType_INT,
				Description: "Total risk score associated with the report.",
			},
			{
				Name:        "disabled",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if the report is disabled.",
			},
			{
				Name:        "disable_feedback",
				Type:        proto.ColumnType_STRING,
				Description: "Feedback on why the report was disabled.",
			},
			{
				Name:        "disabled_at",
				Type:        proto.ColumnType_STRING,
				Description: "Timestamp when the report was disabled.",
			},
			{
				Name:        "internal",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if the report is internal.",
			},
			{
				Name:        "user_vote",
				Type:        proto.ColumnType_INT,
				Description: "User vote on the report.",
			},
			{
				Name:        "extra_data",
				Type:        proto.ColumnType_JSON,
				Description: "Extra data.",
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: "Tags associated with the report.",
			},
			{
				Name:        "impacted",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Time when the issue impacted the cluster.",
			},
		},
	}
}

func listClusterReportsV2(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	const functionName = "listClusterReportsV2"

	// get the cluster ID
	clusterID := d.EqualsQualString("cluster_id")

	if clusterID == "" {
		err := errors.New("you must specify an OCP version")
		pluginLogError(ctx, openshiftInsightsAggregatorV2ClusterReports, functionName, "query_error", err)
		return nil, err
	}

	client, err := connect(ctx, d, defaultTimeout)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsAggregatorV2ClusterReports, functionName, "client_error", err)
		return nil, err
	}

	url := fmt.Sprintf("https://console.redhat.com/api/insights-results-aggregator/v2/cluster/%s/reports", clusterID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsAggregatorV2ClusterReports, functionName, "request_error", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsAggregatorV2ClusterReports, functionName, "api_error", err)
		return nil, err
	}

	defer resp.Body.Close()

	clusterReportsResponse, err := decodeClusterReportsResponseV2(resp.Body)
	if err != nil {
		pluginLogError(ctx, openshiftInsightsAggregatorV2ClusterReports, functionName, "decode_error", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = errors.New(clusterReportsResponse.Status)
		pluginLogError(ctx, openshiftInsightsAggregatorV2ClusterReports, functionName, "api_error", err)
		return nil, err
	}

	for _, report := range clusterReportsResponse.Report.Data {
		report.ClusterID = clusterID
		d.StreamListItem(ctx, report)
	}

	return nil, nil

}

func decodeClusterReportsResponseV2(body io.ReadCloser) (ClusterReportsResponseV2, error) {
	var clusterReportsResponse ClusterReportsResponseV2
	err := json.NewDecoder(body).Decode(&clusterReportsResponse)
	return clusterReportsResponse, err
}
