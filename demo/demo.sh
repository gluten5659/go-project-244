#!/usr/bin/env bash

prompt='$ '
type_delay="${TYPE_DELAY:-0.04}"
pause_after="${PAUSE:-1.2}"

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

type_and_run 'bin/gendiff --help'
type_and_run 'bin/gendiff sample/before.json sample/after.json'
type_and_run 'bin/gendiff --format plain sample/before.json sample/after.json'
type_and_run 'bin/gendiff --format json sample/before.json sample/after.json'
type_and_run 'bin/gendiff sample/before.yaml sample/after.yaml'
type_and_run 'bin/gendiff no-such-file sample/after.json; echo "exit code: $?"'
