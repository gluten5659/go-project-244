#!/usr/bin/env bash

prompt='$ '
type_delay="${TYPE_DELAY:-0.04}"
pause_after="${PAUSE:-1.2}"

sample_dir='sample'

setup_sample() {
	mkdir -p "$sample_dir"

	cat > "$sample_dir/before.json" <<'JSON'
{
  "host": "hexlet.io",
  "timeout": 50,
  "proxy": "123.234.53.22",
  "follow": false,
  "settings": {
    "cache": true,
    "ttl": 60
  }
}
JSON

	cat > "$sample_dir/after.json" <<'JSON'
{
  "host": "hexlet.io",
  "timeout": 20,
  "verbose": true,
  "settings": {
    "cache": false,
    "ttl": 60
  }
}
JSON

	cat > "$sample_dir/before.yaml" <<'YAML'
host: hexlet.io
timeout: 50
follow: false
YAML

	cat > "$sample_dir/after.yaml" <<'YAML'
host: hexlet.io
timeout: 20
verbose: true
YAML
}

cleanup_sample() {
	rm -rf "$sample_dir"
}

type_and_run() {
	local command="$1"
	printf '%s' "$prompt"
	for (( index = 0; index < ${#command}; index++ )); do
		printf '%s' "${command:index:1}"
		sleep "$type_delay"
	done
	printf '\n'
	sleep 0.3
	eval "$command"
	printf '\n'
	sleep "$pause_after"
}

trap cleanup_sample EXIT
setup_sample

type_and_run 'bin/gendiff --help'
type_and_run 'bin/gendiff sample/before.json sample/after.json'
type_and_run 'bin/gendiff --format plain sample/before.json sample/after.json'
type_and_run 'bin/gendiff sample/before.yaml sample/after.yaml'
type_and_run 'bin/gendiff no-such-file sample/after.json; echo "exit code: $?"'
