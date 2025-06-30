#!/bin/bash

# verbose
set -ex

# setup vars
BUILD_VERSION=$(git describe --tags --abbrev=0 --match "v*" | cut -dv -f2)
WORKING_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
OUTPUT_DIR="${WORKING_DIR}/dist"
COMMANDS=$(ls ${WORKING_DIR}/cmd)
PLUGINS=$(ls ${WORKING_DIR}/plugins)
LDFLAGS="-X main.Build=${BUILD_VERSION}"

# reset output dir
rm -rf ${OUTPUT_DIR} > /dev/null 2>&1
mkdir -p ${OUTPUT_DIR}/web ${OUTPUT_DIR}/plugins

# go environment
(cd ${WORKING_DIR} && go version)

# run tests
(cd ${WORKING_DIR} && go test ./internal/pkg/... && go test ./pkg/...)

# build commands
(cd ${WORKING_DIR} && for item in ${COMMANDS}; do go build -ldflags "${LDFLAGS}" -o dist/$item cmd/$item/$item.go; done)

# build plugins
(cd ${WORKING_DIR} && for item in ${PLUGINS}; do go build -buildmode=plugin -ldflags "${LDFLAGS}" -o dist/plugins/$item.so plugins/$item/$item.go; done)

# build web
(cd ${WORKING_DIR}/web && rm -rf node_modules > /dev/null 2>&1 && npm i && npm run build)
