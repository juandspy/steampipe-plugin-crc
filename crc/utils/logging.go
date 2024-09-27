package utils

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// LogErrorUsingSteampipeLogger logs an error using the steampipe logger
func LogErrorUsingSteampipeLogger(ctx context.Context, table, function, errType string, err error) {
	plugin.Logger(ctx).Error(fmt.Sprintf("%s.%s", table, function), errType, err)
}
