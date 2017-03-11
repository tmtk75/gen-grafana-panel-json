package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

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

func NewTargetsSQS(opts *cloudWatchOpts, urls []string, prefix string, rp bool, exclude *regexp.Regexp) []Target {
	targets := make([]Target, 0)
	for i, q := range urls {
		qn := queueName(q)
		if rp {
			qn = removePrefix(prefix, qn)
		}
		if exclude != nil && exclude.Match([]byte(qn)) {
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
