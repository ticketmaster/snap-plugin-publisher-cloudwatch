package cloudwatch

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
	"time"
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
	svc := cloudwatch.New(session.New())

	var metrics []plugin.MetricType

	switch contentType {
	case plugin.SnapGOBContentType:
		dec := gob.NewDecoder(bytes.NewBuffer(content))
		if err := dec.Decode(&metrics); err != nil {
			logger.Printf("Error decoding: error=%v content=%v", err, content)
			return err
		}
	case plugin.SnapJSONContentType:
		err := json.Unmarshal(content, &metrics)
		if err != nil {
			logger.Printf("Error decoding JSON: error=%v content=%v", err, content)
			return err
		}
	default:
		logger.Printf("Error unknown content type '%v'", contentType)
		return fmt.Errorf("Unknown content type '%s'", contentType)
	}

	err := publishDataToCloudWatch(
		metrics,
		svc,
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

func publishDataToCloudWatch(metrics []plugin.MetricType, svc *cloudwatch.CloudWatch, logger *log.Logger) error {
	for _, m := range metrics {
		input := &cloudwatch.PutMetricDataInput{
			MetricData: []*cloudwatch.MetricDatum{
				{
					MetricName: aws.String(m.Namespace()),
					Timestamp: aws.Time(m.Timestamp()),
					Unit: aws.String("StandardUnit"),
					Value: aws.Float64(m.Data()),
				},
			},
			Namespace: aws.String("snap"),
		}

		_, err := svc.PutMetricData(input)
		if err != nil {
			handleErr(err)
		}
	}

	return nil
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}
