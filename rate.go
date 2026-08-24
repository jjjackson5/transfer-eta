package main

import (
	"fmt"
	"strings"
)

// parseRate turns a string like "25MB/s" into bytes per second. Only
// per-second rates are accepted for now — see README for what's not
// supported yet (Mbps/Gbps, per-minute rates, and so on).
func parseRate(s string) (float64, error) {
	s = strings.TrimSpace(s)
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid rate %q: expected format like 25MB/s", s)
	}

	perUnit := strings.ToLower(strings.TrimSpace(parts[1]))
	if perUnit != "s" && perUnit != "sec" && perUnit != "second" {
		return 0, fmt.Errorf("invalid rate %q: only per-second rates are supported (e.g. 25MB/s)", s)
	}

	bytesPerSec, err := parseSize(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid rate %q: %w", s, err)
	}

	return bytesPerSec, nil
}
