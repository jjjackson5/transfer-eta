package main

import (
	"fmt"
	"strings"
)

// formatDuration renders a number of seconds as a compact human string
// like "1h2m3s" or "3d4h". time.Duration isn't used here because transfer
// times can run into weeks at low enough rates, well past what most
// people want to see expressed in nanosecond-precision Duration units.
func formatDuration(seconds float64) string {
	if seconds < 0 {
		seconds = 0
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
