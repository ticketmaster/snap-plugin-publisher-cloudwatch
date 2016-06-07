package cloudwatch

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	name       = "cloudwatch"
	version    = 1
	pluginType = plugin.PublisherPluginType
)

type cloudwatchPublisher struct{}

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

func NewCloudWatchPublisher() *cloudwatchPublisher {
	return &cloudwatchPublisher{}

}

func (p *cloudwatchPublisher) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()

	param1, err := cpolicy.NewStringRule("region", true)
	handleErr(err)
	param1.Description = "AWS Region"

	config.Add(param1)
	cp.Add([]string{""}, config)

	return cp, nil
}

func (p *cloudwatchPublisher) Publish(contentType string, content []byte, config map[string]ctypes.ConfigValue) error {
	logger := log.New()
	svc := cloudwatch.New(session.New(&aws.Config{Region: aws.String(config["region"].(ctypes.ConfigValueStr).Value)}))

	logger.Println("Publishing started")

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
		svc,
		logger,
	)

	return err
}

func publishDataToCloudWatch(metrics []plugin.MetricType, svc *cloudwatch.CloudWatch, logger *log.Logger) error {
	for _, m := range metrics {
		//logger.Println(strings.Join(m.Namespace().Strings(), "."))
		//logger.Println(m.Timestamp().String())
		//logger.Println(m.Data())

		input := &cloudwatch.PutMetricDataInput{
			MetricData: []*cloudwatch.MetricDatum{
				{
					//MetricName: aws.String(strings.Join(m.Namespace().Strings(), ".")),
					//Timestamp: aws.Time(m.Timestamp()),
					//Unit: aws.String("StandardUnit"),
					//Value: aws.Float64(m.Data().(float64)),
					MetricName: aws.String("MetricsName"),
					Timestamp: aws.Time(m.Timestamp()),
					Unit: aws.String("StandardUnit"),
					Value: aws.Float64(1.0),
				},
			},
			Namespace: aws.String("snap"),
		}

		logger.Printf(input)

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
