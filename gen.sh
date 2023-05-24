#!/usr/bin/env bash
set -euo pipefail

require() {
  if ! command -v "$1" &>/dev/null && [[ -n "$2" ]]; then
    go install "$2"
  fi
  if ! command -v "$1" >/dev/null; then
    echo >&2 "[gen.sh] $1 not found"
    if [[ -n "$2" ]]; then
      echo >&2 "[gen.sh] Is \${GOPATH}/bin in your \$PATH?"
    fi
    exit 1
  fi
}

#require protoc
require buf
require protoc-gen-go google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
require protoc-gen-go-grpc google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
require protoc-gen-gotf github.com/travix/protoc-gen-gotf@latest

echo "[gen.sh] generating protobuf and grpc code"
buf mod update
buf generate
## without buf
# module="github.com/travix/gotf-example"
# protoc -I. --go_out=. --go_opt module=${module} --go-grpc_out=. --go-grpc_opt module=${module} example.proto

echo "[gen.sh] generating terraform go code"
buf generate --template buf.gen.tf.yaml
## without buf
# protoc -I. --gotf_out=. --gotf_opt=log_level=debug --gotf_opt module=${module} example.proto
