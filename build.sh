#!/usr/bin/env bash
set -eu

export WORKDIR="$(realpath $(dirname "$0"))"
export GOPATH="$WORKDIR/.build/gopath"
export PATH="$GOPATH/bin:$PATH"
export DISTDIR="$WORKDIR/dist"

export VERSION="$(git describe --first-parent --always)"

rm -rf "$GOPATH"
rm -rf "$DISTDIR"
mkdir -p "$GOPATH"

export GOX_TEMPLATE="$DISTDIR/{{.Dir}}-{{.OS}}-{{.Arch}}-$VERSION/reviewhoo"
export GOX_ARCH="amd64"
export GOX_OS="linux darwin"

go get github.com/mitchellh/gox

mkdir -p "$GOPATH/src/github.com/apparentlymart"
ln -s "$WORKDIR" "$GOPATH/src/github.com/apparentlymart/reviewhoo"

go get -v github.com/apparentlymart/reviewhoo

gox -arch="$GOX_ARCH" -os="$GOX_OS" -output="$GOX_TEMPLATE" github.com/apparentlymart/reviewhoo

cd "$DISTDIR"
for dir in *"-$VERSION"; do
    cp ../README.md "$dir/README.md"
    tar jcf "$dir.tar.bz2" "$dir/README.md"
    echo "dist/$dir.tar.bz2"
done
