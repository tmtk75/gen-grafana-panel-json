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
