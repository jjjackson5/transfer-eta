package main

import "testing"

func TestParseSize(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"1024", 1024},
		{"1024B", 1024},
		{"1KB", 1e3},
		{"4.7GB", 4.7e9},
		{"650MiB", 650 * 1024 * 1024},
		{"1KIB", 1024},
		{"1PiB", 1024 * 1024 * 1024 * 1024 * 1024},
		{"  2MB  ", 2e6},
		{"0", 0},
	}

	for _, c := range cases {
		got, err := parseSize(c.in)
		if err != nil {
			t.Errorf("parseSize(%q) returned error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("parseSize(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseSizeErrors(t *testing.T) {
	cases := []string{
		"",
		"GB",
		"5XB",
		"five GB",
		"-5GB",
	}

	for _, in := range cases {
		if _, err := parseSize(in); err == nil {
			t.Errorf("parseSize(%q) expected an error, got nil", in)
		}
	}
}

func TestFormatSize(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0 B"},
		{999, "999 B"},
		{1000, "1.00 KB"},
		{1500000, "1.50 MB"},
		{4700000000, "4.70 GB"},
	}

	for _, c := range cases {
		got := formatSize(c.in)
		if got != c.want {
			t.Errorf("formatSize(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}
