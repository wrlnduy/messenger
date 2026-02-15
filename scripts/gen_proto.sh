#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

PROTO_DIR="$SCRIPT_DIR/../proto"
OUT_DIR="$PROTO_DIR"

PROTO_FILES=$(find "$PROTO_DIR" -type f -name "*.proto")

protoc \
  --proto_path="$PROTO_DIR" \
  --go_out="$OUT_DIR" --go_opt=paths=source_relative \
  --go-grpc_out="$OUT_DIR" --go-grpc_opt=paths=source_relative \
  $PROTO_FILES