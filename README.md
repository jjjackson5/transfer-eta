# transfer-eta

You know two of {file size, link speed, how long you're willing to wait}
and want the third. This shouldn't require mental math with GB vs GiB or
a trip to a browser calculator, so `transfer-eta` does the conversion:
give it any two of size, rate, and duration, and it solves for the one
you left out.

```
$ transfer-eta --size 4.7GB --rate 25MB/s
4.70 GB at 25.00 MB/s => 3m8s

$ transfer-eta --size 4.7GB --duration 1h
4.70 GB in 1h0m0s => 1.31 MB/s

$ transfer-eta --rate 25MB/s --duration 3m8s
25.00 MB/s for 3m8s => 4.70 GB
```

## Usage

```
transfer-eta --size <size> --rate <rate> [--json]
transfer-eta --size <size> --duration <duration> [--json]
transfer-eta --rate <rate> --duration <duration> [--json]
```

Provide exactly two of `--size`, `--rate`, and `--duration`; the third is
solved for.

- `--size` accepts decimal units (B, KB, MB, GB, TB, PB, 1000-based) and
  binary units (KiB, MiB, GiB, TiB, PiB, 1024-based). A bare number is
  treated as bytes.
- `--rate` accepts a size followed by `/s` (also `/sec`, `/second`),
  e.g. `25MB/s`, `650KiB/s`, `1.5GB/s`.
- `--duration` accepts a combination of `d`, `h`, `m`, `s` units, e.g.
  `1h30m`, `90s`, `3d4h`. A bare number is treated as seconds.

### JSON output

```
$ transfer-eta --size 650MiB --rate 12MB/s --json
{
  "size_bytes": 681574400,
  "rate_bytes_per_second": 12000000,
  "duration_seconds": 56.79786666666667,
  "duration_human": "57s"
}
```

The JSON form is meant for scripting: pipe it into `jq` or read it from
another program instead of parsing the human-readable line.

## Building

```
go build -o transfer-eta .
```

No external dependencies — standard library only.

## Not supported yet

- Rates quoted in bits per second (Mbps, Gbps) — everything is bytes for now
- Sub-second precision for very fast, very small transfers

## License

MIT, see [LICENSE](LICENSE).
