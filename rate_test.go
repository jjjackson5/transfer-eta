package main

import "testing"

func TestParseRate(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"25MB/s", 25e6},
		{"1024B/s", 1024},
		{"500B/sec", 500},
		{"1KB/second", 1e3},
		{"100bps", 100.0 / 8},
		{"100Mbps", 100e6 / 8},
		{"1.5Gbps", 1.5e9 / 8},
		{"  10MB/s  ", 10e6},
	}

	for _, c := range cases {
		got, err := parseRate(c.in)
		if err != nil {
			t.Errorf("parseRate(%q) returned error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("parseRate(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseRateErrors(t *testing.T) {
	cases := []string{
		"",
		"25MB",
		"25MB/min",
		"25MB/hour",
		"bps",
		"XMbps",
	}

	for _, in := range cases {
		if _, err := parseRate(in); err == nil {
			t.Errorf("parseRate(%q) expected an error, got nil", in)
		}
	}
}
