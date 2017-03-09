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
