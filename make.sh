#!/bin/bash

set -ex

cwd=$PWD

for goos in darwin linux windows; do
	for goarch in 386 amd64; do
		(cd ${0%/*} && GOOS=$goos GOARCH=$goarch go build -o $cwd/grpn_${goos}_$goarch)
	done
done
