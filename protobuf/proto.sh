# go install github.com/micro/micro/v5/cmd/protoc-gen-micro@latest
#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

PROTO_DIR="$SCRIPT_DIR/pb"
OUT_DIR="$ROOT_DIR/protobuf/gen"

PROTO_FILES=("$PROTO_DIR"/*.proto)

protoc \
  -I="$PROTO_DIR" \
  --go_out=paths=source_relative:"$OUT_DIR" \
  --micro_out=paths=source_relative:"$OUT_DIR" \
  "${PROTO_FILES[@]}"
