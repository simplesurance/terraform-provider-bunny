#!/usr/bin/env bash

set -euo pipefail

changelog_file="CHANGELOG.md"
gpg_sign_key_id="0xC8B381683DBCEDFE"

err() {
	echo "ERR:" "$@" >&2
}

fatal() {
	err "$@"
	exit 1
}


set_unreleased_date() (
	if ! grep -qxi "^## $version (Unreleased)" "$changelog_file";then
		fatal "unreleased entry for ver $version missing in $changelog_file"
	fi

	ver_date="$(env LANG=en_US.UTF-8 date '+%B %d, %Y')"

	sed  -i "s/## $version (Unreleased)/## $version ($ver_date)/g" "$changelog_file"
)

add_new_release_line() (
	next_ver="$(echo "${version}" | awk -F. -v OFS=. '{$NF++;print}')"
	echo -e "## $next_ver (Unreleased)\n\n$(cat "$changelog_file")" > "$changelog_file"
)

git_tag_exists() {
	git describe --abbrev=0 --match="$tag" --tags &>/dev/null
}

create_git_tag() (
	if ! git describe --abbrev=0 --match="$version" --tags; then
		git tag -u "$gpg_sign_key_id" -s "$tag" -m "version $version"
		echo "git tag $tag created"
	else
		fatal "git tag $tag already exists"
	fi
)

git_worktree_is_clean() (
	if git diff-files --quiet; then
		return 0
	fi

	[ -z "$(git ls-files --other --directory --exclude-standard)" ]
)

usage() {
	echo "usage: $(basename "$0") VERSION"
}

## main

if [ $# -ne 1 ]; then
	err "missing commandline argument"
	echo
	usage >&2
	exit 1
fi

if [[ "$1" = "-h" || "$1" = "--help" ]]; then
	usage
	exit 0
fi

if ! command -v goreleaser &> /dev/null; then
	fatal "goreleaser command not found"
fi

if [[ ! -v GITHUB_TOKEN || -z "$GITHUB_TOKEN" ]]; then
	fatal "GITHUB_TOKEN environment variable is not set"
fi

version="$1"
tag="v$version"

if ! git_worktree_is_clean; then
	fatal "git work tree is dirty, remove or commit changes, see 'git status' output"
fi

if git_tag_exists; then
	fatal "git tag already exists"
fi

set_unreleased_date
git commit -m "changelog: release $version" "$changelog_file"
create_git_tag

GPG_RELEASE_SIGN_KEY_ID="$gpg_sign_key_id" goreleaser release --rm-dist

echo "Github draft release created, please finalize and publish it on the webpage"

add_new_release_line
echo "Entry in $changelog_file created for next release"
git commit -m "changelog: add entry for next version"

git push
git push --tags
