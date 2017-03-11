package main

import (
	"os"

	"github.com/jawher/mow.cli"
)

var (
	versionShort = "" // Expected given by -ldflags
	versionLong  = "" // Expected given by -ldflags
)

func main() {
	app := cli.App("gen-grafana-panel-json", "JSON Generator for Grafana CloudWatch datasource")
	app.Version("version", versionShort)
	app.Command("ec2", "Generate JSON for AWS/EC2", ec2Cmd)
	app.Command("sqs", "Generate JSON for AWS/SQS", sqsCmd)
	app.Command("list-queues", "List queue names for SQS", listQueueCmd)
	app.Command("cloudwatch", "CloudWatch commands", cloudwatchCmd)
	app.Run(os.Args)
}

type cloudwatchOpts struct {
	dsName     *string
	metricName *string
	region     *string
	statistics *string
}

func newCloudwatchOpts(c *cli.Cmd) *cloudwatchOpts {
	return &cloudwatchOpts{
		dsName:     c.String(cli.StringArg{Name: "DATASOURCE_NAME", Desc: "Grafana datasource name"}),
		metricName: c.String(cli.StringOpt{Name: "metric-name m", Value: "CPUUtilization", Desc: "CloudWatch MetricName"}),
		region:     c.String(cli.StringOpt{Name: "region r", Value: "ap-northeast-1", Desc: "AWS region"}),
		statistics: c.String(cli.StringOpt{Name: "statistics s", Value: "Average", Desc: "e.g: Average,Maximum,Minimum,Sum,SampleCount"}),
	}
}
