# transfer-eta

You know the size of a file and the speed of a link, and you want to know
how long the transfer will take. This shouldn't require mental math with
GB vs GiB or a trip to a browser calculator, so `transfer-eta` does the one
conversion: size + rate -> duration.

```
$ transfer-eta --size 4.7GB --rate 25MB/s
4.70 GB at 25.00 MB/s => 3m8s
```

## Usage

```
transfer-eta --size <size> --rate <rate> [--json]
```

- `--size` accepts decimal units (B, KB, MB, GB, TB, PB, 1000-based) and
  binary units (KiB, MiB, GiB, TiB, PiB, 1024-based). A bare number is
  treated as bytes.
- `--rate` accepts a size followed by `/s` (also `/sec`, `/second`),
  e.g. `25MB/s`, `650KiB/s`, `1.5GB/s`.

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

- Solving for size or rate instead of duration (e.g. "I have 2 hours and
  a 4.7GB file, what rate do I need?")
- Rates quoted in bits per second (Mbps, Gbps) — everything is bytes for now
- Sub-second precision for very fast, very small transfers

## License

MIT, see [LICENSE](LICENSE).
