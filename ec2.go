package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	cli "github.com/jawher/mow.cli"
)

type EC2 struct {
	*cloudwatchOpts
	filters []*ec2.Filter
}

func ec2Cmd(c *cli.Cmd) {
	fs := c.String(cli.StringOpt{Name: "filters", Desc: `e.g: "tag:Name,dev-*", "instance-type,m3.large"`})
	opts := newCloudwatchOpts(c)
	c.Spec = "DATASOURCE_NAME [OPTIONS]"
	c.Action = func() {
		p := NewGrafanaPanel(*opts.dsName, "EC2 "+*opts.metricName)
		ec2 := EC2{cloudwatchOpts: opts, filters: parseFilters(*fs)}
		p.Targets = ec2.NewTargets()
		PrintGrafanaPanelJSON(p)
	}
}

func (e *EC2) NewTargets() []Target {
	opts := e.cloudwatchOpts
	f := []*ec2.Filter{
		&ec2.Filter{
			Name:   aws.String("instance-state-name"),
			Values: []*string{aws.String("running")},
		},
	}
	f = append(f, e.filters...)

	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(*opts.region)})
	req := ec2.DescribeInstancesInput{Filters: f}
	res, err := svc.DescribeInstances(&req)
	if err != nil {
		log.Fatalf("failed to describe-instances: %v", err)
		log.Fatalf("failed to DescribeInstances: %v", err)
	}

	ref := 0
	targets := make([]Target, 0)
	for _, res := range res.Reservations {
		//fmt.Println(len(res.Instances))
		for _, i := range res.Instances {
			alias := ""
			for _, t := range i.Tags {
				if *t.Key == "Name" {
					alias = *t.Value
				}
			}
			targets = append(targets, Target{
				Dimensions: map[string]string{"InstanceId": *i.InstanceId},
				MetricName: *opts.metricName,
				Namespace:  "AWS/EC2",
				Period:     "",
				Region:     *opts.region,
				Statistics: []string{
					*opts.statistics,
				},
				RefID: fmt.Sprintf("A%d", ref),
				Alias: alias,
			})
			ref += 1
		}
	}

	sort.Sort(Targets(targets))
	return targets
}

func parseFilters(filters string) []*ec2.Filter {
	f := []*ec2.Filter{}
	fs := strings.Split(filters, ",")
	if len(fs)%2 == 1 {
		log.Fatalln("illegal filters option: it should be even")
	}
	for i := 0; i*2 < len(fs); i++ {
		f = append(f, &ec2.Filter{
			Name:   aws.String(strings.TrimSpace(fs[i*2])),
			Values: []*string{aws.String(strings.TrimSpace(fs[i*2+1]))},
		})
	}
	return f
}
