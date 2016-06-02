package main

import (
	"os"

	// Import the snap plugin library
	"github.com/ticketmaster/snap-plugin-publisher-cloudwatch/cloudwatch"
	"github.com/intelsdi-x/snap/control/plugin"
)

func main() {
	// Three things provided:
	//   the definition of the plugin metadata
	//   the implementation satfiying plugin.CollectorPlugin
	//   the collector configuration policy satifying plugin.ConfigRules

	// Define metadata about Plugin
	meta := cloudwatch.Meta()

	// Start a collector
	plugin.Start(meta, cloudwatch.NewCloudWatchPublisher(), os.Args[1])
}
