#!/bin/bash

if [ $# -eq 0 ]
then
  echo "Usage: $0 version.number"
  exit
fi

version=$1

OSS=(darwin freebsd linux)
ARCHS=(386 amd64)

mkdir -p dist
rm -f dist/botella*

for os in "${OSS[@]}"; do
  for arch in "${ARCHS[@]}"; do
    echo "Building for $os($arch)"
    GOOS=$os GOARCH=$arch go build
    mv botella dist/botella_${version}_$os-$arch
  done
done
