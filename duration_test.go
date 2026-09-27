package main

import "testing"

func TestParseDuration(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"90", 90},
		{"90s", 90},
		{"1h30m", 5400},
		{"3d4h", 3*86400 + 4*3600},
		{"1h30m15s", 5415},
		{"0.5h", 1800},
		{"  90s  ", 90},
	}

	for _, c := range cases {
		got, err := parseDuration(c.in)
		if err != nil {
			t.Errorf("parseDuration(%q) returned error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("parseDuration(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseDurationErrors(t *testing.T) {
	cases := []string{
		"",
		"1h30x",
		"1h1h",
		"abc",
		"h",
	}

	for _, in := range cases {
		if _, err := parseDuration(in); err == nil {
			t.Errorf("parseDuration(%q) expected an error, got nil", in)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0s"},
		{5, "5s"},
		{65, "1m5s"},
		{3661, "1h1m1s"},
		{90061, "1d1h1m1s"},
	}

	for _, c := range cases {
		got := formatDuration(c.in)
		if got != c.want {
			t.Errorf("formatDuration(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}
