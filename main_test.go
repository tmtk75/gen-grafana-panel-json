package main

import "testing"

func TestAlias(t *testing.T) {
	ts := []struct {
		name string
		exp  string
	}{
		{name: "ApproximateAgeOfOldestMessage", exp: "AAOM"},
		{name: "ApproximateNumberOfMessagesVisible", exp: "ANMV"},
	}
	for _, c := range ts {
		r := alias(c.name)
		if r != c.exp {
			t.Fatalf("expect: %v, but actual: %v for %v", c.exp, r, c.name)
		}
	}
}

func TestQueueName(t *testing.T) {
	ts := []struct {
		url string
		exp string
	}{
		{url: "https://sqs.ap-northeast-1.amazonaws.com/112646608144/broker_report-0082", exp: "broker_report-0082"},
	}
	for _, c := range ts {
		r := queueName(c.url)
		if r != c.exp {
			t.Fatalf("expect: %v, but actual: %v for %v", c.exp, r, c.url)
		}
	}
}

func TestRemovePrefix(t *testing.T) {
	ts := []struct {
		prefix string
		input  string
		exp    string
	}{
		{prefix: "dev-", input: "dev-a", exp: "a"},
		{prefix: "stg-", input: "stg-b", exp: "b"},
		{prefix: "aaa-", input: "prd-c", exp: "prd-c"},
		{prefix: "aaa-", input: "aaa-abc", exp: "abc"},
	}
	for _, c := range ts {
		r := removePrefix(c.prefix, c.input)
		if r != c.exp {
			t.Fatalf("expect: %v, but actual: %v for %v", c.exp, r, c.input)
		}
	}
}
