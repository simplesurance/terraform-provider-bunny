#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "$0")" && pwd -P)"
repo_root="$script_dir/.."

cd "$repo_root"
go mod tidy

changes="$(git diff --exit-code "vendor" "go.mod" "go.sum" || true)"
if [ -n "$changes" ]; then
	echo "ERROR: go module definitions is not tidy"
	echo "       Run 'go mod tidy' and commit the changes."
	echo
	echo "      Changes:"
	echo "$changes"
	exit 1
fi

echo "Go module definition is tidy"
