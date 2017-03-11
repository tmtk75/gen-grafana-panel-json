package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func ListMetrics(ns, region string, nextToken *string) {
	svc := cloudwatch.New(session.New(), &aws.Config{Region: aws.String(region)})
	req := cloudwatch.ListMetricsInput{
		Namespace:  aws.String(ns),
		Dimensions: []*cloudwatch.DimensionFilter{
		//&cloudwatch.DimensionFilter{
		//	Name:  aws.String("QueueName"),
		//	Value: aws.String("stg-jp"),
		//},
		},
		NextToken: nextToken,
	}
	res, err := svc.ListMetrics(&req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	for _, m := range res.Metrics {
		fmt.Printf("%v %v\n", *m.Dimensions[0].Value, *m.MetricName)
	}
	if res.NextToken != nil && *res.NextToken != "" {
		ListMetrics(ns, region, res.NextToken)
	}
}
