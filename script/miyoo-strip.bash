#!/bin/bash

set -eu

function die() {
	echo "DIE: $*" >&2
	exit 1
}

for m in */miyoogamelist.xml; do
	d=$(dirname "$m")
	g="$d/gamelist.xml"

	# Assert things are as expected
	[ -f "$m" ] || die "The output file $m is missing"
	[ ! -f "$g" ] || die "The destination file $g already exists"

	mv "$m" "$g"

	xml edit \
		-d "//desc" \
		-d "//rating" \
		-d "//genre" \
		-d "//players" \
		-d "//releasedate" \
		-d "//developer" \
		-d "//publisher" \
		-d "//hash" \
		-d "//thumbnail" \
		-d "//genreid" \
		--subnode "gameList/game[not(image)]" \
		-t elem \
		-n image \
		-v "no-img.png" \
		"$g" >"$m"
done
