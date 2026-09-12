package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// bitRatePattern matches bits-per-second shorthand like "100Mbps",
// "1.5Gbps", or a bare "500bps". Network speeds are almost always quoted
// this way, distinct from the byte-based "25MB/s" form.
var bitRatePattern = regexp.MustCompile(`^([0-9]*\.?[0-9]+)\s*([KMGTPkmgtp]?)bps$`)

var bitUnitMultipliers = map[string]float64{
	"":  1,
	"k": 1e3,
	"m": 1e6,
	"g": 1e9,
	"t": 1e12,
	"p": 1e15,
}

// parseRate turns a string like "25MB/s" or "100Mbps" into bytes per
// second. Per-second rates are accepted for now — see README for what's
// not supported yet (per-minute rates and so on).
func parseRate(s string) (float64, error) {
	s = strings.TrimSpace(s)

	if m := bitRatePattern.FindStringSubmatch(s); m != nil {
		n, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid rate %q: %w", s, err)
		}
		bitsPerSec := n * bitUnitMultipliers[strings.ToLower(m[2])]
		return bitsPerSec / 8, nil
	}

	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid rate %q: expected format like 25MB/s or 100Mbps", s)
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
