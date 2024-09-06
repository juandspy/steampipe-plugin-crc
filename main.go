package main

import (
	"github.com/juandspy/steampipe-plugin-crc/crc"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: crc.Plugin})
}
