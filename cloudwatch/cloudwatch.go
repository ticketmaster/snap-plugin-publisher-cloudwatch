package cloudwatch

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
)

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

type cloudwatchPublisher struct{}

func NewCloudWatchPublisher() *cloudwatchPublisher {
	return &cloudwatchPublisher{}

}

const (
	name       = "cloudwatch"
	version    = 1
	pluginType = plugin.PublisherPluginType
)

func (rmq *cloudwatchPublisher) Publish(contentType string, content []byte, config map[string]ctypes.ConfigValue) error {
	logger := log.New()
	var metrics []plugin.MetricType
	switch contentType {
	case plugin.SnapGOBContentType:
		dec := gob.NewDecoder(bytes.NewBuffer(content))
		if err := dec.Decode(&metrics); err != nil {
			logger.Printf("Error decoding: error=%v content=%v", err, content)
			return err
		}
	default:
		logger.Printf("Error unknown content type '%v'", contentType)
		return fmt.Errorf("Unknown content type '%s'", contentType)
	}
	err := publishDataToCloudWatch(
		metrics,
		logger,
	)
	return err
}

func (rmq *cloudwatchPublisher) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()

	cp.Add([]string{""}, config)
	return cp, nil
}

func publishDataToCloudWatch(metrics []plugin.MetricType, logger *log.Logger) error {

	return nil
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}
