package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type result struct {
	SizeBytes    float64 `json:"size_bytes"`
	RateBytesSec float64 `json:"rate_bytes_per_second"`
	Seconds      float64 `json:"duration_seconds"`
	Human        string  `json:"duration_human"`
}

const usage = "usage: transfer-eta --size <size> --rate <rate> [--json]\n" +
	"       transfer-eta --size <size> --duration <duration> [--json]\n" +
	"       transfer-eta --rate <rate> --duration <duration> [--json]\n"

func main() {
	sizeFlag := flag.String("size", "", "transfer size, e.g. 4.7GB or 650MiB")
	rateFlag := flag.String("rate", "", "transfer rate, e.g. 25MB/s")
	durationFlag := flag.String("duration", "", "transfer duration, e.g. 1h30m or 90s")
	jsonFlag := flag.Bool("json", false, "output machine-readable JSON instead of a human-readable line")
	flag.Parse()

	provided := 0
	for _, s := range []string{*sizeFlag, *rateFlag, *durationFlag} {
		if s != "" {
			provided++
		}
	}
	if provided != 2 {
		fmt.Fprint(os.Stderr, usage)
		fmt.Fprintln(os.Stderr, "provide exactly two of --size, --rate, --duration; the third is solved for")
		flag.PrintDefaults()
		os.Exit(2)
	}

	var (
		bytes       float64
		bytesPerSec float64
		seconds     float64
		label       string
		err         error
	)

	switch {
	case *sizeFlag != "" && *rateFlag != "":
		bytes, err = parseSize(*sizeFlag)
		if err == nil {
			bytesPerSec, err = parseRate(*rateFlag)
		}
		if err == nil && bytesPerSec <= 0 {
			err = fmt.Errorf("rate must be greater than zero")
		}
		if err == nil {
			seconds = bytes / bytesPerSec
		}
		label = "duration"

	case *sizeFlag != "" && *durationFlag != "":
		bytes, err = parseSize(*sizeFlag)
		if err == nil {
			seconds, err = parseDuration(*durationFlag)
		}
		if err == nil && seconds <= 0 {
			err = fmt.Errorf("duration must be greater than zero")
		}
		if err == nil {
			bytesPerSec = bytes / seconds
		}
		label = "rate"

	default: // rate + duration
		bytesPerSec, err = parseRate(*rateFlag)
		if err == nil && bytesPerSec <= 0 {
			err = fmt.Errorf("rate must be greater than zero")
		}
		if err == nil {
			seconds, err = parseDuration(*durationFlag)
		}
		if err == nil && seconds <= 0 {
			err = fmt.Errorf("duration must be greater than zero")
		}
		if err == nil {
			bytes = bytesPerSec * seconds
		}
		label = "size"
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "transfer-eta:", err)
		os.Exit(1)
	}

	human := formatDuration(seconds)

	if *jsonFlag {
		out := result{
			SizeBytes:    bytes,
			RateBytesSec: bytesPerSec,
			Seconds:      seconds,
			Human:        human,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			fmt.Fprintln(os.Stderr, "transfer-eta:", err)
			os.Exit(1)
		}
		return
	}

	switch label {
	case "rate":
		fmt.Printf("%s in %s => %s/s\n", formatSize(bytes), human, formatSize(bytesPerSec))
	case "size":
		fmt.Printf("%s/s for %s => %s\n", formatSize(bytesPerSec), human, formatSize(bytes))
	default:
		fmt.Printf("%s at %s/s => %s\n", formatSize(bytes), formatSize(bytesPerSec), human)
	}
}
