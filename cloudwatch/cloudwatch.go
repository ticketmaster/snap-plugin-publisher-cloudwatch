/*
The MIT License (MIT)

Copyright (c) 2016 Ticketmaster Ltd

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package cloudwatch

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
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

	param2, err := cpolicy.NewStringRule("namespace", true)
	handleErr(err)
	param2.Description = "Metrics Namespace"
	config.Add(param2)

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
		config,
		logger,
	)

	return err
}

func getCloudwatchMetricDimension(m plugin.MetricType, logger *log.Logger) []*cloudwatch.Dimension {
	tags := m.Tags()

	dimensions := make([]*cloudwatch.Dimension, len(tags), len(tags))
	index := 0
	for k, v := range tags {

		dimensions[index] = &cloudwatch.Dimension{
			Name:  aws.String(k),
			Value: aws.String(v),
		}

		index++
	}

	return dimensions
}

func publishDataToCloudWatch(metrics []plugin.MetricType, svc *cloudwatch.CloudWatch, config map[string]ctypes.ConfigValue, logger *log.Logger) error {
	namespace := config["namespace"].(ctypes.ConfigValueStr).Value

	for _, m := range metrics {
		valueString := fmt.Sprintf("%v", m.Data())
		cloudwatchMetricValue, err := strconv.ParseFloat(valueString, 64)
		if err == nil {
			logger.Printf("Converted Value: %s=%f", m.Namespace().String(), cloudwatchMetricValue)
			cloudwatchMetricDimension := getCloudwatchMetricDimension(m, logger)

			input := &cloudwatch.PutMetricDataInput{
				MetricData: []*cloudwatch.MetricDatum{
					{
						MetricName: aws.String(strings.Join(m.Namespace().Strings(), ".")),
						Timestamp:  aws.Time(m.Timestamp()),
						Value:      aws.Float64(cloudwatchMetricValue),
						Dimensions: cloudwatchMetricDimension,
					},
				},
				Namespace: aws.String(namespace),
			}

			_, err := svc.PutMetricData(input)
			if err != nil {
				handleErr(err)
			}
		} else {
			logger.Println("Unable to convert string onto float value: ", err, valueString)
		}
	}

	return nil
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}
