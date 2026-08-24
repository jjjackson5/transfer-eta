package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Decimal units follow the SI convention (1 KB = 1000 B), binary units
// follow IEC (1 KiB = 1024 B). Both are common in the wild depending on
// whether a tool reports disk usage or network throughput, so both need
// to parse cleanly.
var sizeUnits = map[string]float64{
	"B":   1,
	"KB":  1e3,
	"MB":  1e6,
	"GB":  1e9,
	"TB":  1e12,
	"PB":  1e15,
	"KIB": 1024,
	"MIB": 1024 * 1024,
	"GIB": 1024 * 1024 * 1024,
	"TIB": 1024 * 1024 * 1024 * 1024,
	"PIB": 1024 * 1024 * 1024 * 1024 * 1024,
}

var sizePattern = regexp.MustCompile(`^([0-9]*\.?[0-9]+)\s*([A-Za-z]*)$`)

// parseSize turns a string like "4.7GB", "650MiB", or a bare "1024"
// (assumed bytes) into a byte count.
func parseSize(s string) (float64, error) {
	s = strings.TrimSpace(s)
	m := sizePattern.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("invalid size %q", s)
	}

	n, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size %q: %w", s, err)
	}

	unit := strings.ToUpper(m[2])
	if unit == "" {
		unit = "B"
	}
	mult, ok := sizeUnits[unit]
	if !ok {
		return 0, fmt.Errorf("unknown unit %q in %q", m[2], s)
	}

	return n * mult, nil
}

// formatSize renders a byte count using decimal (SI) prefixes, since
// that's what the input side accepts by default and what most transfer
// speeds are quoted in.
func formatSize(bytes float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	value := bytes
	i := 0
	for value >= 1000 && i < len(units)-1 {
		value /= 1000
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%.0f %s", value, units[i])
	}
	return fmt.Sprintf("%.2f %s", value, units[i])
}
