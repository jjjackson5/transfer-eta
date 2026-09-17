package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// formatDuration renders a number of seconds as a compact human string
// like "1h2m3s" or "3d4h". time.Duration isn't used for that range because
// transfer times can run into weeks at low enough rates, well past what
// most people want to see expressed in nanosecond-precision Duration units.
//
// Below one second, rounding to the nearest second would just print "0s"
// for every fast transfer, which is the common case for small files on a
// LAN or loopback link. time.Duration's own formatting already picks a
// sensible unit (ms, µs, ns) there, so it's used directly.
func formatDuration(seconds float64) string {
	if seconds < 0 {
		seconds = 0
	}

	if seconds > 0 && seconds < 1 {
		return time.Duration(seconds * float64(time.Second)).String()
	}

	total := int64(seconds + 0.5) // round to the nearest second
	days := total / 86400
	total %= 86400
	hours := total / 3600
	total %= 3600
	minutes := total / 60
	secs := total % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 || days > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 || hours > 0 || days > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	parts = append(parts, fmt.Sprintf("%ds", secs))

	return strings.Join(parts, "")
}

var durationPattern = regexp.MustCompile(`(\d+(?:\.\d+)?)(d|h|m|s)`)

var durationUnitSeconds = map[string]float64{
	"d": 86400,
	"h": 3600,
	"m": 60,
	"s": 1,
}

// parseDuration turns a string like "1h30m", "90s", or "3d4h" into a
// number of seconds. A bare number is assumed to already be seconds.
func parseDuration(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("invalid duration %q", s)
	}

	if n, err := strconv.ParseFloat(s, 64); err == nil {
		return n, nil
	}

	matches := durationPattern.FindAllStringSubmatch(s, -1)
	if matches == nil {
		return 0, fmt.Errorf("invalid duration %q: expected something like 1h30m, 90s, or 3d4h", s)
	}

	// Require the matches to account for the whole string, so a typo like
	// "1h30x" doesn't silently parse as "1h" with the rest ignored.
	var rebuilt strings.Builder
	for _, m := range matches {
		rebuilt.WriteString(m[0])
	}
	if rebuilt.String() != s {
		return 0, fmt.Errorf("invalid duration %q: expected something like 1h30m, 90s, or 3d4h", s)
	}

	seen := make(map[string]bool, len(matches))
	var total float64
	for _, m := range matches {
		if seen[m[2]] {
			return 0, fmt.Errorf("invalid duration %q: unit %q repeated", s, m[2])
		}
		seen[m[2]] = true

		n, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid duration %q: %w", s, err)
		}
		total += n * durationUnitSeconds[m[2]]
	}

	return total, nil
}
