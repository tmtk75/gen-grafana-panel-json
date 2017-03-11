package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/jawher/mow.cli"
)

var (
	versionShort = "" // Expected given by -ldflags
	versionLong  = "" // Expected given by -ldflags
)

func main() {
	app := cli.App("gen-grafana-panel-json", "JSON Generator for Grafana CloudWatch datasource")
	app.Version("version", versionShort)
	app.Command("ec2", "EC2", func(c *cli.Cmd) {
		filters := c.String(cli.StringOpt{Name: "filters", Desc: `e.g: "tag:Name,dev-*", "instance-type,m3.large"`})
		opts := newCloudWatchOpts(c)
		c.Spec = "DATASOURCE_NAME [OPTIONS]"
		c.Action = func() {
			p := NewGrafanaPanel(*opts.dsName, "EC2 "+*opts.metricName)
			p.Targets = NewTargetsEC2(opts, *filters)
			PrintGrafanaPanelJSON(p)
		}
	})
	app.Command("sqs", "SQS", func(c *cli.Cmd) {
		opts := newCloudWatchOpts(c)
		px := c.String(cli.StringArg{Name: "PREFIX", Desc: "Prefix to filter"})
		rp := c.Bool(cli.BoolOpt{Name: "remove-prefix", Desc: "Prefix"})
		exc := c.String(cli.StringOpt{Name: "exclude", Desc: ""})
		c.Spec = "[OPTIONS] DATASOURCE_NAME PREFIX"
		c.Action = func() {
			var qs []string
			if terminal.IsTerminal(int(os.Stdin.Fd())) {
				qs = ListQueues(*opts.region, *px)
			} else {
				bytes, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					log.Fatalf("failed to read stdin: %v", err)
				}
				qs = strings.Split(strings.Trim(string(bytes), "\n"), "\n")
			}
			var re *regexp.Regexp
			if *exc != "" {
				re = regexp.MustCompile(*exc)
			}
			p := NewGrafanaPanel(*opts.dsName, fmt.Sprintf("SQS %v-* %v", *px, *opts.metricName))
			p.Targets = NewTargetsSQS(opts, qs, *px, *rp, re)
			PrintGrafanaPanelJSON(p)
		}
	})
	app.Command("list-queues", "SQS", func(c *cli.Cmd) {
		var (
			p = c.String(cli.StringArg{Name: "PREFIX", Desc: "Prefix to filter"})
			r = c.String(cli.StringArg{Name: "REGION", Value: "ap-northeast-1", Desc: "ap-northeast-1 by default"})
		)
		c.Spec = "PREFIX [REGION]"
		c.Action = func() {
			for _, q := range ListQueues(*r, *p) {
				fmt.Printf("%v\n", q)
			}
		}
	})
	app.Command("cloudwatch", "CloudWatch", func(c *cli.Cmd) {
		c.Command("list-metrics", "", func(c *cli.Cmd) {
			var (
				ns = c.String(cli.StringArg{Name: "NAMESPACE", Desc: "CloudWatch namespace e.g) AWS/EC2"})
				r  = c.String(cli.StringArg{Name: "REGION", Value: "ap-northeast-1", Desc: "ap-northeast-1 by default"})
			)
			c.Spec = "NAMESPACE [REGION]"
			c.Action = func() {
				ListMetrics(*ns, *r, nil)
			}
		})
	})
	app.Run(os.Args)
}

type cloudWatchOpts struct {
	dsName     *string
	metricName *string
	region     *string
	statistics *string
}

func newCloudWatchOpts(c *cli.Cmd) *cloudWatchOpts {
	return &cloudWatchOpts{
		dsName:     c.String(cli.StringArg{Name: "DATASOURCE_NAME", Desc: "Grafana datasource name"}),
		metricName: c.String(cli.StringOpt{Name: "metric-name m", Value: "CPUUtilization", Desc: "CloudWatch MetricName"}),
		region:     c.String(cli.StringOpt{Name: "region r", Value: "ap-northeast-1", Desc: "AWS region"}),
		statistics: c.String(cli.StringOpt{Name: "statistics s", Value: "Average", Desc: "e.g: Average,Maximum,Minimum,Sum,SampleCount"}),
	}
}
