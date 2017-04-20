package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jawher/mow.cli"
)

type DynamoDB struct {
	*cloudwatchOpts
	prefix string
	tables []string
}

func dynamodbCmd(c *cli.Cmd) {
	var (
		opts = newCloudwatchOpts(c)
		px   = c.String(cli.StringArg{Name: "PREFIX", Desc: "Prefix to filter"})
	)
	c.Spec = "[OPTIONS] DATASOURCE_NAME [PREFIX]"
	c.Action = func() {
		var names []string
		if terminal.IsTerminal(int(os.Stdin.Fd())) {
			names = ListTables(*opts.region)
		} else {
			bytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("failed to read stdin: %v", err)
			}
			names = strings.Split(strings.Trim(string(bytes), "\n"), "\n")
		}

		ns := make([]string, 0)

		var re *regexp.Regexp
		if *px != "" {
			re = regexp.MustCompile("^" + *px)
			for _, n := range names {
				if re.MatchString(n) {
					ns = append(ns, n)
				}
			}
		} else {
			ns = names
		}

		p := NewGrafanaPanel(*opts.dsName, fmt.Sprintf("DynamoDB %v* %v %v", *px, *opts.metricName, *opts.statistics))
		dynamo := DynamoDB{cloudwatchOpts: opts, prefix: *px, tables: ns}
		p.Targets = dynamo.NewTargets()
		PrintGrafanaPanelJSON(p)
	}
}

func listTableCmd(c *cli.Cmd) {
	var (
		r = c.String(cli.StringArg{Name: "REGION", Value: "ap-northeast-1", Desc: "ap-northeast-1 by default"})
	)
	c.Spec = "[REGION]"
	c.Action = func() {
		for _, q := range ListTables(*r) {
			fmt.Printf("%v\n", q)
		}
	}
}

func ListTables(region string) []string {
	svc := dynamodb.New(session.New(), &aws.Config{Region: aws.String(region)})
	req := dynamodb.ListTablesInput{}

	res, err := svc.ListTables(&req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	ns := make([]string, 0)
	for _, t := range res.TableNames {
		ns = append(ns, *t)
	}
	return ns
}

func (e *DynamoDB) NewTargets() []Target {
	opts := e.cloudwatchOpts
	targets := make([]Target, 0)
	for i, n := range e.tables {
		t := Target{
			Dimensions: map[string]string{
				"TableName": n,
			},
			MetricName: *opts.metricName,
			Namespace:  "AWS/DynamoDB",
			Period:     "60",
			RefID:      fmt.Sprintf("ID-%v", i),
			Region:     *opts.region,
			Statistics: []string{
				*opts.statistics,
			},
			Alias: strings.TrimPrefix(n, e.prefix),
		}
		targets = append(targets, t)
	}
	return targets
}
