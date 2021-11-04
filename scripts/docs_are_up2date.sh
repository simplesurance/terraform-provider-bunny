#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "$0")" && pwd -P)"
repo_root="$script_dir/.."

cd "$repo_root"

make docs
echo '--------------------------------------------------------------'
echo

out="$(git status -s docs || true)"
if [ -n "$out" ]; then
	echo 'ERROR: docs/ is not up to date!'
	echo 'The following files are outdated:'
	echo
	echo '--------------------------------------------------------------'
	printf "%s\n" "$out" | sed 's/^/ * /g'
	echo '--------------------------------------------------------------'
	git diff
	echo '--------------------------------------------------------------'
	echo
	echo 'Please run 'make docs' and commit the changes.'
	exit 1
fi

echo 'docs/ are up to date'
