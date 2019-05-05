#!/bin/sh
VERSION=`cat sde.version`
URL="https://www.fuzzwork.co.uk/dump/sde-${VERSION}/postgres-${VERSION}-schema.dmp.bz2"

echo $URL

curl ${URL} -o sde-$VERSION.bz2
bunzip2 sde-$VERSION.bz2
