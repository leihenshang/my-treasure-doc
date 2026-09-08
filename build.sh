#!/bin/bash

VERSION=v0.0.1

docker build -t treasure-doc-${VERSION} .
docker save -o treasure-doc-${VERSION}.tar.gz treasure-doc-${VERSION}