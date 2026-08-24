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

func main() {
	sizeFlag := flag.String("size", "", "transfer size, e.g. 4.7GB or 650MiB")
	rateFlag := flag.String("rate", "", "transfer rate, e.g. 25MB/s")
	jsonFlag := flag.Bool("json", false, "output machine-readable JSON instead of a human-readable line")
	flag.Parse()

	if *sizeFlag == "" || *rateFlag == "" {
		fmt.Fprintln(os.Stderr, "usage: transfer-eta --size 4.7GB --rate 25MB/s [--json]")
		flag.PrintDefaults()
		os.Exit(2)
	}

	bytes, err := parseSize(*sizeFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "transfer-eta:", err)
		os.Exit(1)
	}

	bytesPerSec, err := parseRate(*rateFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "transfer-eta:", err)
		os.Exit(1)
	}
	if bytesPerSec <= 0 {
		fmt.Fprintln(os.Stderr, "transfer-eta: rate must be greater than zero")
		os.Exit(1)
	}

	seconds := bytes / bytesPerSec
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

	fmt.Printf("%s at %s/s => %s\n", formatSize(bytes), formatSize(bytesPerSec), human)
}
