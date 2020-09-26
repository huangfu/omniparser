#!/bin/bash
SCRIPT_DIR=$(pwd `dirname "$0"`)
go run -ldflags "-X main.gitCommit=$(git rev-parse HEAD) -X main.buildEpochSec=$(date +%s)" $SCRIPT_DIR/cli/op.go "$@"