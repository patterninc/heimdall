#!/bin/bash

# verbose
set -ex

# setup vars
BUILD_VERSION=$(git describe --tags --abbrev=0 --match "v*" | cut -dv -f2)
WORKING_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
OUTPUT_DIR="${WORKING_DIR}/dist"
LDFLAGS="-X main.Build=${BUILD_VERSION}"
BUILD_GO=false
BUILD_WEB=false
RUN_TESTS=false
for arg in "$@"; do
  case "${arg}" in
    --go) BUILD_GO=true ;;
    --ui) BUILD_WEB=true ;;
    --test) RUN_TESTS=true ;;
    *) echo "unknown flag: ${arg} (use --go, --ui, --test)" >&2; exit 1 ;;
  esac
done

# reset output dir
rm -rf ${OUTPUT_DIR} > /dev/null 2>&1
mkdir -p ${OUTPUT_DIR}/web ${OUTPUT_DIR}/plugins

if [ "${RUN_TESTS}" = "true" ]; then
  echo "running tests"
  (cd ${WORKING_DIR} && go test ./internal/pkg/... && go test ./pkg/...)
fi

if [ "${BUILD_GO}" = "true" ]; then
  COMMANDS=$(ls ${WORKING_DIR}/cmd)
  PLUGINS=$(ls ${WORKING_DIR}/plugins)



  # go environment
  (cd ${WORKING_DIR} && go version)

  # build commands
  (cd ${WORKING_DIR} && echo "${COMMANDS}" | xargs -P 4 -I {} bash -c 'echo "building command: {}" && go build -ldflags "${LDFLAGS}" -o dist/{} cmd/{}/{}.go')

  # build plugins
  (cd ${WORKING_DIR} && echo "${PLUGINS}" | xargs -P 4 -I {} bash -c 'echo "building plugin: {}" && go build -buildmode=plugin -ldflags "${LDFLAGS}" -o dist/plugins/{}.so plugins/{}/{}.go')
fi

if [ "${BUILD_WEB}" = "true" ]; then
  # build web
  (cd ${WORKING_DIR}/web && corepack enable && pnpm install --frozen-lockfile --ignore-scripts && pnpm run build)
fi
