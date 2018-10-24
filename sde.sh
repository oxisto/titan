#!/bin/sh
VERSION=`cat sde.version`
URL="https://cdn1.eveonline.com/data/sde/tranquility/sde-${VERSION}.zip"
FILES="sde/fsd/blueprints.yaml sde/fsd/typeIDs.yaml sde/fsd/groupIDs.yaml sde/fsd/categoryIDs.yaml sde/bsd/invMetaTypes.yaml"

curl ${URL} -o sde.zip
unzip -o sde.zip ${FILES}
