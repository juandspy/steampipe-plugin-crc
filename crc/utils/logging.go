package utils

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// LogErrorUsingSteampipeLogger logs an error using the steampipe logger
func LogErrorUsingSteampipeLogger(ctx context.Context, table, errType string, err error) {
	plugin.Logger(ctx).Error(table, errType, err)
}
