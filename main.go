package main

import (
	"os"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/ticketmaster/snap-plugin-publisher-cloudwatch/cloudwatch"
)

func main() {
	meta := cloudwatch.Meta()
	plugin.Start(meta, cloudwatch.NewCloudWatchPublisher(), os.Args[1])
}
