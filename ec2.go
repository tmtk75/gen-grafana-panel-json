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
	prefix  string
}

func ec2Cmd(c *cli.Cmd) {
	var (
		fs = c.String(cli.StringOpt{Name: "filters", Desc: `e.g: "tag:Name,dev-*", "instance-type,m3.large", "instance-state-name,running"`})
		rp = c.String(cli.StringOpt{Name: "remove-prefix", Desc: "Prefix to remove in alias"})
	)
	opts := newCloudwatchOpts(c)
	c.Spec = "[OPTIONS] DATASOURCE_NAME"
	c.Action = func() {
		p := NewGrafanaPanel(*opts.dsName, "EC2 "+*opts.metricName)
		ec2 := EC2{cloudwatchOpts: opts, filters: parseFilters(*fs), prefix: *rp}
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

	cfg := &aws.Config{Region: aws.String(*opts.region)}
	svc := ec2.New(session.New(), cfg)
	req := ec2.DescribeInstancesInput{Filters: f}
	res, err := svc.DescribeInstances(&req)
	if err != nil {
		log.Fatalf("failed to describe-instances: %v %v", err, cfg.CredentialsChainVerboseErrors)
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
				alias = strings.Replace(alias, e.prefix, "", 1)
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
	fs := []string{}
	if filters != "" {
		fs = strings.Split(filters, ",")
		if len(fs)%2 == 1 {
			log.Fatalf("illegal filters option: it should be even. %q\n", filters)
		}
	}
	for i := 0; i*2 < len(fs); i++ {
		f = append(f, &ec2.Filter{
			Name:   aws.String(strings.TrimSpace(fs[i*2])),
			Values: []*string{aws.String(strings.TrimSpace(fs[i*2+1]))},
		})
	}
	return f
}
