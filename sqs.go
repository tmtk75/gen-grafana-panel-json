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
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jawher/mow.cli"
)

type SQS struct {
	*cloudwatchOpts
	queues  []string
	prefix  string
	remove  bool
	exclude *regexp.Regexp
}

func sqsCmd(c *cli.Cmd) {
	opts := newCloudwatchOpts(c)
	px := c.String(cli.StringArg{Name: "PREFIX", Desc: "Prefix to filter"})
	rp := c.Bool(cli.BoolOpt{Name: "remove-prefix", Desc: "Remove prefix in display if true"})
	exc := c.String(cli.StringOpt{Name: "exclude", Desc: "Regex for name to exclude"})
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
		sqs := SQS{cloudwatchOpts: opts, queues: qs, prefix: *px, remove: *rp, exclude: re}
		p.Targets = sqs.NewTargets()
		PrintGrafanaPanelJSON(p)
	}
}

func listQueueCmd(c *cli.Cmd) {
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
}

func ListQueues(region, prefix string) []string {
	svc := sqs.New(session.New(), &aws.Config{Region: aws.String(region)})
	req := sqs.ListQueuesInput{
		QueueNamePrefix: aws.String(prefix),
	}
	res, err := svc.ListQueues(&req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	qs := make([]string, 0)
	for _, q := range res.QueueUrls {
		qs = append(qs, *q)
	}
	return qs
}

func queueName(q string) string {
	s := strings.Split(q, "/")
	return s[len(s)-1] // queue name
}

func removePrefix(prefix, s string) string {
	return strings.Replace(s, prefix, "", 1)
}

func (e *SQS) NewTargets() []Target {
	opts := e.cloudwatchOpts
	targets := make([]Target, 0)
	for i, q := range e.queues {
		qn := queueName(q)
		if e.remove {
			qn = removePrefix(e.prefix, qn)
		}
		if e.exclude != nil && e.exclude.Match([]byte(qn)) {
			continue
		}
		t := Target{
			Dimensions: map[string]string{
				"QueueName": queueName(q),
			},
			MetricName: *opts.metricName,
			Namespace:  "AWS/SQS",
			Period:     "300",
			RefID:      fmt.Sprintf("ID-%v", i),
			Region:     *opts.region,
			Statistics: []string{
				*opts.statistics,
			},
			//Alias: alias(*opts.metricName) + ":" + qn,
			Alias: qn,
		}
		targets = append(targets, t)
	}
	return targets //[0:20]
}
