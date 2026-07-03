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
gendiff [--format <fmt>] <first-file> <second-file>
```

Default format is `stylish`:

```bash
gendiff sample/before.json sample/after.json
```

Switch the format with `--format` (or `-f`) — `stylish`, `plain`, or `json`:

```bash
gendiff --format plain sample/before.json sample/after.json
```

JSON and YAML mix freely, so the two files don't have to share a format:

```bash
gendiff sample/before.yaml sample/after.yaml
```

## Development

```bash
make test           # tests with the race detector
make test-coverage  # tests + coverage.out
make lint           # golangci-lint
make fmt            # format
make demo           # rebuild and regenerate demo/demo.gif
```
