#!/bin/bash -eu

source /build-common.sh

BINARY_NAME="modemrebooter"
COMPILE_IN_DIRECTORY="cmd/modemrebooter"
BINTRAY_PROJECT="joonas/modemrebooter"

standardBuildProcess
