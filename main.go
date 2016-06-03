package main

import (
	"os"

	"github.com/ticketmaster/snap-plugin-publisher-cloudwatch/cloudwatch"
	"github.com/intelsdi-x/snap/control/plugin"
)

func main() {
	meta := cloudwatch.Meta()
	plugin.Start(meta, cloudwatch.NewCloudWatchPublisher(), os.Args[1])
}
