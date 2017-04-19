package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jawher/mow.cli"
)

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
