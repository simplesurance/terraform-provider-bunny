#!/usr/bin/env bash

set -euo pipefail

usage() {
	echo "$(basename "$0") PROVIDER-DIR"
	echo
	echo "create a terraform config to use a local provider bin"
}


if [ $# -ne 1 ]; then
	echo "ERR: missing argument " >&2
	usage
	exit 1
fi

if [[ "$1" == "-h" || "$1" == "--help" ]]; then
	usage
	exit 0
fi

provider_dir="$(realpath -m "$1")"

cfg='provider_installation {
  dev_overrides  {
    "simplesurance/bunny" = "%s"
  }
  direct {}
}'

printf "$cfg" "$provider_dir"
