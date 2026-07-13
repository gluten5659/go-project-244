# gendiff

[![Actions Status](https://github.com/gluten5659/go-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/gluten5659/go-project-244/actions)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=gluten5659_go-project-244&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=gluten5659_go-project-244)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=gluten5659_go-project-244&metric=coverage)](https://sonarcloud.io/summary/new_code?id=gluten5659_go-project-244)

Compare two config files and see what changed. Works with JSON and YAML,
handles nested structures, and prints the diff in `stylish`, `plain`, or `json` format.

## Demo

![gendiff demo](demo/demo.gif)

## Install

```bash
make build
```

You get the binary at `bin/gendiff`.

## Usage

```bash
bin/gendiff [--format <fmt>] <first-file> <second-file>
```

Default format is `stylish`:

```bash
bin/gendiff examples/before.json examples/after.json
```

Switch the format with `--format` (or `-f`) — `stylish`, `plain`, or `json`:

```bash
bin/gendiff --format plain examples/before.json examples/after.json
```

The `json` format is machine-readable — the whole diff tree is wrapped in a
`diff` array, and every node carries its `key`, `type`, and value:

```bash
bin/gendiff --format json examples/before.json examples/after.json
```

```json
{
  "diff": [
    {
      "key": "follow",
      "type": "removed",
      "value": false
    },
    {
      "key": "host",
      "type": "unchanged",
      "value": "hexlet.io"
    },
    {
      "key": "proxy",
      "type": "removed",
      "value": "123.234.53.22"
    },
    {
      "children": [
        {
          "key": "cache",
          "newValue": false,
          "oldValue": true,
          "type": "updated"
        },
        {
          "key": "ttl",
          "type": "unchanged",
          "value": 60
        }
      ],
      "key": "settings",
      "type": "nested"
    },
    {
      "key": "timeout",
      "newValue": 20,
      "oldValue": 50,
      "type": "updated"
    },
    {
      "key": "verbose",
      "type": "added",
      "value": true
    }
  ]
}
```

JSON and YAML mix freely, so the two files don't have to share a format:

```bash
bin/gendiff examples/before.yaml examples/after.yaml
```

## Exit codes

`gendiff` follows the BSD `sysexits` conventions, so scripts can tell failures
apart:

| Code | Meaning |
|------|---------|
| 0  | Success |
| 1  | Unexpected error |
| 64 | Usage error — bad flag, unsupported format, or wrong number of paths |
| 65 | Data error — malformed or unsupported input file |
| 66 | Input file not found |
| 74 | I/O error — reading a file or writing the result |
| 77 | Permission denied |

## Development

```bash
make test           # tests with the race detector
make test-coverage  # tests + coverage.out
make lint           # golangci-lint
make fmt            # format
make demo           # rebuild and regenerate demo/demo.gif
```
